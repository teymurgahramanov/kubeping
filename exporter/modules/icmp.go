package modules

import (
	probing "github.com/prometheus-community/pro-bing"
)

func ProbeICMP(address string) (bool, error) {
	pinger, err := probing.NewPinger(address)
	if err != nil {
		return false, err
	}
	pinger.Count = 3
// SET TIMEOUT
	err = pinger.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}