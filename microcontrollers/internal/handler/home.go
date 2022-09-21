package handler

import (
	"errors"
	"fmt"
	"net/http"

	"microcontrollers/internal/entity"
	"microcontrollers/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func (h *Handler) connectHome(ctx *gin.Context) {
	var input entity.CreateHomeInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	_, err := h.service.CreateHome(ctx, *input.HomeId, *input.ClientId)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("could not connect home").Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{1})
}

func (h *Handler) getHomeInfo(ctx *gin.Context) {
	homeId := ctx.Param("id")

	home, err := h.service.GetHome(ctx, homeId)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusInternalServerError, "could not get home info")
		return
	}

	ctx.JSON(http.StatusOK, home)
}
func (h *Handler) getHomeInfoTG(ctx *gin.Context) {
	clientId := ctx.Param("id")

	home, err := h.service.GetHomeTG(ctx, clientId)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(ctx, http.StatusInternalServerError, "could not get home info")
		return
	}

	ctx.JSON(http.StatusOK, home)
}

func (h *Handler) updateHomeInfo(ctx *gin.Context) {
	homeId := ctx.Param("id")

	var input entity.UpdateHomeInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	home, err := h.service.UpdateHome(ctx, homeId, input)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusInternalServerError, "could not update home info")
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
	ctx.JSON(http.StatusOK, home)
}

func (h *Handler) updateHomeCommandInfo(ctx *gin.Context) {
	clientId := ctx.Param("id")

	var input entity.UpdateHomeCommandInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.service.UpdateHomeInfo(ctx, clientId, input)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("could not update home security info").Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{1})
}
