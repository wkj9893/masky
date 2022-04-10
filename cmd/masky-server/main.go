package main

import (
	"os"

	"github.com/wkj9893/masky/internal/server"
)

func main() {
	server.Run(server.ParseArgs(os.Args[1:]))
}
