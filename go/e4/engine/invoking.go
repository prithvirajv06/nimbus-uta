package engine

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InvokeService(c *gin.Context) {
	var nimId = c.Param("nimb_id")
	//get raw body
	payload, err := c.GetRawData()
	response, err := InvokeServiceByNIMBID(nimId, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()
	c.DataFromReader(response.StatusCode, response.ContentLength, response.Header.Get("Content-Type"), response.Body, nil)
}

func InvokeServiceByNIMBID(nimbID string, payload []byte) (*http.Response, error) {
	endpoint, ok := GetServiceEndpoint(nimbID)
	if !ok {
		return nil, fmt.Errorf("service not found for NIMB_ID: %s", nimbID)
	}
	return http.Post(endpoint, "application/json", bytes.NewReader(payload))
}
