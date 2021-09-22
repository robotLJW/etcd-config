package main

import (
	"etcd-config/pkg/calculate"
	"os"
)

func main() {
	if err := calculate.Execute(); err != nil {
		os.Exit(1)
	}
}
