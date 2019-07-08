// Package huobi huobi rest api package
package hitbtc

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	//"fmt"
	"fmt"
	"compress/gzip"
	"bytes"
	"io/ioutil"
	"encoding/json"
	//"time"
)

const (
	connecturl = "wss://api.hitbtc.com/api/2/ws"
)

// Huobi API data
type Hitbtc struct {
	Req string
	Id string
}
type tradejson struct {
	Req string `json:"method"`
	Id paramsjson `json:"params"`
}
type paramsjson struct {
	symbol string `json:"symbol"`
}
// New create new Huobi API data
func New(accessKey string, secretKey string) *Huobi {
	return &Okcoin{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (z *Hitbtc) GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

func (z *Hitbtc) GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func (z *Hitbtc) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	c, _, err := websocket.DefaultDialer.Dial(connecturl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	tradereq := tradejson{}
	tradereq.Req = "subscribeTrades"
	tradereq.Id.symbol =  strings.ToUpper(quote) + strings.ToUpper(base)
	tradeReqJson, _ := json.Marshal(tradereq)
	err = c.WriteMessage(websocket.TextMessage,tradeReqJson)
	if err != nil {
		log.Println("write:", err)
		return
	}
	go func() {
		defer close(done)
		for {
			tt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			mm,err :=GzipDecode(message)
			if err != nil {
				fmt.Println("gzip:",err)
			}
			log.Printf("recv: %s", mm)
			log.Printf("mst: %d", tt)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
			/*case  <-ticker.C:
				//err = c.WriteMessage(websocket.TextMessage, []byte("[\"SubAdd\",{\"subs\":[\"0~Huobi~BTC~USD\"]}]"))
				err := c.WriteMessage(websocket.TextMessage, []byte(`{
	  "req": "market.btcusdt.trade.detail",
	  "id": "id11"
	}`))*/
			if err != nil {
				log.Println("write:", err)
				return
			}
			//fmt.Println(t.String())
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
