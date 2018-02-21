package handler

import (
	"net/http"

	"fmt"
	"log"
	"strconv"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	anzuredis "github.com/InvincibleMan/anzu-chain/redis"
	"github.com/InvincibleMan/anzu-chain/tx"
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

func (a *AccountHandler) UpdateHP(c *gin.Context) {
	hp, err := strconv.Atoi(c.Query("hp"))

	log.Println(c.Params, " hp = ", hp)

	if err != nil {
		log.Println("strconv hp parse error ", err)
	}
	log.Println(c.Request.Body, "Params hp=", hp)
	accoutaccsess := dba.AccountAccess{}
	err = accoutaccsess.UpdateBalance(config.GetConfig().MinorAccountID, int64(hp))
	if err != nil {
		log.Println(err)
	}

	c.String(http.StatusOK, "updated helth point to "+fmt.Sprint(hp))
}

func (a *AccountHandler) Remit(c *gin.Context) {
	toid := c.Query("to")
	fromid := c.Query("from")
	value, err := strconv.Atoi(c.Query("value"))
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprint(err))
		return
	} else if toid == "" || fromid == "" || value == 0 {
		c.String(http.StatusBadRequest, "need to fill toid fromid value")
		return
	}

	tx := tx.Tx{
		toid,
		fromid,
		int64(value),
	}

	if err := anzuredis.AddSetToTxPoolKey(tx); err != nil {
		log.Println(err)
	}

	c.String(http.StatusOK, "remit finished now try to calculate")
}
