package handlers

import (
	"e-wallet/storage/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h handlers) NewUser(c *gin.Context) {

	bodyRequest := models.User{}

	err := c.ShouldBindJSON(&bodyRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.Err{Error: "error while bidnding json"})
		return
	}
	bodyRequest.Id = uuid.New().String()
	user, err := h.repo.NewUser(bodyRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Err{
			Error: err.Error(),
		})
		return
	}

	userByte, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	historyHex := h.auth.HashBody(userByte)

	c.Writer.Header().Set("X-Digest", historyHex)

	c.JSON(http.StatusOK,user)
}
