package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

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
	// fmt.Println("key1", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
	// Output: key value
	// key2 does not exist}
}

func main() {

	http.HandleFunc("/", ExampleClient)
	//db.Ping(ctx)
	fmt.Println("Listen at 8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
	fmt.Println("Done")
}
