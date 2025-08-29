package routes

import (
	"tron-api-go/internal/handlers"
	"tron-api-go/internal/types"
	"tron-api-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// 设置所有路由
func SetupRoutes(r *gin.Engine, config *types.Config) {
	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 中间件
	r.Use(utils.CorsMiddleware())

	// 创建处理器服务
	handlerService := handlers.NewService(config)

	// 主页路由
	r.GET("/", handlerService.IndexHandler)
	r.GET("/doc", handlerService.DocsHandler)

	// API v1 路由组
	v1 := r.Group("/v1")
	{
		// 工具接口
		v1.Any("/status", handlerService.StatusHandler)
		v1.Any("/getApiList", handlerService.GetApiListHandler)

		// 地址生成相关接口
		v1.Any("/createAddress", handlerService.CreateAddressHandler)
		v1.Any("/generateAddressWithMnemonic", handlerService.GenerateAddressWithMnemonicHandler)
		v1.Any("/getAddressByKey", handlerService.GetAddressByKeyHandler)
		v1.Any("/mnemonicToAddress", handlerService.MnemonicToAddressHandler)
		v1.Any("/mnemonicToAddressBatch", handlerService.MnemonicToAddressBatchHandler)
		v1.Any("/privateKeyToAddress", handlerService.PrivateKeyToAddressHandler)

		// 余额查询相关接口
		v1.Any("/getTrxBalance", handlerService.GetTrxBalanceHandler)
		v1.Any("/getTrc20Balance", handlerService.GetTrc20BalanceHandler)
		v1.Any("/getTrc10Info", handlerService.GetTrc10InfoHandler)
		v1.Any("/getTrc10Balance", handlerService.GetTrc10BalanceHandler)

		// 转账相关接口
		v1.Any("/sendTrx", handlerService.SendTrxHandler)
		v1.Any("/sendTrc20", handlerService.SendTrc20Handler)
		v1.Any("/sendTrc10", handlerService.SendTrc10Handler)

		// 交易查询相关接口
		v1.Any("/getTransaction", handlerService.GetTransactionHandler)
		v1.Any("/getTrc20TransactionReceipt", handlerService.GetTrc20TransactionReceiptHandler)

		// 区块链信息查询接口
		v1.Any("/getBlockHeight", handlerService.GetBlockHeightHandler)
		v1.Any("/getBlockByNumber", handlerService.GetBlockByNumberHandler)
	}
}
