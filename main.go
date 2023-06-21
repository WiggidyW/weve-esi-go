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

func main() {
	env.Init()

	run_local := flag.Bool("local", false, "Run using local cache")

	client := client.NewClient(*run_local)
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
