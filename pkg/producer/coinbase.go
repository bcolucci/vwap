package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"golang.org/x/net/websocket"
)

const subscribeUrl = "wss://ws-feed.exchange.coinbase.com"

// CoinbaseUpdateEvent is a Coinbase update event after normalization
// e.g. {ProductId:BTC-USD ChangeType:sell Price:62965.03 Quantity:0.05}
type CoinbaseUpdateEvent struct {
	ProductId  string
	ChangeType string
	Price      float64
	Quantity   float64
}

// ProduceCoinbaseUpdates subscribes to Coinbase real-time updates and produces a channel of Coinbase normalized update events
func ProduceCoinbaseUpdates(ctx context.Context, productIDs []string) (<-chan CoinbaseUpdateEvent, error) {
	wss, err := websocket.Dial(subscribeUrl, "", "http://localhost/")
	if err != nil {
		return nil, fmt.Errorf("error while dialing: %w", err)
	}

	subscribeEvent, _ := json.Marshal(map[string]interface{}{
		"type":        "subscribe",
		"channels":    []string{"level2"},
		"product_ids": productIDs,
	})
	if _, err := wss.Write(subscribeEvent); err != nil {
		return nil, fmt.Errorf("error while subscribing: %w", err)
	}

	ch := make(chan CoinbaseUpdateEvent, 200)
	go func() {
		defer func() {
			log.Println("stop feeding updates")
			wss.Close()
			close(ch)
		}()

		log.Println("start feeding updates")

		var e rawCoinbaseUpdateEvent
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := json.NewDecoder(wss).Decode(&e)
				if err != nil {
					log.Printf("error while decoding message: %s\n", err)
					continue
				}
				if e.Type != "l2update" {
					continue
				}
				log.Printf("received coinbase update %+v\n", e)
				for _, u := range e.ExtractUpdates() {
					ch <- u
				}
			}
		}
	}()

	return ch, nil
}

// raw Coinbase update event
// e.g. {Type:l2update ProductId:BTC-USD Changes:[[buy 62922.10 0.14523324]]}
type rawCoinbaseUpdateEvent struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Changes   [][]string `json:"changes"`
}

// ExtractUpdates maps all the Coinbase raw update event changes to normalized update events
func (e *rawCoinbaseUpdateEvent) ExtractUpdates() []CoinbaseUpdateEvent {
	updates := make([]CoinbaseUpdateEvent, len(e.Changes))
	for i := range e.Changes {
		price, _ := strconv.ParseFloat(e.Changes[i][1], 64)
		quantity, _ := strconv.ParseFloat(e.Changes[i][2], 64)
		updates[i] = CoinbaseUpdateEvent{
			ProductId:  e.ProductId,
			ChangeType: e.Changes[i][0], // buy or sell
			Price:      price,
			Quantity:   quantity,
		}
	}
	return updates
}
