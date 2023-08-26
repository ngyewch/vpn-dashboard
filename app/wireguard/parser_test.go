package wireguard

import (
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/utils"
	"testing"
)

func parseWgShowAllTest() ([]*Peer, error) {
	return ParseWgShowAllFromCommandOutput("/bin/cat", "testdata/wg.txt")
}

func parseWgConfTest() ([]*Peer, error) {
	return ParseWgConfFile("testdata/wireguard.conf")
}

func dumpPeers(peers []*Peer) error {
	for _, peer := range peers {
		json, err := utils.ToJson(*peer)
		if err != nil {
			return err
		}
		println(json)
	}
	return nil
}

func Test1(t *testing.T) {
	configPeers, err := parseWgConfTest()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	runtimePeers, err := parseWgShowAllTest()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	peers, err := getVpnConnections(configPeers, runtimePeers)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = dumpPeers(peers)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	println("-----")

	addresses := getAddresses(peers)
	for _, address := range addresses {
		println(address)
	}
}
