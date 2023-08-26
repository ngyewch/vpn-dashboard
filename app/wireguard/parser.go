package wireguard

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ParseWgShowAll(reader io.Reader) ([]*Peer, error) {
	peers := make([]*Peer, 0)
	var peer *Peer = nil
	bufferedReader := bufio.NewReader(reader)
	for {
		line, err := bufferedReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			part := strings.SplitN(line, ": ", 2)
			if len(part) == 2 {
				key := part[0]
				value := part[1]
				if key == "peer" {
					peer = &Peer{
						PublicKey: value,
					}
					peers = append(peers, peer)
				} else if peer != nil {
					if key == "endpoint" {
						peer.Endpoint = &value
					} else if key == "allowed ips" {
						peer.AllowedIps = strings.Split(value, ", ")
					} else if key == "latest handshake" {
						peer.LatestHandshake = &value
					} else if key == "transfer" {
						peer.Transfer = &value
					}
				}
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	return peers, nil
}

func ParseWgShowAllFromCommandOutput(name string, arg ...string) ([]*Peer, error) {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		cmd := exec.Command(name, arg...)
		cmd.Stdout = pipeWriter
		err := cmd.Run()
		if err != nil {
			// do nothing
		}
		err = pipeWriter.Close()
		if err != nil {
			// do nothing
		}
	}()
	return ParseWgShowAll(pipeReader)
}

func ParseWgConf(reader io.Reader) ([]*Peer, error) {
	peers := make([]*Peer, 0)
	var peer *Peer = nil
	bufferedReader := bufio.NewReader(reader)
	for {
		line, err := bufferedReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			if line == "[Peer]" {
				peer = &Peer{}
				peers = append(peers, peer)
			} else {
				part := strings.SplitN(line, "=", 2)
				if len(part) == 2 {
					key := strings.TrimSpace(part[0])
					value := strings.TrimSpace(part[1])
					if key == "peer" {
						peer = &Peer{
							PublicKey: value,
						}
						peers = append(peers, peer)
					} else if peer != nil {
						if key == "PublicKey" {
							peer.PublicKey = value
						} else if key == "# Name" {
							peer.Name = &value
						} else if key == "Endpoint" {
							peer.Endpoint = &value
						} else if key == "AllowedIPs" {
							peer.AllowedIps = strings.Split(value, ", ")
							for i, allowedIp := range peer.AllowedIps {
								peer.AllowedIps[i] = strings.TrimSpace(allowedIp)
							}
						}
					}
				}
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	return peers, nil
}

func ParseWgConfFile(path string) ([]*Peer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	peers, err := ParseWgConf(f)
	if err != nil {
		return nil, err
	}

	return peers, nil
}
