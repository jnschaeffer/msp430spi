// Package spi exposes functions and data for reading temperatures from an MSP430 over SPI.
package spi

import "golang.org/x/exp/io/spi"

type Reader struct {
	device *spi.Device
}

// NewReader creates a new SPI reader for the given path.
func NewReader(path string) (*Reader, error) {
	dev, errOpen := spi.Open(&spi.Devfs{
		Dev:      "/dev/spidev0.1",
		Mode:     spi.Mode3,
		MaxSpeed: 12000,
	})

	if errOpen != nil {
		return nil, errOpen
	}

	reader := &Reader{
		device: dev,
	}

	return reader, nil
}

// Close releases all values associated with the Reader.
func (r *Reader) Close() error {
	return r.device.Close()
}

// Read reads values from the SPI device into the given byte slice.
func (r *Reader) Read(b []byte) (int, error) {
	tx := make([]byte, len(b))

	errTx := r.device.Tx(tx, b)
	if errTx != nil {
		return 0, errTx
	}

	return len(b), nil
}
