package main

import (
	"context"
	"fmt"
	"github.com/ohhfishal/alice-rest/server"
	"os"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

}
