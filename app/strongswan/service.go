package strongswan

import (
	"github.com/bronze1man/goStrongswanVici"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service struct {
	client *Client
}

type connectionsResponse struct {
	Entries []goStrongswanVici.VpnConnInfo `json:"entries"`
}

func NewService(client *Client) (*Service, error) {
	return &Service{
		client: client,
	}, nil
}

func (service *Service) Install(e *echo.Echo) {
	e.GET("/service/strongswan/connections", service.connections)
}

func (service *Service) connections(c echo.Context) error {
	connections, err := service.client.GetVpnConnections()
	if err != nil {
		return err
	}

	response := connectionsResponse{Entries: connections}
	return c.JSON(http.StatusOK, response)
}
