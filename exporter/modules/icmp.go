package modules

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func ProbeICMP(address string, timeout int) (bool, error) {
	if timeout <= 0 {
		timeout = 5
	}
	pinger, err := probing.NewPinger(address)
	if err != nil {
		return false, err
	}
	pinger.Count = 3
	pinger.Timeout = time.Duration(timeout) * time.Second
	err = pinger.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}