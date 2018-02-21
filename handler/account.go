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
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"time"
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
	account, _ := model.NewAccount()

	conf := config.GetConfig()
	var rc redis.Conn
	var err error
	for i := 0; i < 200; i++ {
		rc, err = redis.Dial("tcp", conf.RedisHost + ":" + conf.RedisPort)
		if err != nil {
			log.Printf("%s:%s", conf.RedisHost, conf.RedisPort)
			log.Printf(err.Error())
			log.Printf("redis connection: retry cnt %d", i)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	defer rc.Close()
	// RedisのKVSにJSONでぶち込む
	// {
	//      "id":id,
	//      "hp":hp,
	//}

	data, err := model.NewAccount()
	if err != nil{
		log.Println(err)
	}
	datajson, err := json.Marshal(data)
	if err != nil{
		log.Print(err)
	}
	rc.Do("SET", account.ID, datajson)
	log.Println("set public key", account.ID, " ", fmt.Sprint(datajson))
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

func (a *AccountHandler) GetHP(c *gin.Context) {
	log.Println(c.Request.URL, " posted ", c.Query("id"))
	id := c.Query("id")
	accountaccess := dba.AccountAccess{}
	healthmodel, err := accountaccess.GetHealth(id)
	if err != nil{
		log.Println(err)
		c.String(http.StatusBadRequest,  fmt.Sprintln(err))
	}
	c.String(http.StatusOK, fmt.Sprint(healthmodel.Hp))
}

func (a *AccountHandler) GetBalance(c *gin.Context) {
	id := c.Query("id")
	accountaccess := dba.AccountAccess{}
	balance, err := accountaccess.GetBalance(id)
	if err != nil{
		log.Println(err)
	}
	c.String(http.StatusOK, fmt.Sprint(balance))
}

func (a *AccountHandler) GetBlock(c *gin.Context){
	blockaccsess := dba.BlockAccess{}
	block, err := blockaccsess.GetAllBlock()
	if err != nil{
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintln(err))
	}
	c.String(http.StatusOK, fmt.Sprint(block))
}