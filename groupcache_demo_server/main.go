package main

import (
	"flag"
	"fmt"
	"github.com/golang/groupcache"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

var (
	peers_addrs = []string{"127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"}
	rpc_addrs   = []string{"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003"}
	index       = flag.Int("index", -1, "peer index")
)

type Server struct {
	group  *groupcache.Group
	client *rpc.Client
}

func (server *Server) GetCache(args *string, result *string) error {
	return server.group.Get(nil, *args, groupcache.StringSink(result))
}

func (server *Server) Get(args *string, result *string) error {
	return server.client.Call("DB.Get", args, result)
}

func (server *Server) Set(args *[2]string, result *int) error {
	return server.client.Call("DB.Set", args, result)
}

func NewServer(group *groupcache.Group, client *rpc.Client) *Server {
	return &Server{group: group, client: client}
}

func main() {
	flag.Parse()
	if *index < 0 || *index >= len(peers_addrs) {
		fmt.Printf("peer_index %d not invalid\n", *index)
		os.Exit(1)
	}

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	peers := groupcache.NewHTTPPool(addrToURL(peers_addrs[*index]))
	peers.Set(addrsToURLs(peers_addrs)...)
	getter := groupcache.GetterFunc(func(ctx groupcache.Context, key string, dest groupcache.Sink) (err error) {
		fmt.Println("Serve key", key)
		var value string
		err = client.Call("DB.Get", &key, &value)
		if err != nil {
			return
		}
		dest.SetString(strconv.Itoa(*index) + ":" + value)
		return nil
	})
	group := groupcache.NewGroup("demo", 1<<20, getter)

	server := NewServer(group, client)
	rpc.Register(server)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", rpc_addrs[*index])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go http.Serve(listener, nil)

	fmt.Println("Start listening on", peers_addrs[*index])
	fmt.Println("Group name:", group.Name())
	log.Fatal(http.ListenAndServe(peers_addrs[*index], peers))
}

func addrToURL(addr string) string {
	return "http://" + addr
}

func addrsToURLs(addrs []string) []string {
	result := make([]string, len(addrs))
	for _, addr := range addrs {
		result = append(result, addrToURL(addr))
	}
	return result
}
