// Command srv exposes an HTTP endpoint, /temperature, which provides the current temperature
// in JSON.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jnschaeffer/msp430spi/handlers"
)

func main() {
	var (
		port     string
		interval int
	)
	flag.StringVar(&port, "http", ":8080", "port to listen on")
	flag.IntVar(&interval, "interval", 5, "interval to poll temperatures at in seconds")
	flag.Parse()

	h := handlers.NewTemperatureHandler(time.Duration(interval) * time.Second)
	defer h.Close()
	http.Handle("/temperature", h)

	log.Printf("listening on %s at /temperature", port)
	errHTTP := http.ListenAndServe(port, nil)
	if errHTTP != nil {
		log.Fatal(errHTTP)
	}
}
