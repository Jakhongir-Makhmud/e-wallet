package handlers

import (
	"e-wallet/storage/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


// @Summary Fill Wallet
// @By this endpoin you can fill or top-up your wallet
// @Security Digest
// @Accept json
// @Produce json
// @Param body body models.WalletFill true "fill wallet"
// @Success 200 {object} models.Wallet
// @Failure 401 {object} models.Err
// @Failure 500 {object} models.Err
// @Router /wallet/fill [post]
func (h handlers) FillWallet(c *gin.Context) {

	bodyRequest := models.WalletFill{}

	err := c.ShouldBindJSON(&bodyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while binding json",
		})
		return
	}

	w, err := h.repo.CheckWalletExists(models.Wallet{Id: bodyRequest.Id})

	if w.Id == "" || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no such wallet, please check it",
		})
		return
	}

	w, err = h.repo.FillWallet(bodyRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
	}

	fillByte, err := json.Marshal(w)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong, please try agan",
		})
		return
	}
	bodyHex := h.auth.HashBody(fillByte)
	c.Writer.Header().Set("X-Digest", bodyHex)

	c.JSON(http.StatusOK,w)
}
