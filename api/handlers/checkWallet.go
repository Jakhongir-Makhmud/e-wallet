package handlers

import (
	"e-wallet/storage/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h handlers) CheckWalletExists(c *gin.Context) {

	body := models.Wallet{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"error in binding json",
		})
		return
	}

	res, err := h.repo.CheckWalletExists(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"something went wrong",
		})
		return
	}
	if res.Balance == 0 && res.Id == "" {
		c.JSON(http.StatusOK,gin.H{
			"message":"such wallet does not exist",
		})
		return
	}

	bodyByte ,err := json.Marshal(res)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"something went wrong, please try agan",
		})
		return
	}
	bodyHex := h.auth.HashBody(string(bodyByte))
	c.Writer.Header().Set("X-Digest",bodyHex)
	
	c.JSON(http.StatusOK,res)
}