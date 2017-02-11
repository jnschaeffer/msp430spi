// Command print reads from the MSP430 controller at the given rate and prints the temperature in CSV format.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jnschaeffer/msp430spi/temperature"
)

func main() {
	var (
		freq   int
		device string
	)

	flag.IntVar(&freq, "frequency", 1, "frequency of readings in seconds")
	flag.StringVar(&device, "device", "/dev/spidev/0.1", "SPI device to listen on")
	flag.Parse()

	if freq <= 0 {
		flag.PrintDefaults()
		log.Fatal("negative frequency value")
	}

	src, errSource := temperature.NewSPISource(device)
	if errSource != nil {
		log.Fatal(errSource)
	}

	defer src.Close()

	fmt.Println("time,temperature")
	for {
		now := time.Now().UTC()
		temp, err := src.DegsC()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s,%.1f\n", now.Format(time.RFC3339), temp)
		time.Sleep(time.Duration(freq) * time.Second)
	}
}
