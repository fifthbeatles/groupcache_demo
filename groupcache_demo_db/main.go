package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

type DB struct {
	data map[string]string
}

func NewDB() *DB {
	return &DB{data: make(map[string]string)}
}

func (db *DB) get(key string) (value string, err error) {
	value, ok := db.data[key]
	if !ok {
		err = errors.New("Key '" + key + "' not set yet!")
		return
	}
	time.Sleep(time.Second)
	fmt.Printf("Getting '%s'='%s'\n", key, value)
	return
}

func (db *DB) set(key, value string) {
	db.data[key] = value
	fmt.Printf("Setting '%s'='%s'\n", key, value)
}

func (db *DB) Get(args *string, result *string) (err error) {
	*result, err = db.get(*args)
	return
}

func (db *DB) Set(args *[2]string, result *int) (err error) {
	db.set((*args)[0], (*args)[1])
	return nil
}

func main() {
	db := NewDB()
	rpc.Register(db)
	rpc.HandleHTTP()
	addr := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Serve on", addr)
	http.Serve(listener, nil)
}
