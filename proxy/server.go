package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "flowright/proxy/reece.ooo/flowright"
)

var (
	version          = "0.1.0"
	name             = flag.String("name", "proxy-1", "The unique name of this proxy")
	controlPlaneHost = flag.String("control-plane-host", "127.0.0.1", "The host for the control plane.")
	controlPlanePort = flag.Uint64("control-plane-port", 50051, "The port for the control plane.")
	host             = flag.String("host", "127.0.0.1", "The host to listen on.")
	publicHost       = flag.String("public-host", "127.0.0.1", "The public endpoint of this proxy.")
	port             = flag.Uint64("port", 9000, "The proxy port to listen on.")
	grpcPort         = flag.Uint64("grpc-port", 9001, "The grpc port to listen on.")
	db_path          = flag.String("db-path", "proxy.db", "Path to the database.")
	upgrader         = websocket.FastHTTPUpgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin:       func(ctx *fasthttp.RequestCtx) bool { return true },
		EnableCompression: true,
	}
	db *gorm.DB
)

type ProxyServer struct {
	pb.UnimplementedProxyServer
}

func (s *ProxyServer) GetInfo(ctx context.Context, in *pb.Empty) (*pb.ProxyInfo, error) {
	return &pb.ProxyInfo{
		Name:      *name,
		Host:      *publicHost,
		GrpcPort:  uint32(*grpcPort),
		ProxyPort: uint32(*port),
		Version:   version,
	}, nil
}

func (s *ProxyServer) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	return &pb.HeartbeatResponse{}, nil
}

func (s *ProxyServer) CreateOrUpdateRoute(ctx context.Context, in *pb.RoutingMap) (*pb.RoutingRule, error) {
	tx := db.Save(&Route{
		Owner:     in.Owner,
		Project:   in.Project,
		ProxyName: in.ProxyName,
		ProxyHost: in.ProxyHost,
		ProxyPort: int(in.ProxyPort),
	})
	if tx.Error != nil {
		return nil, status.Error(codes.Internal, "Failed to save route")
	}

	return &pb.RoutingRule{
		ProxyName: in.ProxyName,
		ProxyHost: in.ProxyHost,
		ProxyPort: in.ProxyPort,
	}, nil
}

