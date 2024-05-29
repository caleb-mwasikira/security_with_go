package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {
	// list all interfaces on current machine
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatalf("failed to list network interfaces: %v", err)
	}

	for _, device := range devices {
		fmt.Println("Device Name: ", device.Name)
		fmt.Println("Device Description: ", device.Description)
		if len(device.Addresses) == 0 {
			println()
		}

		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet Mask: ", address.Netmask)
			println()
		}
	}
}
