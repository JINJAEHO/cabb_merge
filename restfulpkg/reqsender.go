package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func rest() {
	i := 0
	for {
		buff := bytes.NewBuffer([]byte(`
        {
            “timestamp” : ` + fmt.Sprint(i) + `,
            “clientID” : “a”,
            “operation” : “a”
        }
        `))
		http.Post("http://192.168.10.67:10000/req", "application/json", buff)
		fmt.Println(i)
		i++
		time.Sleep(time.Millisecond * 100)
	}
}
