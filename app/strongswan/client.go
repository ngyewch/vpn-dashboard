package strongswan

import (
	"github.com/bronze1man/goStrongswanVici"
	"strings"
)

type Client struct {
}

func NewClient() (*Client, error) {
	return &Client{}, nil
}

func (client *Client) GetVpnConnections() ([]goStrongswanVici.VpnConnInfo, error) {
	viciClient, err := goStrongswanVici.NewClientConnFromDefaultSocket()
	if err != nil {
		return nil, err
	}
	defer viciClient.Close()

	connections, err := viciClient.ListAllVpnConnInfo()
	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (client *Client) GetAddresses() ([]string, error) {
	connections, err := client.GetVpnConnections()
	if err != nil {
		return nil, err
	}

	var addresses []string
	suffix := "/32"
	for _, connection := range connections {
		for _, subnet := range connection.Remote_ts {
			if strings.HasSuffix(subnet, suffix) {
				addr := subnet[0 : len(subnet)-len(suffix)]
				addresses = append(addresses, addr)
			}
		}
	}

	return addresses, nil
}
