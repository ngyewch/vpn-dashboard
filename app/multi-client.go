package main

import (
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/strongswan"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/wireguard"
)

type MultiClient struct {
	StrongswanClient *strongswan.Client
	WireguardClient  *wireguard.Client
}

func NewMultiClient(strongswanClient *strongswan.Client, wireguardClient *wireguard.Client) (*MultiClient, error) {
	return &MultiClient{
		StrongswanClient: strongswanClient,
		WireguardClient:  wireguardClient,
	}, nil
}

func (client *MultiClient) GetAddresses() ([]string, error) {
	addresses := make([]string, 0)
	if client.StrongswanClient != nil {
		strongswanAddresses, _ := client.StrongswanClient.GetAddresses()
		if strongswanAddresses != nil {
			addresses = append(addresses, strongswanAddresses...)
		}
	}
	if client.WireguardClient != nil {
		wireguardAddresses, _ := client.WireguardClient.GetAddresses()
		if wireguardAddresses != nil {
			addresses = append(addresses, wireguardAddresses...)
		}
	}
	return addresses, nil
}
