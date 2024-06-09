package main

import (
	"AppStore4/daemon"
	"AppStore4/server"
	beatrix "github.com/eXtern-OS/Beatrix"
	"github.com/eXtern-OS/core9-common/db"
	"github.com/eXtern-OS/core9-common/utils"
	"github.com/gin-gonic/gin"
	"log"
)

type Config struct {
	MongoURI         string `json:"mongo_uri"`
	BeatrixToken     string `json:"beatrix_token"`
	BeatrixChannelId string `json:"beatrix_channel_id"`
}

func main() {
	var c Config
	err := utils.ReadConfig(&c)
	if err != nil {
		log.Panicln(err)
	}
	db.Init(c.MongoURI)
	beatrix.Init("AppStore4", c.BeatrixToken, c.BeatrixChannelId)
	go daemon.StartDaemon()
	r := gin.Default()
	server.SetServer(r)
	r.Run()
}
