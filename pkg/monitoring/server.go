package monitoring

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bcolucci/zerohash/pkg/producer"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartMonitoringServer(points <-chan producer.VWAPUpdateEvent) {
	mux := sync.RWMutex{}

	// trading pairs current VWAP values
	current := map[string]float64{}
	go func() {
		for {
			point, open := <-points
			if !open {
				return
			}
			mux.Lock()
			current[point.ProductId] = point.VWAP
			mux.Unlock()
		}
	}()

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("error while subscribing: %s\n", err)
		}

		// every second, we send the new state of all trading pairs
		t := time.NewTicker(time.Second)
		for range t.C {
			mux.RLock()
			log.Printf("sending %v\n", current)
			b, _ := json.Marshal(current)
			if err := conn.WriteMessage(1, b); err != nil {
				log.Fatalf("error while sending current values: %s\n", err)
			}
			mux.RUnlock()
		}
	})

	// serves static files of the http/ dir
	fs := http.FileServer(http.Dir("./http"))
	http.Handle("/", fs)

	fmt.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
