package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/el-j/ts2go/pkg/adapters/driving/web"
)

func main() {
	// Parse command line flags
	addr := flag.String("addr", "localhost:8080", "HTTP server address")
	flag.Parse()

	fmt.Println("🌐 TS2Go Web API Server")
	fmt.Println("========================")
	fmt.Println()

	// Create and start the web server
	server, err := web.NewServer(*addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
