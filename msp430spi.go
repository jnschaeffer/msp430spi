// Package msp430spi defines functions and data for interacting with an MSP430-powered
// temperature sensor.

package msp430spi

// TemperatureSource represents a source of temperature values.
type TemperatureSource interface {
	DegsC() (float64, error)
	Close() error
}
