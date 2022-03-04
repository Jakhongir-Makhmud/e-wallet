package handlers

import (
	"e-wallet/storage/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h handlers) GetHistory(c *gin.Context) {

	bodyRequest := models.Wallet{}

	err := c.ShouldBindJSON(&bodyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"error while binding json",
		})
		return
	}

	w,err := h.repo.CheckWalletExists(bodyRequest)

	if w.Id == "" || err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"no such wallet",
		})
		return
	}


	history,err := h.repo.GetHistory(bodyRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"something went wrong, please try again",
		})
		return
	}

	historyByte,err := json.Marshal(history)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"something went wrong,please try again",
		})
		return
	}

	historyHex := h.auth.HashBody(historyByte)

	c.Writer.Header().Set("X-Digest",historyHex)

	c.JSON(http.StatusOK,history)
}
