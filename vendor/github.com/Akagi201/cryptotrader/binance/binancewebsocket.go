// Package huobi huobi rest api package
package binance

import (

	"log"

	"time"
	"os"
	"os/signal"
	"github.com/gorilla/websocket"
	"github.com/Akagi201/cryptotrader/model"
	"strings"
)

const (
	connecturl = "wss://stream.binance.com:9443/ws/"
)

// Huobi API data
type Binance struct {
	Req string
	Id string
}
type tradejson struct {
	Req string `json:"req"`
	Id string `json:"id"`
}
// New create new Huobi API data
func NewBinance(accessKey string, secretKey string) *Binance {
	return &Binance{
		Req: accessKey,
		Id: secretKey,
	}
}

// GetTicker 行情
func (z *Binance) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Trade, error) {
	connecturl1 := connecturl + strings.ToLower(quote) + strings.ToLower(base) + "@aggTrade"
	var trades []model.Trade
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	c, _, err := websocket.DefaultDialer.Dial(connecturl1, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			tt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", string(message))
			log.Printf("mst: %d", tt)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return trades,nil
		/*case  <-ticker.C:
			//err = c.WriteMessage(websocket.TextMessage, []byte("[\"SubAdd\",{\"subs\":[\"0~Huobi~BTC~USD\"]}]"))
			err := c.WriteMessage(websocket.TextMessage, []byte(`{
  "req": "market.btcusdt.trade.detail",
  "id": "id11"
}`))*/
			if err != nil {
				log.Println("write:", err)
				return trades ,err
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
