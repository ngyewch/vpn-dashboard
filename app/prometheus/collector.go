package prometheus

import (
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/network_util"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type CustomGatherer struct {
	addressProvider func() ([]string, error)
}

func (gatherer CustomGatherer) Gather() ([]*dto.MetricFamily, error) {
	registry := prometheus.NewRegistry()

	addresses, err := gatherer.addressProvider()
	if err != nil {
		return nil, err
	}

	pingWorker := network_util.NewPingWorker(addresses)
	pingWorker.RunAndWait()

	for _, result := range pingWorker.Results {
		if result.Status == "finished" && result.Statistics != nil {
			if result.Statistics.PacketsRecv > 0 {
				minGauge := prometheus.NewGauge(prometheus.GaugeOpts{
					Name:        "ping_min_ms",
					Help:        "Ping minimum RTT (ms)",
					ConstLabels: prometheus.Labels{"addr": result.Address},
				})
				avgGauge := prometheus.NewGauge(prometheus.GaugeOpts{
					Name:        "ping_avg_ms",
					Help:        "Ping average RTT (ms)",
					ConstLabels: prometheus.Labels{"addr": result.Address},
				})
				maxGauge := prometheus.NewGauge(prometheus.GaugeOpts{
					Name:        "ping_max_ms",
					Help:        "Ping maximum RTT (ms)",
					ConstLabels: prometheus.Labels{"addr": result.Address},
				})
				stdDevGauge := prometheus.NewGauge(prometheus.GaugeOpts{
					Name:        "ping_stddev_ms",
					Help:        "Ping RTT standard deviation (ms)",
					ConstLabels: prometheus.Labels{"addr": result.Address},
				})

				registry.MustRegister(minGauge)
				registry.MustRegister(avgGauge)
				registry.MustRegister(maxGauge)
				registry.MustRegister(stdDevGauge)

				minGauge.Set(float64(result.Statistics.MinRtt.Microseconds()) / 1000)
				avgGauge.Set(float64(result.Statistics.AvgRtt.Microseconds()) / 1000)
				maxGauge.Set(float64(result.Statistics.MaxRtt.Microseconds()) / 1000)
				stdDevGauge.Set(float64(result.Statistics.StdDevRtt.Microseconds()) / 1000)
			}

			packetLossGauge := prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "ping_packet_loss_percent",
				Help:        "Ping packet loss (percent)",
				ConstLabels: prometheus.Labels{"addr": result.Address},
			})

			registry.MustRegister(packetLossGauge)

			packetLossGauge.Set(result.Statistics.PacketLoss)
		}
	}

	collected, err := registry.Gather()
	if err != nil {
		return nil, err
	}

	return collected, nil
}
