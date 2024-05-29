package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

var (
	device      = "lo"
	snapshotLen = 1024
	promiscuous = false
	err         error
	timeout     = 30 * time.Second
	handle      *pcap.Handle
	packetCount = 0
)

// list all interfaces on current machine
func listNetworkDevice(name string) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatalf("failed to list network interfaces: %v", err)
	}

	for _, device := range devices {
		if device.Name == name {
			fmt.Printf("%v [%v] (%v)\n", device.Name, device.Flags, device.Description)

			for _, address := range device.Addresses {
				fmt.Println("- IP address: ", address.IP)
				fmt.Println("- Subnet Mask: ", address.Netmask)
			}

			println()
		}
	}
}

func main() {
	listNetworkDevice(device)

	// open device in live capture mode
	handle, err = pcap.OpenLive(device, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		log.Fatalf("failed to open network interface %v; %v", device, err)
	}
	defer handle.Close()
	fmt.Printf("sniffing interface %v...\n", device)

	// set filter
	filter := "tcp and port 8080"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatalf("failed to set BPF filter on captured packets; %v", err)
	}

	// create new output pcap file and write header
	filename := strings.ReplaceAll(filter, " ", "_")
	filename = path.Join("build", fmt.Sprintf("%v.pcap", filename))
	fmt.Printf("creating pcap file %v...\n", filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create output pcap file; %v", err)
	}
	pcap_writer := pcapgo.NewWriter(file)
	pcap_writer.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
	defer file.Close()

	packetSrc := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSrc.Packets() {
		pcap_writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data())

		// only capture 50 packets and the stop
		if packetCount > 50 {
			break
		}
		packetCount++
	}
}
