package calculate

import (
	"context"
	"etcd-config/pkg/config"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	defaultCalculatePath = "/configs"
	defaultKey           = "/ljw"
	defaultValue         = "/value"
	requestTimeout       = 10 * time.Second
)

var startTime int64
var endTime int64
var count int64
var data = make(chan int64, 10)

func Execute() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	configPath := pwd + defaultCalculatePath
	config.ReadCalculateConfig(configPath, "calculate", "yaml")
	endpoints := strings.Split(config.CalculateConfigMsg.EndPoints, ",")
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: config.CalculateConfigMsg.DialTimeout,
	})
	defer client.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	go updateKV(client)
	if config.CalculateConfigMsg.ChannelCount > 0 {
		data = make(chan int64, config.CalculateConfigMsg.ChannelCount)
	}
	http.HandleFunc("/calculate", calculate)
	err = http.ListenAndServe("0.0.0.0:8888", nil)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func updateKV(client *clientv3.Client) {
	key := defaultKey
	if config.CalculateConfigMsg.DataKey != "" {
		key = config.CalculateConfigMsg.DataKey
	}
	value := defaultValue
	if config.CalculateConfigMsg.DataValue != "" {
		value = config.CalculateConfigMsg.DataValue
	}
	// 目前只更新 2 次，一次创建，一次更新
	for i := 1; i <= 2; i++ {
		ctx, cancel := context.WithTimeout(context.TODO(), requestTimeout)
		startTime = time.Now().UnixNano()
		_, err := client.Put(ctx, key, value)
		if err != nil {
			log.Fatal(err)
		}
		cancel()
		time.Sleep(config.CalculateConfigMsg.UpdateInterval)
		value = value + "1"

	}
}

func calculate(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	t, _ := strconv.ParseInt(values["receiveTime"][0], 10, 64)
	countNub(t)
	fmt.Printf("receive time : %v\n", t)
	w.Write([]byte("calculate finish!"))
}

func countNub(receiveTime int64) {
	if count > 0 {
		return
	}
	select {
	case data <- receiveTime:
		if len(data) == config.CalculateConfigMsg.ChannelCount {
			endTime = time.Now().UnixNano()
			atomic.AddInt64(&count, 1)
			close(data)
			compute()
		}
	default:
		endTime = time.Now().UnixNano()
		atomic.AddInt64(&count, 1)
		close(data)
		compute()
	}

}

func compute() {
	var maxEndTime int64
	fmt.Println("compute......")
	for {
		if d, ok := <-data; ok {
			if d > maxEndTime {
				maxEndTime = d
			}
		} else {
			break
		}
	}
	useTime := (maxEndTime - startTime) / 1000
	fmt.Printf("startTime: %v\n", startTime)
	fmt.Printf("maxEndTime: %v\n", maxEndTime)
	fmt.Printf("endTime: %v\n", endTime)
	fmt.Printf("use time: %v\n", useTime)
	if useTime!=0{
		fmt.Printf("thoughout %v\n", int64(config.CalculateConfigMsg.ChannelCount)*1000/useTime)
	}
}
