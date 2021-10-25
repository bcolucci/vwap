package producer

import (
	"context"
	"log"

	"github.com/bcolucci/zerohash/pkg/util"
)

// VWAPUpdateEvent is an event triggered when we calculate a new VWAP value
type VWAPUpdateEvent struct {
	ProductId string
	NbPoints  int
	VWAP      float64
}

// ProduceVWAPUpdates consumes coinbase update events and produces a channel of VWAP update events
func ProduceVWAPUpdates(ctx context.Context, updates <-chan CoinbaseUpdateEvent, window int) <-chan VWAPUpdateEvent {
	ch := make(chan VWAPUpdateEvent, window)

	go func() {
		defer func() {
			log.Println("stop feeding vwap points")
			close(ch)
		}()

		log.Println("start feeding vwap points")

		queues := map[string]*util.CapQueue{}
		nbPoints := map[string]int{}

		var points []CoinbaseUpdateEvent

		for {
			select {
			case <-ctx.Done():
				return
			case u := <-updates:
				log.Printf("received update %+v\n", u)

				q, f1 := queues[u.ProductId]
				if !f1 {
					q = util.NewCapQueue(window)
					queues[u.ProductId] = q
				}
				q.Append(u)

				nbPoints[u.ProductId]++

				q.CopyTo(&points)
				ch <- VWAPUpdateEvent{
					ProductId: u.ProductId,
					NbPoints:  nbPoints[u.ProductId],
					VWAP:      computeVWAP(points),
				}
			}
		}
	}()

	return ch
}

// computes VWAP value from an array of coinbase update events
func computeVWAP(points []CoinbaseUpdateEvent) float64 {
	if len(points) == 0 {
		return 0.0
	}
	n := 0.0
	d := 0.0
	for _, point := range points {
		n += point.Price * point.Quantity
		d += point.Quantity
	}
	return n / d
}
