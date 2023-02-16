package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/kelchy/go-lib/redis"
)

func main() {
	uri := os.Getenv("REDIS_URI")
	clientCert := os.Getenv("REDIS_CLIENT_CERT")
	clientKey := os.Getenv("REDIS_CLIENT_KEY")

	var err error
	var tlsCert tls.Certificate

	if clientCert != "" && clientKey != "" {
		tlsCert, err = tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			fmt.Println("err", err)
		}
	}

	redisclient, e := redis.New(uri, tlsCert)
	if e != nil {
		fmt.Println(e)
		return
	}

	keys, _ := redisclient.Keys(nil, "*")
	fmt.Println("keys", keys)

	res, _ := redisclient.Set(nil, "key", "value2", 10*time.Second)
	fmt.Println("result", res)

	resi, _ := redisclient.Del(nil, "key")
	fmt.Println("result int", resi)

	resb, _ := redisclient.SetNX(nil, "key", "value2", 10*time.Second)
	fmt.Println("result bool", resb)

	val, _ := redisclient.Get(nil, "key")
	fmt.Println("result value", val)

	lock, _ := redisclient.Lock(nil, "locktest", 20*time.Second)
	fmt.Println("result lock", lock)

	unlock, _ := redisclient.Unlock(nil, "locktest")
	fmt.Println("result unlock", unlock)
}
