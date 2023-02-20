package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetIndex(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{"name": "lxy"})	
}