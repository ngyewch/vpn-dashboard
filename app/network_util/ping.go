package network_util

import (
	"github.com/prometheus-community/pro-bing"
	"time"
)

type PingWorker struct {
	Channels map[string]chan PingResult
	Results  map[string]PingResult `json:"results"`
}

type PingResult struct {
	Address    string              `json:"address"`
	Status     string              `json:"status"`
	Error      string              `json:"error"`
	Statistics *probing.Statistics `json:"statistics"`
}

func NewPingWorker(addrs []string) PingWorker {
	channels := make(map[string]chan PingResult)
	results := make(map[string]PingResult)
	for _, addr := range addrs {
		channels[addr] = make(chan PingResult)
		results[addr] = PingResult{Address: addr, Status: "started", Error: ""}
	}
	return PingWorker{Channels: channels, Results: results}
}

func (worker PingWorker) Run() {
	for addr, result := range worker.Results {
		c := worker.Channels[addr]
		go worker.Ping(result, c)
	}
}

func (worker PingWorker) RunAndWait() {
	worker.Run()

	for addr, channel := range worker.Channels {
		result := <-channel
		worker.Results[addr] = result
	}
}

func (worker PingWorker) Ping(result PingResult, c chan PingResult) {
	pinger, err := probing.NewPinger(result.Address)
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		c <- result
		close(c)
		return
	}

	pinger.Debug = true
	pinger.Count = 5
	pinger.Timeout = 5 * time.Second

	err = pinger.Run()
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		c <- result
		close(c)
		return
	}

	result.Status = "finished"
	result.Error = ""
	result.Statistics = pinger.Statistics()
	c <- result
	close(c)
}
