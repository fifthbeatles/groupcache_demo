package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
)

var (
	rpc_addrs = []string{"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003"}
	index     = flag.Int("index", -1, "peer index")
)

func main() {
	flag.Parse()
	if *index < 0 || *index >= len(rpc_addrs) {
		fmt.Printf("peer_index %d not invalid\n", *index)
		os.Exit(1)
	}

	client, err := rpc.DialHTTP("tcp", rpc_addrs[*index])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(9 - i)
		var result int
		err = client.Call("Server.Set", &[2]string{key, value}, &result)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(i)
		var value string
		err = client.Call("Server.GetCache", &key, &value)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Got '%s'='%s'\n", key, value)
	}
	
	for i := 3; i < 6; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(9 - i)
		var result int
		err = client.Call("Server.Set", &[2]string{key, value}, &result)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(i)
		var value string
		err = client.Call("Server.GetCache", &key, &value)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Got '%s'='%s'\n", key, value)
	}
}
