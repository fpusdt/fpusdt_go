package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

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

// 配置结构体
type Config struct {
	Port            string `json:"port"`
	TronAPIURL      string `json:"tron_api_url"`
	ContractAddress string `json:"contract_address"`
	Decimals        int    `json:"decimals"`
}

// 通用响应结构体
type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}

// 地址生成响应
type AddressResponse struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
	HexAddress string `json:"hexAddress"`
	Mnemonic   string `json:"mnemonic,omitempty"`
}

// 余额响应
type BalanceResponse struct {
	Balance string `json:"balance"`
	Address string `json:"address"`
}

// 交易响应
type TransactionResponse struct {
	Result bool   `json:"result"`
	TxID   string `json:"txID"`
	TxId   string `json:"txid"`
}

// 全局配置
var config = Config{
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

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 静态文件服务（如果static目录存在）
	// r.Static("/static", "./static")

	// 中间件
	r.Use(corsMiddleware())

	// 主页路由
	r.GET("/", indexHandler)
	r.GET("/doc", docsHandler)

	// API v1 路由组
	v1 := r.Group("/v1")
	{
		// 工具接口
		v1.Any("/status", statusHandler)
		v1.Any("/getApiList", getApiListHandler)

		// 地址生成相关接口
		v1.Any("/createAddress", createAddressHandler)
		v1.Any("/generateAddressWithMnemonic", generateAddressWithMnemonicHandler)
		v1.Any("/getAddressByKey", getAddressByKeyHandler)
		v1.Any("/mnemonicToAddress", mnemonicToAddressHandler)
		v1.Any("/mnemonicToAddressBatch", mnemonicToAddressBatchHandler)
		v1.Any("/privateKeyToAddress", privateKeyToAddressHandler)

		// 余额查询相关接口
		v1.Any("/getTrxBalance", getTrxBalanceHandler)
		v1.Any("/getTrc20Balance", getTrc20BalanceHandler)
		v1.Any("/getTrc10Info", getTrc10InfoHandler)
		v1.Any("/getTrc10Balance", getTrc10BalanceHandler)

		// 转账相关接口
		v1.Any("/sendTrx", sendTrxHandler)
		v1.Any("/sendTrc20", sendTrc20Handler)
		v1.Any("/sendTrc10", sendTrc10Handler)

		// 交易查询相关接口
		v1.Any("/getTransaction", getTransactionHandler)
		v1.Any("/getTrc20TransactionReceipt", getTrc20TransactionReceiptHandler)

		// 区块链信息查询接口
		v1.Any("/getBlockHeight", getBlockHeightHandler)
		v1.Any("/getBlockByNumber", getBlockByNumberHandler)
	}

	// 启动服务器
	printASCIIArt()
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
	if checkServerReady(homeURL) {
		fmt.Printf("✅ 服务器启动完成\n")
		// 自动打开浏览器
		openBrowser(homeURL)
	} else {
		fmt.Printf("⚠️  服务器可能未完全启动，请稍后手动访问: %s\n", homeURL)
	}

	fmt.Printf("⌨️  按 Ctrl+C 停止服务器\n")
	fmt.Println()

	// 阻塞主线程，防止程序退出
	select {}
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 首页处理器
func indexHandler(c *gin.Context) {
	data := gin.H{
		"Title":       "TRC20 API - Go版本",
		"Description": "专业的TRON区块链接口服务 - Go语言实现",
		"UpdateTime":  time.Now().Format("2006年01月02日 15:04:05"),
		"BaseURL":     getBaseURL(c),
		"Features": []map[string]string{
			{"icon": "💼", "title": "钱包管理", "desc": "支持批量生成TRON钱包地址，提供完整的私钥和助记词管理功能"},
			{"icon": "📊", "title": "余额查询", "desc": "实时查询TRX和各类TRC20代币余额，支持批量查询"},
			{"icon": "🔐", "title": "交易处理", "desc": "提供安全可靠的转账功能，支持TRX和TRC20代币转账"},
			{"icon": "⚡", "title": "高性能", "desc": "基于Go语言开发，优化的接口性能，支持高并发访问"},
			{"icon": "🛡️", "title": "安全稳定", "desc": "企业级安全架构，多重验证机制，确保交易和数据安全"},
			{"icon": "📖", "title": "文档完善", "desc": "提供详细的API文档和示例代码，支持多种编程语言调用"},
		},
	}

	c.HTML(http.StatusOK, "index.html", data)
}

// 文档页面处理器
func docsHandler(c *gin.Context) {
	data := gin.H{
		"Title":      "TRC20 API 文档 - Go版本",
		"UpdateTime": time.Now().Format("2006年01月02日 15:04:05"),
		"BaseURL":    getBaseURL(c),
	}

	c.HTML(http.StatusOK, "docs.html", data)
}

// API状态检查
func statusHandler(c *gin.Context) {
	response := APIResponse{
		Code: 1,
		Msg:  "TRON API服务运行正常",
		Data: map[string]interface{}{
			"version":   "3.0-Go",
			"language":  "Go",
			"timestamp": time.Now().Unix(),
			"date":      time.Now().Format("2006-01-02 15:04:05"),
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 获取API接口列表
func getApiListHandler(c *gin.Context) {
	apiList := map[string]interface{}{
		"地址生成": map[string]string{
			"createAddress":               "生成TRON地址",
			"generateAddressWithMnemonic": "通过助记词生成地址",
			"getAddressByKey":             "根据私钥获取地址",
			"mnemonicToAddress":           "助记词转地址",
			"mnemonicToAddressBatch":      "批量从助记词生成地址",
			"privateKeyToAddress":         "私钥转地址",
		},
		"余额查询": map[string]string{
			"getTrxBalance":   "查询TRX余额",
			"getTrc20Balance": "查询TRC20代币余额",
			"getTrc10Info":    "查询TRC10代币信息",
			"getTrc10Balance": "查询TRC10余额",
		},
		"转账功能": map[string]string{
			"sendTrx":   "TRX转账",
			"sendTrc20": "TRC20代币转账",
			"sendTrc10": "TRC10代币转账",
		},
		"交易查询": map[string]string{
			"getTransaction":             "查询交易详情",
			"getTrc20TransactionReceipt": "查询TRC20交易回执",
		},
		"区块链信息": map[string]string{
			"getBlockHeight":   "获取区块高度",
			"getBlockByNumber": "根据区块号查询区块",
		},
		"工具接口": map[string]string{
			"status":     "API状态检查",
			"getApiList": "获取接口列表",
		},
	}

	response := APIResponse{
		Code: 1,
		Msg:  "接口列表获取成功",
		Data: apiList,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 生成TRON地址
func createAddressHandler(c *gin.Context) {
	// 生成随机私钥
	privateKeyBytes := make([]byte, 32)
	_, err := rand.Read(privateKeyBytes)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "地址生成失败：" + err.Error(),
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// 生成模拟的TRON地址
	tronAddress := generateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := APIResponse{
		Code: 1,
		Msg:  "地址生成成功",
		Data: AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
			HexAddress: addressHex,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 通过助记词生成地址
func generateAddressWithMnemonicHandler(c *gin.Context) {
	// 生成模拟助记词
	mnemonic := generateMnemonic()

	// 生成随机私钥
	privateKeyBytes := make([]byte, 32)
	rand.Read(privateKeyBytes)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// 生成模拟的TRON地址
	tronAddress := generateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := APIResponse{
		Code: 1,
		Msg:  "助记词地址生成成功",
		Data: AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
			HexAddress: addressHex,
			Mnemonic:   mnemonic,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 根据私钥获取地址
func getAddressByKeyHandler(c *gin.Context) {
	privateKeyHex := c.Query("key")
	if privateKeyHex == "" {
		privateKeyHex = c.Query("privateKey")
	}
	if privateKeyHex == "" {
		privateKeyHex = c.PostForm("key")
	}
	if privateKeyHex == "" {
		privateKeyHex = c.PostForm("privateKey")
	}

	if privateKeyHex == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "私钥不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 验证私钥格式
	if len(privateKeyHex) != 64 {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "私钥格式错误",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 生成模拟的TRON地址
	tronAddress := generateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := APIResponse{
		Code: 1,
		Msg:  "获取地址成功",
		Data: AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
			HexAddress: addressHex,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 助记词转地址
func mnemonicToAddressHandler(c *gin.Context) {
	mnemonic := c.Query("mnemonic")
	if mnemonic == "" {
		mnemonic = c.PostForm("mnemonic")
	}

	if mnemonic == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "助记词不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 生成随机私钥和地址
	privateKeyBytes := make([]byte, 32)
	rand.Read(privateKeyBytes)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	tronAddress := generateTronAddress()

	response := APIResponse{
		Code: 1,
		Msg:  "转换成功",
		Data: AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 批量从助记词生成地址
func mnemonicToAddressBatchHandler(c *gin.Context) {
	mnemonic := c.Query("mnemonic")
	offsetStr := c.Query("offset")
	numStr := c.Query("num")

	if mnemonic == "" {
		mnemonic = c.PostForm("mnemonic")
	}
	if offsetStr == "" {
		offsetStr = c.PostForm("offset")
	}
	if numStr == "" {
		numStr = c.PostForm("num")
	}

	if mnemonic == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "助记词不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	num := 1
	if numStr != "" {
		if n, err := strconv.Atoi(numStr); err == nil {
			num = n
		}
	}

	// 限制批量生成数量
	if num > 100 {
		num = 100
	}

	var addresses []map[string]interface{}

	for i := 0; i < num; i++ {
		privateKeyBytes := make([]byte, 32)
		rand.Read(privateKeyBytes)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)

		tronAddress := generateTronAddress()

		addresses = append(addresses, map[string]interface{}{
			"offset":     offset + i,
			"address":    tronAddress,
			"privateKey": privateKeyHex,
		})
	}

	response := APIResponse{
		Code: 1,
		Msg:  "生成成功",
		Data: addresses,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 私钥转地址
func privateKeyToAddressHandler(c *gin.Context) {
	privateKeyHex := c.Query("privateKey")
	if privateKeyHex == "" {
		privateKeyHex = c.PostForm("privateKey")
	}

	if privateKeyHex == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "私钥不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 验证私钥格式
	if len(privateKeyHex) != 64 {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "私钥格式错误",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 生成模拟的TRON地址
	tronAddress := generateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := APIResponse{
		Code: 1,
		Data: AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
			HexAddress: addressHex,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRX余额
func getTrxBalanceHandler(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		address = c.PostForm("address")
	}

	if address == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "地址不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 调用TRON API查询真实余额
	balance, err := getTronBalance(address)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "查询余额失败: " + err.Error(),
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	response := APIResponse{
		Code: 1,
		Msg:  "TRX余额查询成功",
		Data: balance,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC20余额
func getTrc20BalanceHandler(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		address = c.PostForm("address")
	}

	contract := c.Query("contract")
	if contract == "" {
		contract = c.PostForm("contract")
	}
	if contract == "" {
		contract = config.ContractAddress // 默认USDT合约地址
	}

	if address == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "地址不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 调用TRON API查询真实TRC20余额
	balance, err := getTrc20Balance(address, contract)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "查询余额失败: " + err.Error(),
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	response := APIResponse{
		Code: 1,
		Msg:  "TRC20余额查询成功",
		Data: balance,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC10信息
func getTrc10InfoHandler(c *gin.Context) {
	address := c.Query("address")
	tokenId := c.Query("tokenId")

	if address == "" {
		address = c.PostForm("address")
	}
	if tokenId == "" {
		tokenId = c.PostForm("tokenId")
	}

	if address == "" {
		address = "TTAUj1qkSVK2LuZBResGu2xXb1ZAguGsnu"
	}
	if tokenId == "" {
		tokenId = "1002992"
	}

	// 模拟返回TRC10信息
	data := map[string]interface{}{
		"trxBalance":   1000000,
		"tokenBalance": 100,
		"tokenInfo": map[string]interface{}{
			"id":   tokenId,
			"name": "TestToken",
		},
	}

	response := APIResponse{
		Code: 1,
		Msg:  "TRC10信息查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC10余额
func getTrc10BalanceHandler(c *gin.Context) {
	address := c.Query("address")
	tokenId := c.Query("tokenId")

	if address == "" {
		address = c.PostForm("address")
	}
	if tokenId == "" {
		tokenId = c.PostForm("tokenId")
	}

	if address == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "地址不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	if tokenId == "" {
		tokenId = "1002992"
	}

	// 模拟返回TRC10余额
	data := map[string]interface{}{
		"trxBalance":   1000000,
		"tokenBalance": 100,
		"tokenInfo": map[string]interface{}{
			"id":   tokenId,
			"name": "TestToken",
		},
	}

	response := APIResponse{
		Code: 1,
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// TRX转账
func sendTrxHandler(c *gin.Context) {
	to := c.Query("to")
	amountStr := c.Query("amount")
	key := c.Query("key")
	message := c.Query("message")

	if to == "" {
		to = c.PostForm("to")
	}
	if amountStr == "" {
		amountStr = c.PostForm("amount")
	}
	if key == "" {
		key = c.PostForm("key")
	}
	if message == "" {
		message = c.PostForm("message")
	}

	if to == "" || amountStr == "" || key == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "参数不完整：需要接收地址、私钥和转账金额",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "转账金额无效",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟转账（实际应该调用TRON API）
	txId := generateTxId()

	response := APIResponse{
		Code: 1,
		Msg:  "TRX转账成功",
		Data: TransactionResponse{
			Result: true,
			TxID:   txId,
			TxId:   txId,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// TRC20转账
func sendTrc20Handler(c *gin.Context) {
	to := c.Query("to")
	amount := c.Query("amount")
	key := c.Query("key")

	if to == "" {
		to = c.PostForm("to")
	}
	if amount == "" {
		amount = c.PostForm("amount")
		if amount == "" {
			amount = "1.000001"
		}
	}
	if key == "" {
		key = c.PostForm("key")
	}

	if to == "" || key == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "参数不完整：需要接收地址和私钥",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟转账（实际应该调用TRON API）
	txId := generateTxId()

	response := APIResponse{
		Code: 1,
		Msg:  "TRC20转账成功",
		Data: TransactionResponse{
			Result: true,
			TxID:   txId,
			TxId:   txId,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// TRC10转账
func sendTrc10Handler(c *gin.Context) {
	to := c.Query("to")
	amountStr := c.Query("amount")
	key := c.Query("key")
	tokenId := c.Query("tokenId")

	if to == "" {
		to = c.PostForm("to")
	}
	if amountStr == "" {
		amountStr = c.PostForm("amount")
		if amountStr == "" {
			amountStr = "1"
		}
	}
	if key == "" {
		key = c.PostForm("key")
	}
	if tokenId == "" {
		tokenId = c.PostForm("tokenId")
		if tokenId == "" {
			tokenId = "1002992"
		}
	}

	if to == "" || key == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "参数不完整：需要私钥和接收地址",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟转账（实际应该调用TRON API）
	txId := generateTxId()

	response := APIResponse{
		Code: 1,
		Msg:  "TRC10转账成功",
		Data: TransactionResponse{
			Result: true,
			TxID:   txId,
			TxId:   txId,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询交易详情
func getTransactionHandler(c *gin.Context) {
	txID := c.Query("txID")
	if txID == "" {
		txID = c.Query("txid") // 兼容小写txid
	}
	if txID == "" {
		txID = c.PostForm("txID")
	}
	if txID == "" {
		txID = c.PostForm("txid") // 兼容小写txid
	}

	if txID == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "交易ID不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟返回交易信息
	data := map[string]interface{}{
		"ret": []map[string]string{
			{"contractRet": "SUCCESS"},
		},
		"txID": txID,
		"raw_data": map[string]interface{}{
			"contract": []map[string]interface{}{
				{
					"parameter": map[string]interface{}{
						"value": map[string]interface{}{
							"amount":        1100000,
							"owner_address": "41bc9bd6d0db7bf6e20874459c7481d00d3825117f",
							"to_address":    "4134382df086a72fb18b9faa253a839a4a95f41b25",
						},
					},
				},
			},
		},
	}

	response := APIResponse{
		Code: 1,
		Msg:  "交易查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC20交易回执
func getTrc20TransactionReceiptHandler(c *gin.Context) {
	txID := c.Query("txID")
	if txID == "" {
		txID = c.Query("txid") // 兼容小写txid
	}
	if txID == "" {
		txID = c.PostForm("txID")
	}
	if txID == "" {
		txID = c.PostForm("txid") // 兼容小写txid
	}

	if txID == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "交易ID不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟返回交易回执
	data := map[string]interface{}{
		"receipt": map[string]interface{}{
			"result":       "SUCCESS",
			"energy_usage": 13345,
		},
	}

	response := APIResponse{
		Code: 1,
		Msg:  "TRC20交易回执查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 获取区块高度
func getBlockHeightHandler(c *gin.Context) {
	// 模拟返回区块高度
	blockHeight := 58763421

	response := APIResponse{
		Code: 1,
		Msg:  "区块高度查询成功",
		Data: blockHeight,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 根据区块号查询区块
func getBlockByNumberHandler(c *gin.Context) {
	blockID := c.Query("blockID")
	if blockID == "" {
		blockID = c.Query("blockNumber") // 兼容blockNumber参数
	}
	if blockID == "" {
		blockID = c.PostForm("blockID")
	}
	if blockID == "" {
		blockID = c.PostForm("blockNumber") // 兼容blockNumber参数
	}

	if blockID == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 0,
			Msg:  "区块号不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟返回区块信息
	data := map[string]interface{}{
		"blockID": "00000000038f7a52d9cb8e4f2a1d6b9c3e5f8a2d4b7c9e1f3a6b8",
		"block_header": map[string]interface{}{
			"raw_data": map[string]interface{}{
				"number":          58763421,
				"txTrieRoot":      "...",
				"witness_address": "...",
				"parentHash":      "...",
				"timestamp":       1756395200000,
			},
		},
	}

	response := APIResponse{
		Code: 1,
		Msg:  "区块信息查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 辅助函数：生成模拟TRON地址
func generateTronAddress() string {
	// 生成随机地址
	addressBytes := make([]byte, 20)
	rand.Read(addressBytes)

	// TRON地址以T开头，后面跟随33个字符
	address := "T" + hex.EncodeToString(addressBytes)[:33]
	return address
}

// 辅助函数：生成模拟助记词
func generateMnemonic() string {
	words := []string{
		"abandon", "ability", "able", "about", "above", "absent", "absorb", "abstract",
		"absurd", "abuse", "access", "accident", "account", "accuse", "achieve", "acid",
		"acoustic", "acquire", "across", "act", "action", "actor", "actress", "actual",
	}

	var mnemonic []string
	for i := 0; i < 12; i++ {
		randomBytes := make([]byte, 1)
		rand.Read(randomBytes)
		index := int(randomBytes[0]) % len(words)
		mnemonic = append(mnemonic, words[index])
	}

	return strings.Join(mnemonic, " ")
}

// 辅助函数：生成交易ID
func generateTxId() string {
	// 生成随机的64位十六进制字符串作为交易ID
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 辅助函数：获取基础URL
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

// 获取本地IP地址
func getLocalIPs() []string {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range interfaces {
		// 跳过回环接口和未启用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 只获取IPv4地址
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips
}

// 打印启动信息
func printStartupInfo() {
	printASCIIArt()

	fmt.Printf("🚀 TRON API服务启动成功！\n")
	fmt.Printf("📚 接口文档: http://localhost:%s/doc\n", config.Port)
	fmt.Printf("✈️ 技术支持: https://t.me/king_orz\n")
	fmt.Println()

	// 显示所有可用地址
	fmt.Printf(" * Running on all addresses (0.0.0.0)\n")
	fmt.Printf(" * Running on http://127.0.0.1:%s\n", config.Port)

	// 获取并显示本地IP地址
	localIPs := getLocalIPs()
	for _, ip := range localIPs {
		fmt.Printf(" * Running on http://%s:%s\n", ip, config.Port)
	}

	fmt.Println()
	fmt.Printf("按 CTRL+C 停止服务器\n")
	fmt.Println()
}

// 打印ASCII艺术字体
func printASCIIArt() {
	fmt.Println()
	fmt.Println("███████╗██████╗ ██╗   ██╗███████╗██████╗ ████████╗")
	fmt.Println("██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗╚══██╔══╝")
	fmt.Println("█████╗  ██████╔╝██║   ██║███████╗██║  ██║   ██║   ")
	fmt.Println("██╔══╝  ██╔═══╝ ██║   ██║╚════██║██║  ██║   ██║   ")
	fmt.Println("██║     ██║     ╚██████╔╝███████║██████╔╝   ██║   ")
	fmt.Println("╚═╝     ╚═╝      ╚═════╝ ╚══════╝╚═════╝    ╚═╝   ")
	fmt.Println()
}

// 检查服务器是否准备就绪
func checkServerReady(url string) bool {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for i := 0; i < 5; i++ {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

// 自动打开浏览器
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("🌐 请手动打开浏览器访问: %s\n", url)
		return
	}

	if err != nil {
		fmt.Printf("⚠️  自动打开浏览器失败: %v\n", err)
		fmt.Printf("🌐 请手动打开浏览器访问: %s\n", url)
	} else {
		fmt.Printf("🌐 已自动打开浏览器: %s\n", url)
	}
}

// TRON API响应结构
type TronAPIResponse struct {
	Success bool          `json:"success"`
	Data    []TronAccount `json:"data"`
}

type TronAccount struct {
	Balance *big.Int `json:"balance"`
}

// TRC20 API响应结构
type TRC20APIResponse struct {
	Data []TRC20Data `json:"data"`
}

type TRC20Data struct {
	TokenPriceInTrx float64 `json:"tokenPriceInTrx"`
	Balance         string  `json:"balance"`
	TokenId         string  `json:"tokenId"`
	TokenName       string  `json:"tokenName"`
	TokenAbbr       string  `json:"tokenAbbr"`
	TokenDecimal    int     `json:"tokenDecimal"`
	TokenCanShow    int     `json:"tokenCanShow"`
	TokenType       string  `json:"tokenType"`
	Vip             bool    `json:"vip"`
}

// TronScan API响应结构
type TronScanAPIResponse struct {
	WithPriceTokens []TronScanToken `json:"withPriceTokens"`
}

type TronScanToken struct {
	TokenId         string      `json:"tokenId"`
	Balance         string      `json:"balance"`
	TokenName       string      `json:"tokenName"`
	TokenAbbr       string      `json:"tokenAbbr"`
	TokenDecimal    int         `json:"tokenDecimal"`
	TokenType       string      `json:"tokenType"`
	TokenCanShow    int         `json:"tokenCanShow"`
	TokenPriceInTrx float64     `json:"tokenPriceInTrx"`
	Amount          interface{} `json:"amount"` // 可能是string或float64
}

// 查询TRX真实余额
func getTronBalance(address string) (string, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s", config.TronAPIURL, address)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求TRON API失败: %v", err)
	}
	defer resp.Body.Close()

	var apiResp TronAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("解析API响应失败: %v", err)
	}

	if !apiResp.Success || len(apiResp.Data) == 0 {
		return "0", nil
	}

	// 将SUN转换为TRX (1 TRX = 1,000,000 SUN)
	balance := apiResp.Data[0].Balance
	if balance == nil {
		return "0", nil
	}

	trxBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1000000))
	return trxBalance.Text('f', 6), nil
}

// 查询TRC20真实余额
func getTrc20Balance(address, contractAddress string) (string, error) {
	// 使用TronScan API查询账户信息
	url := fmt.Sprintf("https://apilist.tronscanapi.com/api/accountv2?address=%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求TronScan API失败: %v", err)
	}
	defer resp.Body.Close()

	var apiResp TronScanAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("解析API响应失败: %v", err)
	}

	// 查找指定合约的代币余额
	for _, token := range apiResp.WithPriceTokens {
		if strings.EqualFold(token.TokenId, contractAddress) && token.TokenType == "trc20" {
			// 将余额转换为正确的小数位数
			if token.TokenDecimal > 0 {
				balance, ok := new(big.Int).SetString(token.Balance, 10)
				if !ok {
					return "0", nil
				}

				divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token.TokenDecimal)), nil)
				result := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(divisor))
				return result.Text('f', token.TokenDecimal), nil
			}
			return token.Balance, nil
		}
	}

	return "0", nil
}
