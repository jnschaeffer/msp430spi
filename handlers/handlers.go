// Package handlers contains functions and data for handling temperature readings from the MSP430.
package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/jnschaeffer/msp430spi"
)

// Temperature represents a single temperature reading from the MSP430 over SPI.
type Temperature struct {
	DegsC    float64 `json:"tempC"`
	ReadTime string  `json:"readTime"`
	Error    error   `json:"error,omitempty"`
}

// TemperatureHandler represents a single HTTP temperature handler.
type TemperatureHandler struct {
	src     msp430spi.TemperatureSource
	q       chan struct{}
	t       *time.Ticker
	rwMutex sync.RWMutex
	curTemp Temperature
}

// NewTemperatureHandler creates a new TemperatureHandler that reads from the SPI
// device on a periodic basis.
func NewTemperatureHandler(src msp430spi.TemperatureSource, d time.Duration) *TemperatureHandler {

	var h TemperatureHandler
	h.src = src
	h.q = make(chan struct{})
	h.t = time.NewTicker(d)

	go h.run()

	return &h
}

func (h *TemperatureHandler) run() {
	defer h.Close()

	for {
		select {

		// Termination: Close the reader and exit.
		case <-h.q:
			return

		// Temperature read: Get the latest temperature value with each tick.
		case <-h.t.C:
			now := time.Now().UTC()
			degsC, errRead := h.src.DegsC()
			h.rwMutex.Lock()
			// Rather than quit on error, we'll store the error itself,
			// as it could be a transient issue.
			h.curTemp = Temperature{
				DegsC:    degsC,
				ReadTime: now.Format(time.RFC3339),
				Error:    errRead,
			}
			h.rwMutex.Unlock()
		}
	}
}

// Close closes the temperature handler. Subsequent calls will result in a panic.
func (h *TemperatureHandler) Close() error {
	close(h.q)

	return nil
}

// ServeHTTP serves a temperature reading from the SPI device in JSON.
func (h *TemperatureHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	if h.curTemp.Error != nil {
		rw.WriteHeader(500)
	}

	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	enc.Encode(h.curTemp)
}
