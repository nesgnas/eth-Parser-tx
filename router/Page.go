package router

import (
	"github.com/gin-gonic/gin"
	"serverETH/controler"
)

func Page(incomingRouter *gin.Engine) {
	incomingRouter.GET("/", controler.CheckServerStatus())
	incomingRouter.GET("/transactions/:address", controler.GetOutTransaction())
	incomingRouter.GET("/blockNumber", controler.GetCurrentBlockNumber())
	incomingRouter.POST("/subscriptions", controler.AddSubscriber())
	incomingRouter.PUT("/subscriptions", controler.UnsubscribeSubscriber())
}
