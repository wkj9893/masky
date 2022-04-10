package main

import (
	"os"

	"github.com/wkj9893/masky/internal/client"
)

func main() {
	client.Run(client.ParseArgs(os.Args[1:]))
}
