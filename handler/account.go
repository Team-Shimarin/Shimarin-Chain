package handler

import (
	"net/http"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"log"
	"fmt"
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
	publicKey := c.PostForm("publicKey")
	log.Print("pubkey", publicKey)
	_, _ = model.NewAccount(publicKey)

	// TODO: RedisのKVSにJSONでぶち込む

	c.String(http.StatusOK, "please wait...")
}

func (a *AccountHandler) UpdateHP(c *gin.Context){
	id := c.Query("id")
	hp, err := strconv.Atoi(c.Query("hp"))

	log.Println(c.Params, "id =", id, " hp = ", hp)

	if  err != nil {
		log.Println("strconv hp parse error ", err)
	}
	log.Println(c.Request.Body, "Params id=", id, " hp=", hp)
	accoutaccsess := dba.AccountAccess{}
	err = accoutaccsess.UpdataBalance(id, int64(hp))
	if err != nil{
		log.Println(err)
	}

	c.String(http.StatusOK, "updated helth point to " + fmt.Sprint(hp))
}