func BuildAndRunEnvironment(owner string, project string, requirements string, sourceTarBall []byte) (string, string, int, error) {
	log.Printf("Got tarball of size: %d\n", len(sourceTarBall))
	// TODO: migrate to Podman Go Bindings
	// baseImage := "python:3.11-alpine3.18"
	baseImage := "python:3.11"

	// make sure the container is present
	pullCommand := exec.Command("podman", "pull", baseImage)
	pullCommand.Stderr = os.Stderr
	pullCommand.Stdout = os.Stdout
	err := pullCommand.Run()
	if err != nil {
		log.Println("Error pulling "+baseImage+" base image:", err)
		return "", "", -1, err
	}

	containerName := uniuri.NewLenChars(10, []byte("abcdefghijklmnopqrstuvxyz123456789")) + "-flowright-" + owner + "-" + project
	buildImageName := containerName + "-build"

	defer func() {
		exec.Command("podman", "stop", containerName).Run()
		exec.Command("podman", "rm", "-f", containerName).Run()
	}()

	// create a container
	runCommand := exec.Command("podman", "run", "-d", "--name", containerName, baseImage, "sleep", "180")
	runCommand.Stderr = os.Stderr
	runCommand.Stdout = os.Stdout
	if err := runCommand.Run(); err != nil {
		log.Println("Error creating container:", err)
		return "", "", -1, err
	}

	// install requirements
	tmpRequirementsFile := "/tmp/" + containerName + "_requirements.txt"
	if err := os.WriteFile(tmpRequirementsFile, []byte(requirements), 0644); err != nil {
		log.Println("Error writing requirements file:", err)
		return "", "", -1, err
	}

	copyCommand := exec.Command("podman", "cp", tmpRequirementsFile, containerName+":/requirements.txt")
	copyCommand.Stderr = os.Stderr
	copyCommand.Stdout = os.Stdout
	if err := copyCommand.Run(); err != nil {
		log.Println("Error copying requirements file:", err)
		return "", "", -1, err
	}

	installCommand := exec.Command("podman", "exec", containerName, "pip", "install", "-r", "/requirements.txt")
	installCommand.Stderr = os.Stderr
	installCommand.Stdout = os.Stdout
	if err := installCommand.Run(); err != nil {
		log.Println("Error installing requirements:", err)
		return "", "", -1, err
	}

	// ensure flowright is installed (this issue mostly occurs with flowright local development)
	execCommand := exec.Command("podman", "exec", containerName, "pip", "install", "flowright")
	execCommand.Stderr = os.Stderr
	execCommand.Stdout = os.Stdout
	if err := execCommand.Run(); err != nil {
		log.Println("Error installing flowright:", err)
		return "", "", -1, err
	}

	// copy source
	tmpSourceFile := "/tmp/" + containerName + "_source.tar.gz"
	if err := os.WriteFile(tmpSourceFile, sourceTarBall, 0644); err != nil {
		log.Println("Error writing source file:", err)
		return "", "", -1, err
	}

	copyCommand = exec.Command("podman", "cp", tmpSourceFile, containerName+":/source.tar.gz")
	copyCommand.Stderr = os.Stderr
	copyCommand.Stdout = os.Stdout
	if err := copyCommand.Run(); err != nil {
		log.Println("Error copying source:", err)
		return "", "", -1, err
	}

	execCommand = exec.Command("podman", "exec", containerName, "mkdir", "-p", "/flowright_app")
	execCommand.Stderr = os.Stderr
	execCommand.Stdout = os.Stdout
	if err := execCommand.Run(); err != nil {
		log.Println("Error creating app directory:", err)
		return "", "", -1, err
	}

	extractCommand := exec.Command("podman", "exec", containerName, "tar", "-xzf", "/source.tar.gz", "-C", "/flowright_app")
	extractCommand.Stderr = os.Stderr
	extractCommand.Stdout = os.Stdout
	if err := extractCommand.Run(); err != nil {
		log.Println("Failed to extract source:", err)
		return "", "", -1, err
	}

	// commit changes to image
	commitCommand := exec.Command("podman", "commit", "-p", containerName, buildImageName)
	commitCommand.Stdout = os.Stdout // TODO: capture this output
	commitCommand.Stderr = os.Stderr
	if err := commitCommand.Run(); err != nil {
		log.Println("Error committing changes:", err)
		return "", "", -1, err
	}

	// stop base container
	stopCommand := exec.Command("podman", "stop", containerName)
	stopCommand.Stderr = os.Stderr
	stopCommand.Stdout = os.Stdout
	if err := stopCommand.Run(); err != nil {
		log.Println("Error stopping container:", err)
		return "", "", -1, err
	}

	// now run app
	// NOTE: it appears impossible to expose a UDS from within the container. The proxy will nonetheless support UDS backends until a solution exists.

	// https://coolaj86.com/articles/how-to-test-if-a-port-is-available-in-go/
	portSearchStart := 41000
	portSearchEnd := 42000
	port := portSearchStart
	for ; port < portSearchEnd; port++ {
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			continue
		}
		ln.Close()
		break
	}

	containerRunName := containerName + "-run"
	runCommand = exec.Command("podman", "run", "-d", "-p", strconv.Itoa(port)+":8000", "--name", containerRunName, buildImageName, "flowright", "run", "/flowright_app", "--host=0.0.0.0")
	runCommand.Stderr = os.Stderr
	runCommand.Stdout = os.Stdout
	if err := runCommand.Run(); err != nil {
		log.Println("Error running app:", err)
		return "", "", -1, err
	}
	// install route
	// AddRoutingRule(owner, project, "localhost:"+strconv.Itoa(port), false, "")

	return buildImageName, containerRunName, port, nil
}

func (s *ProxyServer) CreateEnvironment(ctx context.Context, in *pb.EnvironmentCreateRequest) (*pb.EnvironmentCreateResponse, error) {
	imageName, containerName, port, err := BuildAndRunEnvironment(in.Owner, in.Project, in.Requirements, in.ProjectTarball)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to create environment: %v\n", err))
	}

	tx := db.Save(&Environment{
		Owner:               in.Owner,
		Project:             in.Project,
		PodmanImageName:     imageName,
		PodmanContainerName: containerName,
		PodmanContainerHost: "localhost",
		PodmanContainerPort: port,
	})

	if tx.Error != nil {
		// delete container
		exec.Command("podman", "container", "stop", containerName).Run()
		exec.Command("podman", "container", "rm", containerName).Run()

		return nil, status.Error(codes.Internal, "Failed to save environment info")
	}

	return &pb.EnvironmentCreateResponse{
		Success:    true,
		Message:    fmt.Sprintf("Running container %s using image %s on port %d", imageName, containerName, port),
		ProxxyName: *name,
	}, nil
}

type ProxyURIParseError struct {
	Reason string
}

func (err *ProxyURIParseError) Error() string {
	return err.Reason
}

type ResolvedRoute struct {
	Host     string
	Port     int
	BasePath string
}

func (r *ResolvedRoute) Endpoint() string {
	return fmt.Sprintf("%s:%d%s", r.Host, r.Port, r.BasePath)
}

