package wireguard

type Peer struct {
	PublicKey       string   `json:"publicKey"`
	AllowedIps      []string `json:"allowedIps"`
	Name            *string  `json:"name"`
	Endpoint        *string  `json:"endpoint"`
	LatestHandshake *string  `json:"latestHandshake"`
	Transfer        *string  `json:"transfer"`
}
