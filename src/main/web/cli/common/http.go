package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func HttpGet(url string) []byte {
	defer func() {
		if pr := recover(); pr != nil {
			fmt.Printf("panic recover: %v\r\n", pr)
		}
	}()
	log.Printf("url: %s", url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	ErrorHandler(err)
	content, err := ioutil.ReadAll(resp.Body)
	return content
}
