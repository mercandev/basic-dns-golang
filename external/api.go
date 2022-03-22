package external

import (
	"dns/model"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func DnsRequestExternalHost(host string) (model.DnsCheckResponse, error) {
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
