package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	mc "microcontrollers"
	"net/http"
)

func (h *Handler) connectHome(c *gin.Context) {
	var input mc.CreateHomeInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	ok := h.services.Home.CreateHome(*input.HomeId, *input.ClientId)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("could not connect home").Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{1})
}

func (h *Handler) getHomeInfo(c *gin.Context) {
	homeId := c.Param("id")

	home, ok := h.services.Home.GetHome(homeId)

	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("could not get home info").Error())
		return
	}

	c.JSON(http.StatusOK, home)
}
func (h *Handler) getHomeInfoTG(c *gin.Context) {
	clientId := c.Param("id")

	home, ok := h.services.Home.GetHomeTG(clientId)

	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("could not get home info").Error())
		return
	}

	c.JSON(http.StatusOK, home)
}

func (h *Handler) updateHomeInfo(c *gin.Context) {
	homeId := c.Param("id")

	var input mc.UpdateHomeInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if ok := h.services.Home.UpdateHome(homeId, input); !ok {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("could not update home info").Error())
		return
	}

	//c.JSON(http.StatusOK, statusResponse{1})

	home, ok := h.services.Home.GetHome(homeId)

	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("could not get home info").Error())
		return
	}
	if *input.IsRobbery {
		baseUrl := `https://api.telegram.org/bot5282040454:AAFJ7IbLtuFtI0EH8TG5JJw6rChnasuUlek/sendMessage`
		api := fmt.Sprintf("%s?chat_id=%s&text=%s", baseUrl, home.ClientId, "Your house is being robbed")
		resp, err := http.Get(api)
		defer resp.Body.Close()

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(resp)
		}

	}
	c.JSON(http.StatusOK, home)
}

func (h *Handler) updateHomeCommandInfo(c *gin.Context) {
	clientId := c.Param("id")

	var input mc.UpdateHomeCommandInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if ok := h.services.Home.UpdateHomeInfo(clientId, input); !ok {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("could not update home security info").Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{1})
}
