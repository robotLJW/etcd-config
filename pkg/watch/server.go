package watch

import (
	"context"
	"etcd-config/pkg/config"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultWatchPath   = "/configs"
	defaultDialTimeout = 5 * time.Second
)

func Execute() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	configPath := pwd + defaultWatchPath
	config.ReadWatchConfig(configPath, "watch", "yaml")

	endPoints := strings.Split(config.WatchConfigMsg.EndPoints, ",")
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: defaultDialTimeout,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer client.Close()
	count := 0
	watchChan := client.Watch(context.TODO(), config.WatchConfigMsg.WatchKey)
	for response := range watchChan {
		for _, ev := range response.Events {
			receiveTime := time.Now().UnixNano()
			count++
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if count == 2 {
				url := fmt.Sprintf("%s?receiveTime=%v", config.WatchConfigMsg.PostAddress, receiveTime)
				fmt.Println(url)
				http.Post(url, "application/json", nil)
			}
		}
	}
	return nil
}
