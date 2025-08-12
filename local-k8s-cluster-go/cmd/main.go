package main

import (
	"os"

	"github.com/romdj/local-k8s-cluster-go/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