type Route struct {
	// gorm.Model
	Owner     string `gorm:"primaryKey;autoIncrement:false"`
	Project   string `gorm:"primaryKey;autoIncrement:false"`
	ProxyName string
	ProxyHost string
	ProxyPort int
}

func (r *Route) Endpoint() string {
	return fmt.Sprintf("%s:%d", r.ProxyHost, r.ProxyPort)
}

type Environment struct {
	// gorm.Model
	Owner               string `gorm:"primaryKey;autoIncrement:false"`
	Project             string `gorm:"primaryKey;autoIncrement:false"`
	PodmanImageName     string
	PodmanContainerName string
	PodmanContainerHost string
	PodmanContainerPort int
}

func getRoute(owner string, project string) (*ResolvedRoute, error) {
	var route Route

	tx := db.Where(&Route{
		Owner:   owner,
		Project: project,
	}).Order("created_at DESC").First(&route)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if route.ProxyName != *name {
		// redirect to other proxy
		return &ResolvedRoute{
			Host:     route.ProxyHost,
			Port:     route.ProxyPort,
			BasePath: fmt.Sprintf("/%s/%s", owner, project),
		}, nil
	}

	// lookup environment
	var environment Environment
	tx = db.Where(&Environment{
		Owner:   owner,
		Project: project,
	}).Order("created_at DESC").First(&environment)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ResolvedRoute{
		Host:     environment.PodmanContainerHost,
		Port:     environment.PodmanContainerPort,
		BasePath: "",
	}, nil
}

func initDB() error {
	db_ptr, err := gorm.Open(sqlite.Open(*db_path), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db = db_ptr
	return nil
}

func initContainers() error {
	tx := db.Find(&Environment{})
	if tx.Error != nil {
		log.Printf("Could not get local environments on startup! %v\n", tx.Error)
		return tx.Error
	}

	rows, err := tx.Rows()

	if err != nil {
		log.Printf("Could not get local environments on startup! %v\n", err)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var environment Environment
		tx.ScanRows(rows, &environment)
		log.Printf("Found environment %s/%s\n", environment.Owner, environment.Project)

		// startup environment
		runCommand := exec.Command("podman", "start", environment.PodmanContainerName)
		runCommand.Stdout = os.Stdout
		runCommand.Stderr = os.Stderr
		if err := runCommand.Run(); err != nil {
			log.Printf("Failed to start environment %s/%s\n", environment.Owner, environment.Project)
			// TODO: notify control plane
		}

		log.Printf("Started environment %s/%s\n", environment.Owner, environment.Project)
	}

	return nil
}

func initControlPlane() error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *controlPlaneHost, *controlPlanePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewRoutingControllerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = client.RegisterProxy(ctx, &pb.ProxyRegisterRequest{
		Name:     *name,
		Host:     *publicHost,
		Port:     uint32(*port),
		GrpcPort: uint32(*grpcPort),
	})

	if err != nil {
		return err
	}

	return nil
}

func serveGrpc(s *grpc.Server, lis net.Listener) {
	err := s.Serve(lis)
	if err != nil {
		log.Fatal(err)
		panic("Failed to serve grpc")
	}
}

func main() {
	flag.Parse()
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Route{})
	db.AutoMigrate(&Environment{})

	err = initContainers()
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *grpcPort))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterProxyServer(s, &ProxyServer{})
	log.Printf("[GRPC] Listening at %v\n", lis.Addr())
	go serveGrpc(s, lis)

	time.Sleep(time.Second)

	err = initControlPlane()
	if err != nil {
		log.Fatal(err)
	}

	if err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), requestHandler); err != nil {
		fmt.Printf("Error occurred: %v", err)
	}
}

func getProxyPathComponents(uri *fasthttp.URI) (base string, project string, path string, err error) {
	capped_path := string(uri.Path()[1:])
	components := strings.SplitN(capped_path, "/", 3)
	if len(components) != 3 {
		return "", "", "", &ProxyURIParseError{Reason: "Invalid URI"}
	}
	return components[0], components[1], components[2] + uri.QueryArgs().String(), nil
}

func handleWebsocketProxyBackend(server_conn *websocket.Conn, client_conn *websocket.Conn, done chan struct{}, closed *bool) {
	defer func() {
		close(done)
		*closed = true
	}()

	for {
		_, message, err := client_conn.ReadMessage()
		if err != nil {
			if !*closed {
				log.Println("Error reading client socket:", err)
			}
			if websocket.IsCloseError(err) {
				ws_err := err.(*websocket.CloseError)
				server_conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(ws_err.Code, ws_err.Text), time.Now().Add(1*time.Second))
			} else {
				server_conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Backend service read error"), time.Now().Add(1*time.Second))
			}
			return
		}
		// log.Println(string(message))
		err = server_conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			if !*closed {
				log.Println("Error writing server socket:", err)
			}
			if websocket.IsCloseError(err) {
				ws_err := err.(*websocket.CloseError)
				client_conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(ws_err.Code, ws_err.Text), time.Now().Add(1*time.Second))
			} else {
				client_conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Web client write error"), time.Now().Add(1*time.Second))
			}
			return
		}
	}
}

