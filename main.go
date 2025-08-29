package main

import (
	"fmt"
	"log"
	"time"

	"tron-api-go/internal/routes"
	"tron-api-go/internal/types"
	"tron-api-go/internal/utils"

	"github.com/gin-gonic/gin"
)

/**
 * TRON区块链API接口服务 - Go版本
 *
 * 功能说明：
 * - 支持TRC10代币操作
 * - 支持TRC20代币操作（包括USDT）
 * - 支持TRX原生代币操作
 * - 支持助记词生成地址
 * - 支持区块链查询功能
 *
 * 作者：纸飞机(Telegram): https://t.me/king_orz
 * 日期：2025年8月
 *
 * 温馨提示：接受各种代码定制
 */

// 全局配置
var config = &types.Config{
	Port:            "9527",
	TronAPIURL:      "https://api.trongrid.io",
	ContractAddress: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", // USDT TRC20 合约地址
	Decimals:        6,                                    // USDT 精度
}

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r, config)

	// 启动服务器
	utils.PrintASCIIArt()
	fmt.Printf("🚀 TRON API服务启动成功！\n")
	fmt.Printf("📍 服务地址: http://localhost:%s\n", config.Port)
	fmt.Printf("📚 接口文档: http://localhost:%s/doc\n", config.Port)
	fmt.Printf("✈️ 技术支持: https://t.me/king_orz\n")
	fmt.Println()

	// 在后台启动服务器
	go func() {
		if err := r.Run(":" + config.Port); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待服务器启动完成
	fmt.Printf("⏳ 正在启动服务器...\n")
	time.Sleep(2 * time.Second)

	// 验证服务器是否启动成功
	homeURL := fmt.Sprintf("http://localhost:%s", config.Port)
	if utils.CheckServerReady(homeURL) {
		fmt.Printf("✅ 服务器启动完成\n")
		// 自动打开浏览器
		utils.OpenBrowser(homeURL)
	} else {
		fmt.Printf("⚠️  服务器可能未完全启动，请稍后手动访问: %s\n", homeURL)
	}

	fmt.Printf("⌨️  按 Ctrl+C 停止服务器\n")
	fmt.Println()

	// 阻塞主线程，防止程序退出
	select {}
}
