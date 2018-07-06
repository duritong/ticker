package api

import (
	"net/http"
	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

//GetInitHandler returns the basic settings for the ticker.
func GetInitHandler(c *gin.Context) {
	domain, err := GetDomain(c)

	type settings struct {
		RefreshInterval  int              `json:"refresh_interval,omitempty"`
		InactiveSettings InactiveSettings `json:"inactive_settings,omitempty"`
	}

	s := settings{
		RefreshInterval: 10000,
	}

	ticker, err := FindTicker(domain)
	if err != nil || !ticker.Active {
		s.InactiveSettings = *DefaultInactiveSettings()

		c.JSON(http.StatusOK, JSONResponse{
			Data:   map[string]interface{}{"ticker": nil, "settings": s},
			Status: ResponseSuccess,
			Error:  nil,
		})
		return
	}

	c.JSON(http.StatusOK, JSONResponse{
		//TODO: Build NewTickerPublicResponse to hide unnecessary information
		Data:   map[string]interface{}{"ticker": NewTickerResponse(ticker), "settings": s},
		Status: ResponseSuccess,
		Error:  nil,
	})
	return
}
