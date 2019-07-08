// Package huobi huobi rest api package
package huobi

import (

	"log"
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
	"github.com/Akagi201/cryptotrader/model"
	"strings"
	//"time"
)

const (
    connecturl = "wss://api.huobi.pro/ws"
)


type tradejson struct {
	Req string `json:"req"`
	Id string `json:"id"`
}
// New create new Huobi API data
func NewWS(accessKey string, secretKey string) *Huobi {
	return &Huobi{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (z *Huobi) GzipEncode(in []byte) ([]byte, error) {
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

func (z *Huobi) GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func (z *Huobi) GetWSTrade(base string, quote string, typ string, since int, size int) ([]model.Trade, error) {
	var trades []model.Trade
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	c, _, err := websocket.DefaultDialer.Dial(connecturl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	tradereq := tradejson{}
	tradereq.Req = "market." + strings.ToLower(quote) + strings.ToLower(base) + ".trade.detail"
	tradereq.Id = "id11"
	tradeReqJson, _ := json.Marshal(tradereq)
	err = c.WriteMessage(websocket.TextMessage,tradeReqJson)
	if err != nil {
		log.Println("write:", err)
		return trades,err
	}
	go func() {
		defer close(done)
		for {
			tt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			mm,err :=z.GzipDecode(message)
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
			return trades,nil
		case  <-ticker.C:
			//err = c.WriteMessage(websocket.TextMessage, []byte("[\"SubAdd\",{\"subs\":[\"0~Huobi~BTC~USD\"]}]"))
			err := c.WriteMessage(websocket.TextMessage, []byte(`{
  "req": "market.btcusdt.trade.detail",
  "id": "id11"
}`))
			if err != nil {
				log.Println("write:", err)
				return trades,nil
			}
			//fmt.Println(t.String())
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return trades,err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return trades,nil
		}
	}
}
