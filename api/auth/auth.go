package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"e-wallet/config"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	Cfg config.Config
}

func (auth Auth) Auth(c *gin.Context) {

	hashSum := c.GetHeader("X-Digest")

	if len([]byte(hashSum)) != 20 {
		c.AbortWithStatus(401)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	h := hmac.New(sha1.New,nil)

	h.Write(body)

	hashedBody := hex.EncodeToString(h.Sum(nil))

	if hashedBody != hashSum {
		c.AbortWithStatus(401)
	}

}

func (auth Auth) HashBody(body []byte) string {

	h := hmac.New(sha1.New,nil)

	_,err :=h.Write(body)

	if err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}


