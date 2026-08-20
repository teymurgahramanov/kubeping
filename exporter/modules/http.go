package modules

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// ProbeHTTP is for probe HTTP endpoints. When insecureTLS is true the
// server's certificate is not verified.
func ProbeHTTP(address string, timeout int, insecureTLS bool) (bool, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	resp, err := client.Get(address)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}
	return true, nil
}
