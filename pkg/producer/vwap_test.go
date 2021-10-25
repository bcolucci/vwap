package producer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestComputeVWAP(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0.0, computeVWAP([]CoinbaseUpdateEvent{}))
	assert.Equal(1.0, computeVWAP([]CoinbaseUpdateEvent{
		{Price: 1, Quantity: 1},
	}))
	assert.Equal(1.0, computeVWAP([]CoinbaseUpdateEvent{
		{Price: 1, Quantity: 1}, // n = 1*1 / q = 1
		{Price: 1, Quantity: 2}, // n = 1+1*2 / q = 1+2
	}))
	assert.Equal(1.25, computeVWAP([]CoinbaseUpdateEvent{
		{Price: 1, Quantity: 1}, // n = 1*1 / q = 1
		{Price: 1, Quantity: 2}, // n = 1+1*2 / q = 1+2
		{Price: 2, Quantity: 1}, // n = 3+2*1 / q = 3+1
	}))
}

func TestProduceVWAPUpdates(t *testing.T) {
	assert := assert.New(t)

	coinbase := make(chan CoinbaseUpdateEvent, 1)

	ctx, cancel := context.WithCancel(context.Background())
	updates := ProduceVWAPUpdates(ctx, coinbase, 2)

	values := []float64{}
	go func() {
		for v := range updates {
			values = append(values, v.VWAP)
		}
	}()

	coinbase <- CoinbaseUpdateEvent{ProductId: "test", Price: 1, Quantity: 2} // vwap = 1*2/2 = 1
	coinbase <- CoinbaseUpdateEvent{ProductId: "test", Price: 3, Quantity: 1} // vwap = (1*2+3*1)/(2+1) = 1.6666666666666667

	// here the first value is removed because of the window
	coinbase <- CoinbaseUpdateEvent{ProductId: "test", Price: 1, Quantity: 3} // vwap = (3*1+1*3)/(1+3) = 1.5

	time.Sleep(time.Second)
	cancel()

	assert.Equal([]float64{1, 1.6666666666666667, 1.5}, values)
}
