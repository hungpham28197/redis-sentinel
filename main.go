package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ExampleClient(w http.ResponseWriter, r *http.Request) {
	rdb := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName: "mymaster",
		SentinelAddrs: []string{
			os.Getenv("redis-sentinel-1") + ":26379",
			os.Getenv("redis-sentinel-2") + ":26379",
			os.Getenv("redis-sentinel-3") + ":26379",
		},
		Password:         "123",
		SentinelPassword: "123",
		RouteRandomly:    true,
		ReadTimeout:      10 * time.Second,
		WriteTimeout:     10 * time.Second,
		PoolTimeout:      10 * time.Second,
	})
	defer rdb.Close()

	rdb.Ping(ctx)

	err := rdb.Set(ctx, "key1", rand.Int(), 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key1").Result()
	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintf(w, val)
}

func ExampleRedisCluster(w http.ResponseWriter, r *http.Request) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			os.Getenv("redis-node-0") + ":6379",
			os.Getenv("redis-node-1") + ":6379",
			os.Getenv("redis-node-2") + ":6379",
			os.Getenv("redis-node-3") + ":6379",
			os.Getenv("redis-node-4") + ":6379",
			os.Getenv("redis-node-5") + ":6379",
		},
		Password:      "bitnami",
		RouteRandomly: true,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		PoolTimeout:   10 * time.Second,
	})
	defer rdb.Close()

	rdb.Ping(ctx)

	for i := 0; i < 100; i++ {
		err := rdb.Set(ctx, "redis-cluster-"+strconv.Itoa(i), "redis-cluster-test", 0).Err()
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 100; i++ {
		val, err := rdb.Get(ctx, "redis-cluster-"+strconv.Itoa(i)).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Key: " + "redis-cluster-" + strconv.Itoa(i) + "   Value: " + val)
	}

	_, _ = fmt.Fprintf(w, "Done")
}

func main() {

	http.HandleFunc("/", ExampleRedisCluster)
	//db.Ping(ctx)
	//http.HandleFunc("/", ExampleDockerSecret)
	fmt.Println("Listen at 8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
	fmt.Println("Done")
}
