package strongswan

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"strconv"
)

var (
	connectionEstablishedDesc = prometheus.NewDesc(
		"strongswan_connection_established",
		"Connection established (s)",
		[]string{"remoteID"},
		nil)
	connectionBytesInDesc = prometheus.NewDesc(
		"strongswan_connection_bytes_in",
		"Bytes in",
		[]string{"remoteID"},
		nil)
	connectionPacketsInDesc = prometheus.NewDesc(
		"strongswan_connection_packets_in",
		"Packets in",
		[]string{"remoteID"},
		nil)
	connectionBytesOutDesc = prometheus.NewDesc(
		"strongswan_connection_bytes_out",
		"Bytes out",
		[]string{"remoteID"},
		nil)
	connectionPacketsOutDesc = prometheus.NewDesc(
		"strongswan_connection_packets_out",
		"Packets out",
		[]string{"remoteID"},
		nil)
)

type Collector struct {
	client *Client
}

func NewCollector(client *Client) *Collector {
	return &Collector{
		client: client,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- connectionEstablishedDesc
	ch <- connectionBytesInDesc
	ch <- connectionPacketsInDesc
	ch <- connectionBytesOutDesc
	ch <- connectionPacketsOutDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	connections, err := c.client.GetVpnConnections()
	if err != nil {
		log.LogAttrs(context.Background(), slog.LevelError, "error collecting metrics",
			slog.Any("err", err),
		)
		return
	}

	for _, connection := range connections {
		established, err := strconv.ParseFloat(connection.Established, 64)
		if err == nil {
			m, err := prometheus.NewConstMetric(
				connectionEstablishedDesc,
				prometheus.GaugeValue,
				established,
				connection.Remote_id,
			)
			if err != nil {
				log.LogAttrs(context.Background(), slog.LevelError, "error creating metric",
					slog.Any("err", err),
				)
			} else {
				ch <- m
			}
		}

		var bytesIn float64
		var packetsIn float64
		var bytesOut float64
		var packetsOut float64

		if connection.IkeSa.Child_sas != nil {
			for _, childSas := range connection.IkeSa.Child_sas {
				bytesIn += float64(childSas.GetBytesIn())
				packetsIn += float64(childSas.GetPacketsIn())
				bytesOut += float64(childSas.GetBytesOut())
				packetsOut += float64(childSas.GetPacketsOut())
			}
		} else {
			bytesIn = float64(connection.GetBytesIn())
			packetsIn = float64(connection.GetPacketsIn())
			bytesOut = float64(connection.GetBytesOut())
			packetsOut = float64(connection.GetPacketsOut())
		}

		m, err := prometheus.NewConstMetric(
			connectionBytesInDesc,
			prometheus.GaugeValue,
			bytesIn,
			connection.Remote_id,
		)
		if err != nil {
			log.LogAttrs(context.Background(), slog.LevelError, "error creating metric",
				slog.Any("err", err),
			)
		} else {
			ch <- m
		}

		m, err = prometheus.NewConstMetric(
			connectionPacketsInDesc,
			prometheus.GaugeValue,
			packetsIn,
			connection.Remote_id,
		)
		if err != nil {
			log.LogAttrs(context.Background(), slog.LevelError, "error creating metric",
				slog.Any("err", err),
			)
		} else {
			ch <- m
		}

		m, err = prometheus.NewConstMetric(
			connectionBytesOutDesc,
			prometheus.GaugeValue,
			bytesOut,
			connection.Remote_id,
		)
		if err != nil {
			log.LogAttrs(context.Background(), slog.LevelError, "error creating metric",
				slog.Any("err", err),
			)
		} else {
			ch <- m
		}

		m, err = prometheus.NewConstMetric(
			connectionPacketsOutDesc,
			prometheus.GaugeValue,
			packetsOut,
			connection.Remote_id,
		)
		if err != nil {
			log.LogAttrs(context.Background(), slog.LevelError, "error creating metric",
				slog.Any("err", err),
			)
		} else {
			ch <- m
		}
	}
}
