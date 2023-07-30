package main

import (
	"context"
	"errors"
	"flag"
	pb "flowright/controlplane/reece.ooo/flowright"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	host    = flag.String("host", "127.0.0.1", "Host to listen on.")
	port    = flag.Int("port", 50051, "GRPC server port.")
	db_path = flag.String("db", "control_plane.db", "Path to the sqlite database.")
	db      *gorm.DB
)

type Proxy struct {
	gorm.Model
	Name             string `gorm:"primaryKey;autoIncrement:false"`
	Host             string
	Port             uint32
	GrpcPort         uint32
	MissedHeartbeats int
	Alive            bool
}

type Route struct {
	gorm.Model
	Owner     string `gorm:"primaryKey;autoIncrement:false"`
	Project   string `gorm:"primaryKey;autoIncrement:false"`
	ProxyName string
	Proxy     Proxy
}

type Environment struct {
	gorm.Model
	Owner     string `gorm:"primaryKey;autoIncrement:false"`
	Project   string `gorm:"primaryKey;autoIncrement:false"`
	ProxyName string
	Proxy     Proxy
}

type RoutingControllerServer struct {
	pb.UnimplementedRoutingControllerServer
}

func tryConnectProxy(host string, port uint32) bool {
	log.Printf("Connecting to %s:%d\n", host, port)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
		return false
	}
	defer conn.Close()

	c := pb.NewProxyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetInfo(ctx, &pb.Empty{})

	if err != nil {
		log.Print(err)
		return false
	}

	log.Printf("\n---\nProxy:\nName: %s\nHost: %s\nProxy Port: %d\nGRPC Port: %d\nVersion: %s\n---\n", r.Name, r.Host, r.ProxyPort, r.GrpcPort, r.Version)

	return true
}

func (s *RoutingControllerServer) RegisterProxy(ctx context.Context, in *pb.ProxyRegisterRequest) (*pb.ProxyRegisterResponse, error) {
	if !tryConnectProxy(in.Host, in.GrpcPort) {
		return &pb.ProxyRegisterResponse{Success: false}, status.Error(codes.Unavailable, "Could not connect to proxy grpc client")
	}

	tx := db.Save(&Proxy{
		Name:             in.Name,
		Host:             in.Host,
		Port:             in.Port,
		GrpcPort:         in.GrpcPort,
		MissedHeartbeats: 0,
		Alive:            true,
	})

	if tx.Error != nil {
		log.Fatal(tx.Error)
		return &pb.ProxyRegisterResponse{Success: false}, status.Error(codes.Unknown, "Failed to save proxy information to database")
	}

	return &pb.ProxyRegisterResponse{Success: true}, nil
}

func (s *RoutingControllerServer) GetInfo(ctx context.Context, in *pb.InfoRequest) (*pb.ControllerInfo, error) {
	return &pb.ControllerInfo{
		Host:     *host,
		GrpcPort: uint32(*port),
	}, nil
}

func (s *RoutingControllerServer) CreateOrUpdateRoute(ctx context.Context, in *pb.RoutingMap) (*pb.RoutingRule, error) {
	var proxy Proxy
	tx := db.Where(&Proxy{
		Name: in.ProxyName,
	}).First(&proxy)
	if tx.Error != nil {
		return nil, status.Error(codes.NotFound, "Invalid Proxy")
	}

	tx = db.Save(&Route{
		Owner:     in.Owner,
		Project:   in.Project,
		ProxyName: in.ProxyName,
	})
	if tx.Error != nil {
		return nil, status.Error(codes.Internal, "Failed to save route information to database")
	}

	// install route in proxy
	log.Printf("Installing route in %s\n", proxy.Name)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", proxy.Host, proxy.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
		log.Print("Failed to connect to proxy")
		tx.Rollback()
		return nil, status.Error(codes.Unavailable, "Could not connect to target proxy")
	}
	defer conn.Close()

	c := pb.NewProxyClient(conn)
	cctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.CreateOrUpdateRoute(cctx, &pb.RoutingMap{
		Owner:     in.Owner,
		Project:   in.Project,
		ProxyName: proxy.Name,
		ProxyHost: proxy.Host,
		ProxyPort: proxy.Port,
	})
	if err != nil {
		log.Print(err)
		log.Print("Failed to install route in target proxy")
		tx.Rollback()
		return nil, err
	}
	log.Printf("Successfully installed route in %s\n", proxy.Name)

	return &pb.RoutingRule{
		ProxyName: proxy.Name,
		ProxyHost: proxy.Host,
		ProxyPort: proxy.Port,
	}, nil
}

