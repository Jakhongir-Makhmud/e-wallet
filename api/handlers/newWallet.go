package handlers

import (
	"e-wallet/storage/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func (h handlers) NewWallet(c *gin.Context) {

	bodyRequest := models.NewWallet{}

	err := c.ShouldBindJSON(&bodyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Err{
			Error: "error while biding json",
		})
		return
	}
	bodyRequest.WalletId = uuid.New().String()

	w, err := h.repo.NewWallet(bodyRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Err{
			Error: err.Error(),
		})
		return
	}

	walletByte, err := json.Marshal(w)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	historyHex := h.auth.HashBody(walletByte)

	c.Writer.Header().Set("X-Digest", historyHex)

	c.JSON(http.StatusOK,w)
}
