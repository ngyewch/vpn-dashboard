package web_service

import (
	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	"github.com/labstack/echo/v4"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/network_util"
	"net/http"
)

type PingService struct {
	addressProvider func() ([]string, error)
	cache           *lru.Cache
}

type PingResponse struct {
	Id string `json:"id"`
}

type PingResultResponse struct {
	Results map[string]network_util.PingResult `json:"results"`
}

func NewPingService(addressProvider func() ([]string, error)) (*PingService, error) {
	cache, err := lru.New(128)
	if err != nil {
		return nil, err
	}
	return &PingService{
		addressProvider: addressProvider,
		cache:           cache,
	}, nil
}

func (service *PingService) Install(e *echo.Echo) {
	e.GET("/service/ping", service.Ping)
	e.GET("/service/pingResult", service.PingResult)
}

func (service *PingService) Ping(c echo.Context) error {
	addresses, err := service.addressProvider()
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	pingWorker := network_util.NewPingWorker(addresses)
	service.cache.Add(id.String(), pingWorker)
	pingWorker.Run()

	response := PingResponse{Id: id.String()}

	return c.JSON(http.StatusOK, response)
}

func (service *PingService) PingResult(c echo.Context) error {
	id := c.QueryParam("id")
	cached, ok := service.cache.Get(id)
	response := PingResultResponse{}
	if ok {
		worker := cached.(network_util.PingWorker)
		for addr, c := range worker.Channels {
			result, ok := <-c
			if ok {
				worker.Results[addr] = result
			}
		}
		response.Results = worker.Results
	}
	return c.JSON(http.StatusOK, response)
}
