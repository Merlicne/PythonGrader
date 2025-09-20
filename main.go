package main

import (
	"context"
	"fmt"
	"os"
)

func main() {

	app := createCommand()
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	_ = app
}
