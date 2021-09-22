package main

import (
	"etcd-config/pkg/watch"
	"os"
)

func main() {
	if err := watch.Execute(); err != nil {
		os.Exit(1)
	}
}
