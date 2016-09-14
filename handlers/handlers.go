// Package handlers contains functions and data for handling temperature readings from the MSP430.
package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/jnschaeffer/msp430spi/spi"
)

// Temperature represents a single temperature reading from the MSP430 over SPI.
type Temperature struct {
	DegsC    float64 `json:"tempC"`
	ReadTime string  `json:"readTime"`
}

var (
	curTemp Temperature
	errTemp error
	rwMutex sync.RWMutex
)

// TemperatureHandler represents a single HTTP temperature handler.
type TemperatureHandler struct {
	q       chan struct{}
	t       *time.Ticker
	rwMutex sync.RWMutex
	curTemp Temperature
	errTemp error
}

// NewTemperatureHandler creates a new TemperatureHandler that reads from the SPI
// device on a periodic basis.
func NewTemperatureHandler(d time.Duration) *TemperatureHandler {

	var h TemperatureHandler
	h.q = make(chan struct{})
	h.t = time.NewTicker(d)

	go func() {
		for {
			select {
			case <-h.q:
				return
			case <-h.t.C:
				now := time.Now().UTC()
				degsC, errRead := spi.ReadTemperature()
				h.rwMutex.Lock()
				if errRead != nil {
					h.errTemp = errRead
				}
				h.curTemp = Temperature{
					DegsC:    degsC,
					ReadTime: now.Format(time.RFC3339),
				}
				h.rwMutex.Unlock()
			}
		}
	}()

	return &h
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

	if h.errTemp != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(errTemp.Error()))
	}
	temp := h.curTemp

	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	enc.Encode(temp)
}
