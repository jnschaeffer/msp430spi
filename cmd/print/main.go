// Command print reads from the MSP430 controller at the given rate and prints the temperature in CSV format.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jnschaeffer/msp430spi/spi"
)

func main() {
	var freq int

	flag.IntVar(&freq, "frequency", 1, "frequency of readings in seconds")
	flag.Parse()

	if freq <= 0 {
		flag.PrintDefaults()
		log.Fatal("negative frequency value")
	}

	fmt.Println("time,temperature")
	for {
		now := time.Now().UTC()
		temp, err := spi.ReadTemperature()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s,%.1f\n", now.Format(time.RFC3339), temp)
		time.Sleep(time.Duration(freq) * time.Second)
	}
}
