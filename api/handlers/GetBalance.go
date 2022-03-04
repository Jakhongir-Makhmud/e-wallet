package handlers

import (
	"e-wallet/storage/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handlers) GetBalance(c *gin.Context) {

	bodyRequest := models.Wallet{}

	err := c.ShouldBindJSON(&bodyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while binding json",
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


	balance, err := h.repo.GetBalance(bodyRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "something went wrong",
		})
		return
	}

	balanceByte, err := json.Marshal(balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong, please try agan",
		})
		return
	}
	bodyHex := h.auth.HashBody(balanceByte)
	c.Writer.Header().Set("X-Digest", bodyHex)

	c.JSON(http.StatusOK,balance)
}
