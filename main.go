package main

import (
	"dns/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
)

func main() {

	//Listen on UDP Port
	addr := net.UDPAddr{
		Port: 90,
		IP:   net.ParseIP("127.0.0.1"),
	}
	u, errUdp := net.ListenUDP("udp", &addr)
	if errUdp != nil {
		panic(errUdp)
	}

	// Wait to get request on that port
	for {
		tmp := make([]byte, 1024)
		_, addr, _ := u.ReadFrom(tmp)
		clientAddr := addr
		packet := gopacket.NewPacket(tmp, layers.LayerTypeDNS, gopacket.Default)
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		tcp, _ := dnsPacket.(*layers.DNS)
		serveDNS(u, clientAddr, tcp)
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	replyMess := request
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var err error
	var ipList []net.IP

	dnsResponse, errDns := callAPI(string(request.Questions[0].Name))

	if errDns != nil {
		panic(errDns)
	}
	if dnsResponse.Status != "success" {
		panic("no dns record")
	}
	fmt.Println(PrettyPrint(dnsResponse))

	for i := 0; i < len(dnsResponse.Data.Answer); i++ {
		ip := dnsResponse.Data.Answer[i]
		a, _, _ := net.ParseCIDR(ip + "/24")
		ipList = append(ipList, a)
	}

	dnsAnswer.IP = ipList[0]
	//a, _, _ := net.ParseCIDR(ip + "/24")
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.Name = []byte(request.Questions[0].Name)
	fmt.Println(request.Questions[0].Name)
	dnsAnswer.Class = layers.DNSClassIN
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

func callAPI(host string) (model.DnsCheckResponse, error) {
	var a model.DnsCheckResponse
	hostName := "https://host-t.com/A/" + host
	req, err := http.NewRequest("GET", hostName, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &a)
	return a, err
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
