package server

import "github.com/gin-gonic/gin"

func SetServer(r *gin.Engine) {
	r.POST("/search", HandleSearch)
}
