package router

import (
	"net/http"
	"time"

	"smzdtz-server/internal/middleware"
	"smzdtz-server/internal/router/fund"
	"smzdtz-server/internal/router/install"
	"smzdtz-server/internal/router/stock"
	"smzdtz-server/internal/router/user"
	"smzdtz-server/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter() (r *gin.Engine) {
	router := gin.New()

	logger, _ := zap.NewProduction()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(log.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(log.RecoveryWithZap(logger, true))

	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	router.Use(middleware.Cors())
	// 安全，限制接口访问白名单
	// TODO：未来要实现限流
	router.Use(middleware.IPWhiteList())

	baseApi := router.Group("api")
	{
		baseApi.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
	}
	// 系统安装
	installApi := baseApi.Group("install")
	{
		// TODO
		installApi.GET("/status", install.Status)
		installApi.POST("/store", install.Store)
	}
	// 用户
	userApi := baseApi.Group("user")
	{
		userApi.POST("", user.Index)
		userApi.POST("/store", user.Store)
		userApi.POST("/login", user.ValidateLogin)
		userApi.GET("/:id", user.Detail)
	}
	// 股票
	stockApi := baseApi.Group("stock")
	{
		stockApi.GET("/getEMProfile", stock.GetProfile)
		stockApi.GET("/getEMStockNews", stock.GetStockNews)
		stockApi.GET("/getEMZongHePingJia", stock.GetZongHePingJia)
		stockApi.GET("getEMFreeHolderse", stock.GetFreeHolderse)
		stockApi.GET("/getEMIndicator", stock.GetIndicator)
		stockApi.GET("/getEMJiaZhiPingGu", stock.GetJiaZhiPingGu)
		stockApi.GET("/getEMStockTrends", stock.GetStockTrends)

		stockApi.GET("/getZSXGCommentNew", stock.GetCommentNew)
		stockApi.GET("/getCNINFOStockList", stock.GetStockList)

		stockApi.GET("search", stock.SearchStock)
	}
	// 基金
	fundApi := baseApi.Group("fund")
	{
		fundApi.GET("/getEMInfo", fund.GetFundInfo)
	}

	return router
}