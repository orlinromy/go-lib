package main

import (
	"time"
	"fmt"
	"os"
	"github.com/kelchy/go-lib/redis"
)

func main() {
	uri := os.Getenv("REDIS_URI")
	redisclient, e := redis.New(uri)
        if e != nil {
		fmt.Println(e)
		return
	}

	keys, _ := redisclient.Keys("*")
	fmt.Println("keys", keys)

	res, _ := redisclient.Set("key", "value2", 10 * time.Second)
	fmt.Println("result", res)

	resi, _ := redisclient.Del("key")
	fmt.Println("result int", resi)

	resb, _ := redisclient.SetNX("key", "value2", 10 * time.Second)
	fmt.Println("result bool", resb)

	val, _ := redisclient.Get("key")
	fmt.Println("result value", val)

	lock, _ := redisclient.Lock("locktest", 20 * time.Second)
	fmt.Println("result lock", lock)

	unlock, _ := redisclient.Unlock("locktest")
	fmt.Println("result unlock", unlock)
}