func (s *RoutingControllerServer) GetRoute(ctx context.Context, in *pb.RoutingParams) (*pb.RoutingRule, error) {
	// parse out owner and project from url
	u, err := url.ParseRequestURI(in.Url)
	if err != nil {
		log.Printf("Could not parse URL: %s\n", in.Url)
		return nil, status.Error(codes.InvalidArgument, "Could not parse URL")
	}

	components := strings.SplitN(u.EscapedPath()[1:], "/", 3)
	if len(components) < 2 || len(components) > 3 {
		fmt.Printf("Invalid url: %s\n", in.Url)
		return nil, status.Error(codes.InvalidArgument, "Missing owner or project portion of URL")
	}
	if len(components) == 2 {
		components = append(components, "")
	}

	var route Route
	tx := db.Preload("Proxy").Where(&Route{
		Owner:   components[0],
		Project: components[1],
	}).Limit(1).First(&route)
	if tx.Error != nil && errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		log.Printf("Could not find route for (owner, project) = (%s, %s)", components[0], components[1])
		return nil, status.Error(codes.Unavailable, "No route exists for owner/project tuple")
	} else if tx.Error != nil {
		log.Print("Database error occurred looking for route")
		return nil, status.Error(codes.Internal, "Error occurred looking up route")
	}

	return &pb.RoutingRule{
		ProxyName: route.Proxy.Name,
		ProxyHost: route.Proxy.Host,
		ProxyPort: route.Proxy.Port,
	}, nil
}

func (s *RoutingControllerServer) CreateEnvironment(ctx context.Context, in *pb.EnvironmentCreateRequest) (*pb.EnvironmentCreateResponse, error) {
	// find a proxy that is available
	var proxy Proxy
	tx := db.First(&proxy)
	if tx.Error != nil {
		log.Print("No proxies available to create environment")
		return nil, status.Error(codes.Unavailable, "No proxy servers available")
	}

	log.Printf("Creating environment for (%s, %s) in %s\n", in.Owner, in.Project, proxy.Name)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", proxy.Host, proxy.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
		log.Print("Failed to connect to proxy")
		return nil, status.Error(codes.Unavailable, "Could not connect to target proxy")
	}
	defer conn.Close()

	c := pb.NewProxyClient(conn)
	cctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	r, err := c.CreateEnvironment(cctx, &pb.EnvironmentCreateRequest{
		Owner:          in.Owner,
		Project:        in.Project,
		ProjectTarball: in.ProjectTarball,
		Requirements:   in.Requirements,
	})
	if err != nil {
		log.Printf("Failed to create environment for (%s, %s) on %s\n", in.Owner, in.Project, proxy.Name)
		return nil, err
	}

	tx = db.Save(&Environment{
		Owner:     in.Owner,
		Project:   in.Project,
		ProxyName: proxy.Name,
	})

	if tx.Error != nil {
		log.Printf("Failed to save environment infos %v\n", tx.Error)
		return nil, status.Error(codes.Internal, "Failed to save environment info to control plane database")
	}

	return r, nil
}

func main() {
	flag.Parse()

	// open database connection
	db_ptr, err := gorm.Open(sqlite.Open(*db_path), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Failed to open database connection")
	}
	db = db_ptr
	db.AutoMigrate(&Proxy{})
	db.AutoMigrate(&Route{})
	db.AutoMigrate(&Environment{})

	// open
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
		panic("Failed to listen to port!")
	}
	s := grpc.NewServer()
	pb.RegisterRoutingControllerServer(s, &RoutingControllerServer{})
	log.Printf("Listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
		panic("Failed to serve GRPC server")
	}
}
