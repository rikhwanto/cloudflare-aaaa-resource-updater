package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ipAddress struct {
	Address string `json:"ip"`
}

type dnsRecord struct {
	ResourceName string `json:"name"`
	ResourceType string `json:"type"`
	Address      string `json:"content"`
	ProxyStatus  bool   `json:"proxied"`
}

func getMachineAdress() ipAddress {
	requestURL := "https://api6.ipify.org?format=json"
	res, getErr := http.Get(requestURL)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	ipResponse := ipAddress{}

	jsonErr := json.Unmarshal(body, &ipResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return ipResponse
}

func setDnsRecord(machineAdress ipAddress) {

	zoneId := os.Getenv("ZONE_ID")
	dnsRecordId := os.Getenv("DNS_RECORD_ID")
	apiKey := os.Getenv("CLOUDFLARE_API_KEY")
	resourceName := os.Getenv("RESOURCE_NAME")

	apiUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneId, dnsRecordId)
	apiHeader := fmt.Sprintf("Bearer %s", apiKey)

	content := dnsRecord{
		ResourceName: resourceName,
		ResourceType: "AAAA",
		Address:      machineAdress.Address,
		ProxyStatus:  true,
	}
	body, jsonErr := json.Marshal(content)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Println(string(body))

	req, requestErr := http.NewRequest(http.MethodPatch, apiUrl, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", apiHeader)

	if requestErr != nil {
		log.Fatal(requestErr)
	}

	client := &http.Client{}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		log.Fatal(clientErr)
	}

	defer resp.Body.Close()

	body, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		log.Fatal(ioErr)
	}
	log.Println(string(body))
}

func main() {
	setDnsRecord(getMachineAdress())
}
