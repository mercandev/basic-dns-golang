package dns

import (
	"dns/external"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const DEFAULT_DNS_TTL uint32 = 128
const IP_SUBNET string = "/24"

func ServeDnsRequest(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	var dnsAnswer layers.DNSResourceRecord
	var err error

	dnsResponse, errDns := external.DnsRequestExternalHost(string(request.Questions[0].Name)) //check host dns address
	//fmt.Println(util.PrettyPrint(dnsResponse))

	if errDns != nil {
		panic(errDns)
	}

	if dnsResponse.Status != "success" {
		panic("no dns record")
	}

	ip := dnsResponse.Data.Answer[0]
	a, _, _ := net.ParseCIDR(ip + IP_SUBNET)

	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.TTL = DEFAULT_DNS_TTL
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.Name = []byte(request.Questions[0].Name)
	dnsAnswer.Class = layers.DNSClassIN

	replyMess := request
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err = replyMess.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	u.WriteTo(buf.Bytes(), clientAddr)
}
