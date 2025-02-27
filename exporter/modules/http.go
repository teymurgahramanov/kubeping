package modules

import (
	"crypto/tls"
	"net/http"
	"time"
)

// ProbeHTTP is for probe HTTP endpoints
func ProbeHTTP(address string, timeout int) (bool,error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Duration(timeout)*time.Second,
	}
	resp, err := client.Get(address)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, nil
	} 
	return true, nil
}