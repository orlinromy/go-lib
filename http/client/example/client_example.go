package main

import (
	"fmt"
	"encoding/json"
        "github.com/kelchy/go-lib/http/client"
)

func main() {
	c, _ := client.New()
	r := c.Get("https://www.google.com", nil, nil, 0)
	fmt.Println("SAMPLE HTML HEADER", r.Response)
	html, _ := r.HTML()
	fmt.Println("SAMPLE HTML FIRST 64 CHARS", html[:64])

	// get json payload
	j := c.Get("https://jsonplaceholder.typicode.com/todos/1", nil, nil, 0)
	fmt.Println("SAMPLE JSON HEADER", j.Response)
	js, _ := j.JSON()
	jsons, _ := json.Marshal(js)
	fmt.Println("SAMPLE JSON PAYLOAD", string(jsons))

	// simulate a timeout, but we dont want to see the error stack
	c.SetLogger("empty")
	e := c.Get("https://www.google.com", nil, nil, 10)
	fmt.Println("SAMPLE ERROR TIMEOUT", e.Error)
}
