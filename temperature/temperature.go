// Package temperature provides functions and data for getting temperature data
// over SPI.

package temperature

import (
	"github.com/jnschaeffer/msp430spi"
	"github.com/jnschaeffer/msp430spi/spi"
)

// SPISource represents a temperature source that gets temperature
// readings over SPI.
type SPISource struct {
	reader *spi.Reader
}

// NewSPISource creates a new SPI temperature source.
func NewSPISource(path string) (msp430spi.TemperatureSource, error) {
	reader, errReader := spi.NewReader(path)
	if errReader != nil {
		return nil, errReader
	}

	source := &SPISource{
		reader: reader,
	}

	return source, nil
}

// Close releases all resources associate with the temperature source.
func (s *SPISource) Close() error {
	return s.reader.Close()
}

// DegsC gets the current temperature in Celsius over SPI.
func (s *SPISource) DegsC() (float64, error) {
	// Our SPI device sends the temperature out as a two-byte
	// word, most significant byte first, with the value being
	// the current temperature multiplied by a factor of 10.
	rx := make([]byte, 2)

	_, errRead := s.reader.Read(rx)
	if errRead != nil {
		return 0, errRead
	}

	tempX10 := (int(rx[0]) << 8) + int(rx[1])

	temp := float64(tempX10) / 10

	return temp, nil
}
