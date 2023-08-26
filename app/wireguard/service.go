package wireguard

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service struct {
	client *Client
}

type connectionsResponse struct {
	Entries []*Peer `json:"entries"`
}

func NewService(client *Client) (*Service, error) {
	return &Service{
		client: client,
	}, nil
}

func (service *Service) Install(r *gin.Engine) {
	r.GET("/service/wireguard/connections", service.connections)
}

func (service *Service) connections(c *gin.Context) {
	peers, err := service.client.GetVpnConnections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := connectionsResponse{Entries: peers}

	c.JSON(200, response)
}
