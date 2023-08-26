package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func RouterPrometheus(addressProvider func() ([]string, error)) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	gatherer := CustomGatherer{
		addressProvider: addressProvider,
	}
	handlerOpts := promhttp.HandlerOpts{}
	r.Any("/metrics", gin.WrapH(promhttp.HandlerFor(gatherer, handlerOpts)))

	r.GET("/service/ipAddresses", func(c *gin.Context) {
		addresses, err := addressProvider()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"ipAddresses": addresses,
		})
	})

	return r
}
