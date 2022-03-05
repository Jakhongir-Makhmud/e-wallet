package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"e-wallet/config"
	"e-wallet/storage/repo"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	Cfg  config.Config
	Repo repo.Repo
}

func (auth Auth) Auth(c *gin.Context) {

	if auth.isPermitted(c.Request.URL.Path) {
		return
	}
	if strings.Contains(c.Request.URL.Path,"/swagger/") {
		return
	}
	hashSum := c.GetHeader("X-Digest")
	userId := c.GetHeader("X-UserId")

	// length of hmac-sha1 is 20 bytes, not less not more
	if len([]byte(hashSum)) != 20 {
		c.AbortWithStatus(401)
		return
	}

	if userId == "" {
		c.AbortWithStatus(401)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	h := hmac.New(sha1.New, nil)

	h.Write(body)

	hashedBody := hex.EncodeToString(h.Sum(nil))

	if hashedBody != hashSum {
		c.AbortWithStatus(401)
	}

	isExists, err := auth.Repo.CheckUserById(userId)

	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	if !isExists {
		c.AbortWithStatus(401)
		return
	}
}

func (auth Auth) HashBody(body []byte) string {

	h := hmac.New(sha1.New, nil)

	_, err := h.Write(body)

	if err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

func (a Auth) isPermitted(path string) bool {

	for _, v := range urls {
		if path == v {
			return true
		}
	}

	return false

}
