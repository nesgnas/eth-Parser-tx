package router

import (
	"github.com/gin-gonic/gin"
	"serverETH/controler"
)

func Page(incomingRouter *gin.Engine) {
	incomingRouter.GET("/", controler.CheckServerStatus())
	incomingRouter.GET("/ethTransaction/:address", controler.GetOutTransaction())
	incomingRouter.GET("/currentBlockNumber", controler.GetCurrentBlockNumber())
	incomingRouter.POST("/subscribeToServer", controler.AddSubscriber())
	incomingRouter.PUT("/subscribeToServer", controler.UnsubscribeSubscriber())
}
