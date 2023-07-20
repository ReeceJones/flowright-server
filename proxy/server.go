package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"

	"shared"
)

type ProxyURIParseError struct {
	Reason string
}

func (err *ProxyURIParseError) Error() string {
	return err.Reason
}

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin:       func(ctx *fasthttp.RequestCtx) bool { return true },
	EnableCompression: true,
}

func main() {
	shared.Init()
	shared.AddRoutingRule("reece", "test", "localhost:8000", false, "")
	shared.AddRoutingRule("reece", "unix_test", "unix", true, "/Users/reece/Documents/Projects/flowright-test/test.socket")
	if err := fasthttp.ListenAndServe("localhost:9000", requestHandler); err != nil {
		fmt.Printf("Error occurred: %v", err)
	}
}

func getProxyPathComponents(uri *fasthttp.URI) (base string, project string, path string, err error) {
	capped_path := string(uri.Path()[1:])
	components := strings.SplitN(capped_path, "/", 3)
	if len(components) < 2 || len(components) > 3 {
		return "", "", "", &ProxyURIParseError{Reason: "Invalid URI"}
	}
	if len(components) == 2 {
		components = append(components, "")
	}
	return components[0], components[1], components[2], nil
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
		return
	}
	route, err := shared.GetRoute(owner, project)
	if err != nil {
		log.Println("Invalid backend:", owner, project)
		log.Println("Cause:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	log.Println("Endpoint:", route)
	uri := fasthttp.AcquireURI()
	err = uri.Parse([]byte(route.Endpoint), []byte(path))
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
		if route.IsUnixSocket {
			dialer = &websocket.Dialer{
				NetDial: func(network, addr string) (net.Conn, error) {
					return net.Dial("unix", route.UnixSocketPath) // TODO, if horizontally scaled this will be issue
				},
			}
		} else {
			dialer = websocket.DefaultDialer
		}

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
	route, err := shared.GetRoute(owner, project)
	if err != nil {
		log.Println("Invalid backend:", owner, project)
		log.Println("Cause:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	log.Println("Endpoint:", route)
	uri := fasthttp.AcquireURI()
	err = uri.Parse([]byte(route.Endpoint), []byte(path))
	if err != nil {
		log.Println("Error parsing request uri:", err)
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	defer fasthttp.ReleaseURI(uri)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(uri.String())
	ctx.Request.Header.VisitAll(func(key []byte, value []byte) {
		req.Header.Set(string(key), string(value))
	})
	resp := fasthttp.AcquireResponse()
	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			if route.IsUnixSocket {
				return net.Dial("unix", route.UnixSocketPath)
			}
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