func handleWebSocketProxyClient(server_conn *websocket.Conn, client_conn *websocket.Conn, done chan struct{}, closed *bool) {
	for {
		select {
		case <-done:
			log.Println("flowright backend closed connection")
			return
		default:
			_, message, err := server_conn.ReadMessage()

			if err != nil {
				if !*closed || websocket.IsUnexpectedCloseError(err) {
					log.Println("Error reading server socket:", err)
				}
				*closed = true
				return
			}
			err = client_conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				if !*closed || websocket.IsUnexpectedCloseError(err) {
					log.Println("Error writing client socket")
				}
				*closed = true
				return
			}
		}
	}
}

func handleWebsocketProxy(ctx *fasthttp.RequestCtx) {
	owner, project, path, err := getProxyPathComponents(ctx.URI())
	if err != nil {
		log.Println("Invalid URI:", string(ctx.URI().Path()))
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	route, err := getRoute(owner, project)
	if err != nil {
		log.Println("Invalid backend:", owner, project)
		log.Println("Cause:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	uri := fasthttp.AcquireURI()
	err = uri.Parse([]byte(route.Endpoint()), []byte(path))
	if err != nil {
		log.Println("Error parsing request uri:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	uri.SetScheme("ws")
	// upgrade connection
	err = upgrader.Upgrade(ctx, func(server_conn *websocket.Conn) {
		defer server_conn.Close()
		// connect to flowright backend
		dialer := &websocket.Dialer{}
		// if route.IsUnixSocket {
		// 	dialer = &websocket.Dialer{
		// 		NetDial: func(network, addr string) (net.Conn, error) {
		// 			return net.Dial("unix", route.UnixSocketPath) // TODO, if horizontally scaled this will be issue
		// 		},
		// 	}
		// } else {
		// 	dialer = websocket.DefaultDialer
		// }
		dialer = websocket.DefaultDialer

		client_conn, _, err := dialer.Dial(uri.String(), nil)
		if err != nil {
			log.Fatal("Error dialing websocket:", err)
			return
		}
		defer client_conn.Close()

		done := make(chan struct{})
		closed := false

		go handleWebsocketProxyBackend(server_conn, client_conn, done, &closed)
		handleWebSocketProxyClient(server_conn, client_conn, done, &closed)
	})
	if err != nil {
		log.Println("Error upgrading connection:", err)
	}
}

func handleHttpProxy(ctx *fasthttp.RequestCtx) {
	owner, project, path, err := getProxyPathComponents(ctx.URI())
	if err != nil {
		log.Println("Invalid URI:", string(ctx.URI().Path()))
		return
	}
	route, err := getRoute(owner, project)
	if err != nil {
		log.Println("Invalid backend:", owner, project)
		log.Println("Cause:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	log.Println("Endpoint:", route)
	uri := fasthttp.AcquireURI()
	err = uri.Parse([]byte(route.Endpoint()), []byte(path))
	if err != nil {
		log.Println("Error parsing request uri:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	defer fasthttp.ReleaseURI(uri)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(uri.String())
	log.Printf("Translated URI: %s -> %s\n", ctx.URI(), uri.String())
	ctx.Request.Header.VisitAll(func(key []byte, value []byte) {
		req.Header.Set(string(key), string(value))
	})
	resp := fasthttp.AcquireResponse()
	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			// if route.IsUnixSocket {
			// 	return net.Dial("unix", route.UnixSocketPath)
			// }
			return net.Dial("tcp", addr)
		},
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}
	err = client.Do(req, resp)
	if err != nil {
		log.Println("Error servicing request", err)
		ctx.Response.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}
	ctx.Response.SetBody(resp.Body())
	resp.Header.VisitAll(func(key []byte, value []byte) {
		ctx.Response.Header.Set(string(key), string(value))
	})
	fasthttp.ReleaseResponse(resp)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	if websocket.FastHTTPIsWebSocketUpgrade(ctx) {
		log.Printf("%s\t[websocket] %s", ctx.RemoteIP(), ctx.Path())
		handleWebsocketProxy(ctx)
	} else {
		log.Printf("%s\t%s %s", ctx.RemoteIP(), ctx.Method(), ctx.Path())
		handleHttpProxy(ctx)
	}
}
