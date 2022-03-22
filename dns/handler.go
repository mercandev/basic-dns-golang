package dns

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const LISTEN_IP = "127.0.0.1"
const DNS_PORT = 53

func ListenAndHandle() {
	addr := net.UDPAddr{
		Port: DNS_PORT, //dns port
		IP:   net.ParseIP(LISTEN_IP),
	}

	u, errUdp := net.ListenUDP("udp", &addr)

	if errUdp != nil {
		panic(errUdp)
	}

	for {
		tmp := make([]byte, 1024)
		_, addr, _ := u.ReadFrom(tmp)
		clientAddr := addr
		packet := gopacket.NewPacket(tmp, layers.LayerTypeDNS, gopacket.Default)
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		tcp, _ := dnsPacket.(*layers.DNS)
		ServeDnsRequest(u, clientAddr, tcp)
	}
}
