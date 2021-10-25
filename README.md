# VWAP

This project is a real-time VWAP (volume-weighted average price) calculation engine.

## How does it works

The project basically use channels to pipe trading pair update events from a Coinbase websocket to a VWAP calculator.

The steps are:

1. Subscribe to a Coinbase websocket to be notified when there is an update on our trading pairs (buy or sell) ;
2. Normalize/extract data from these events and return a channel of clean update events ;
3. Listen to the normalized channel of updates and return a new channel of VWAP values using a sliding window of 200 values.
4. Print the values on **console** mode or serve them (WebSocket) to a webapp on **serve** mode.

## Trading pairs

Currently, the tradings pairs are hardcoded to:
* BTC-USD
* ETH-USD
* ETH-BTC

## Dependencies

No external dependencies have been used in order to create the calculator.

For the tests, the project uses [testify](github.com/stretchr/testify).
For the monitoring webapp, the project uses [Gorilla WebSocket](github.com/gorilla/websocket) and [chart.js](https://www.chartjs.org/).

## Running modes

There are two running modes:

* **console**: (default) Prints current VWAP values every time we receive an update.
* **server**: Serve an HTML monitoring page with VWAP line charts.

### Console mode

Simply run: `make run`

You should see something like: 

```
[...]
ETH-USD: 4232.7580      (1303 pts)
ETH-USD: 4232.7588      (1304 pts)
BTC-USD: 63639.8765     (982 pts)
BTC-USD: 63639.8531     (983 pts)
BTC-USD: 63639.8495     (984 pts)
ETH-USD: 4232.7365      (1305 pts)
BTC-USD: 63639.8423     (985 pts)
BTC-USD: 63639.8330     (986 pts)
```

**Currently, by default, we stop the console run after 30 seconds because it's enought time to see that it works**

### Server mode

In this mode, the application will serve a monitoring page on http://localhost:8080.

Simply run: `make serve`

You should see something like:

![charts screenshot](./charts.gif)

### Test

`make test`

### Build

`make` or `make build`

```bash
> bin/vwap                    
ETH-USD: 4219.2200      (2 pts)
ETH-USD: 4218.1603      (3 pts)
ETH-USD: 4221.1899      (4 pts)
ETH-USD: 4221.1698      (5 pts)
BTC-USD: 63436.5500     (1 pts)
BTC-USD: 63487.6358     (2 pts)
```

## Choices

* I decided to keep things at simple as possible (Go mindset). 
* Even if 100% of the code is not covered, I decided to cover the most important parts.
* I added a very little webapp so we can visualize the three charts in real-time.

I would say that I keep a kind of "lean" mindset. So I started to do something that works, then I extract some code so it could be tested, and I finally a some visualization tool.
