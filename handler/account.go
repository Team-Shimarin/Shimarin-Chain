package handler

import (
	"net/http"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	conf          *config.Config
	accountAccess dba.AccountAccess
}

func NewAccountHandler(conf *config.Config, accountAccess dba.AccountAccess) *AccountHandler {
	return &AccountHandler{
		conf:          conf,
		accountAccess: accountAccess,
	}
}

func (a *AccountHandler) Register(c *gin.Context) {
	publicKey := c.Param("publickey")

	_, _ = model.NewAccount(publicKey)

	// TODO: RedisのKVSにJSONでぶち込む

	c.String(http.StatusOK, "please wait...")
}
