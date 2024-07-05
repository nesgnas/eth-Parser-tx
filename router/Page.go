package router

import (
	"github.com/gin-gonic/gin"
	"serverETH/controler"
)

func Page(incomingRouter *gin.Engine) {
	incomingRouter.GET("/eth", controler.Gathering())
	incomingRouter.GET("/ethTransaction", controler.GetOutTransaction())
}
