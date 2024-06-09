package server

import (
	"AppStore4/daemon"
	"context"
	"fmt"
	beatrix "github.com/eXtern-OS/Beatrix"
	"github.com/eXtern-OS/core9-common/db"
	"github.com/eXtern-OS/core9-common/models/app"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
)

type searchQuery struct {
	Query string `json:"query"`
}

func HandleSearch(c *gin.Context) {
	var q searchQuery
	if err := c.BindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	filter := bson.D{{"name", bson.D{{"$regex", strings.ToLower(q.Query)}}}}
	filter2 := bson.D{{"name", bson.D{{"$regex", strings.ToUpper(string(q.Query[0])) + q.Query[1:]}}}}
	daemon.D.CheckRunning()
	cur, err := db.DefaultClient.FindMany(filter, "AppStore", "Snaps")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var res []app.App
	for cur.Next(context.Background()) {
		var a app.Snap
		if err := cur.Decode(&a); err != nil {
			go beatrix.SendError(fmt.Sprintf("Failed to unpack snap app from database %e", err.Error()), "AppStore4.HandleSearch")
			continue
		}
		res = append(res, a.Export())
	}
	cur, err = db.DefaultClient.FindMany(filter2, "AppStore", "Flatpaks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	for cur.Next(context.Background()) {
		var a app.Flatpak
		if err := cur.Decode(&a); err != nil {
			go beatrix.SendError(fmt.Sprintf("Failed to unpack flatpak app from database %e", err.Error()), "AppStore4.HandleSearch")
			continue
		}
		res = append(res, a.Export())
	}
	c.JSON(http.StatusOK, gin.H{"results": res})
}
