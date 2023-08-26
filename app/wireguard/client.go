package wireguard

import "strings"

type Client struct {
}

func NewClient() (*Client, error) {
	return &Client{}, nil
}

func (client *Client) GetVpnConnections() ([]*Peer, error) {
	configPeers, err := ParseWgConfFile("/etc/wireguard/wg0.conf")
	if err != nil {
		return nil, err
	}

	runtimePeers, err := ParseWgShowAllFromCommandOutput("/usr/bin/wg", "show", "all")
	if err != nil {
		return nil, err
	}

	return getVpnConnections(configPeers, runtimePeers)
}

func (client *Client) GetAddresses() ([]string, error) {
	peers, err := client.GetVpnConnections()
	if err != nil {
		return nil, err
	}
	return getAddresses(peers), nil
}

func getVpnConnections(configPeers []*Peer, runtimePeers []*Peer) ([]*Peer, error) {
	configPeerMap := make(map[string]*Peer, 0)
	for _, peer := range configPeers {
		configPeerMap[peer.PublicKey] = peer
	}

	peers := make([]*Peer, 0)
	for _, peer := range runtimePeers {
		if peer.Endpoint != nil {
			configPeer, ok := configPeerMap[peer.PublicKey]
			if ok {
				peer.Name = configPeer.Name
			}
			peers = append(peers, peer)
		}
	}

	return peers, nil
}

func getAddresses(peers []*Peer) []string {
	var addresses []string
	suffix := "/32"
	for _, peer := range peers {
		for _, subnet := range peer.AllowedIps {
			if strings.HasSuffix(subnet, suffix) {
				addr := subnet[0 : len(subnet)-len(suffix)]
				addresses = append(addresses, addr)
			}
		}
	}
	return addresses
}
