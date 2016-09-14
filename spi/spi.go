// Package spi exposes functions and data for reading temperatures from an MSP430 over SPI.
package spi

import "golang.org/x/exp/io/spi"

// ReadTemperature reads 2 bytes over SPI from /dev/spidev0.1 and returns
// the 16-bit temperature in Celsius as a float64.
func ReadTemperature() (float64, error) {
	dev, errOpen := spi.Open(&spi.Devfs{
		Dev:      "/dev/spidev0.1",
		Mode:     spi.Mode3,
		MaxSpeed: 12000,
	})
	if errOpen != nil {
		return 0, errOpen
	}

	defer dev.Close()

	tx := make([]byte, 2)
	rx := make([]byte, 2)
	errTx := dev.Transfer(tx, rx)
	if errTx != nil {
		return 0, errTx
	}

	tempX10 := (int(rx[0]) << 8) + int(rx[1])

	temp := float64(tempX10) / 10

	return temp, nil
}
