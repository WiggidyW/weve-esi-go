package main

import (
	"google.golang.org/grpc"

	"flag"
	"fmt"
	"net"

	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/env"
	"github.com/WiggidyW/weve-esi/proto"
)

var runLocal = flag.Bool("local", false, "Run using local cache")

func main() {
	flag.Parse()
	env.Init()

	if *runLocal {
		fmt.Println("Starting server in local mode")
	} else {
		fmt.Println("Starting server")
	}

	client := client.NewClient(*runLocal)
	server := grpc.NewServer()

	proto.RegisterWeveEsiServer(server, client)

	listener, err := net.Listen("tcp", env.LISTEN_ADDRESS)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to listen on %s: %e",
			env.LISTEN_ADDRESS,
			err,
		))
	}

	server.Serve(listener)
}
