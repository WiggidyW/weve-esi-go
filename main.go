package main

import (
	"github.com/rs/cors"

	"flag"
	"fmt"
	"net/http"

	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/env"
	pb "github.com/WiggidyW/weve-esi/proto"
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

	service := client.NewClient(*runLocal)
	server := pb.NewWeveEsiServer(service)

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST"},
		AllowedHeaders: []string{"Content-Type"},
	})
	handler := corsWrapper.Handler(server)

	http.ListenAndServe(env.LISTEN_ADDRESS, handler)
}
