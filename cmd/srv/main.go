// Command srv exposes an HTTP endpoint, /temperature, which provides the current temperature
// in JSON.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jnschaeffer/msp430spi/handlers"
	"github.com/jnschaeffer/msp430spi/temperature"
)

func main() {
	var (
		port     string
		interval int
		device   string
	)
	flag.StringVar(&port, "http", ":8080", "port to listen on")
	flag.IntVar(&interval, "interval", 5, "interval to poll temperatures at in seconds")
	flag.StringVar(&device, "device", "/dev/spidev/0.1", "SPI device to listen on")

	flag.Parse()

	src, errSource := temperature.NewSPISource(device)
	if errSource != nil {
		log.Fatal(errSource)
	}

	h := handlers.NewTemperatureHandler(src, time.Duration(interval)*time.Second)
	defer h.Close()

	http.Handle("/temperature", h)

	log.Printf("listening on %s at /temperature", port)
	errHTTP := http.ListenAndServe(port, nil)
	if errHTTP != nil {
		log.Fatal(errHTTP)
	}
}
