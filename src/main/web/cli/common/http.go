package common

import (
	"io/ioutil"
	"log"
	"net/http"
)

func HttpGet(url string) []byte {
	log.Printf("url: %s", url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	ErrorHandler(err)
	content, err := ioutil.ReadAll(resp.Body)
	return content
}
