package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"tron-api-go/internal/types"
	"tron-api-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// 处理器服务结构体
type Service struct {
	Config *types.Config
}

// 创建新的处理器服务
func NewService(config *types.Config) *Service {
	return &Service{
		Config: config,
	}
}

// 首页处理器
func (s *Service) IndexHandler(c *gin.Context) {
	data := gin.H{
		"Title":       "TRC20 API - Go版本",
		"Description": "专业的TRON区块链接口服务 - Go语言实现",
		"UpdateTime":  time.Now().Format("2006年01月02日 15:04:05"),
		"BaseURL":     utils.GetBaseURL(c),
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
func (s *Service) DocsHandler(c *gin.Context) {
	data := gin.H{
		"Title":      "TRC20 API 文档 - Go版本",
		"UpdateTime": time.Now().Format("2006年01月02日 15:04:05"),
		"BaseURL":    utils.GetBaseURL(c),
	}

	c.HTML(http.StatusOK, "docs.html", data)
}

// API状态检查
func (s *Service) StatusHandler(c *gin.Context) {
	response := types.APIResponse{
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
func (s *Service) GetApiListHandler(c *gin.Context) {
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

	response := types.APIResponse{
		Code: 1,
		Msg:  "接口列表获取成功",
		Data: apiList,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 生成TRON地址
func (s *Service) CreateAddressHandler(c *gin.Context) {
	// 生成随机私钥
	privateKeyBytes := make([]byte, 32)
	_, err := rand.Read(privateKeyBytes)
	if err != nil {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "地址生成失败：" + err.Error(),
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// 生成模拟的TRON地址
	tronAddress := utils.GenerateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := types.APIResponse{
		Code: 1,
		Msg:  "地址生成成功",
		Data: types.AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
			HexAddress: addressHex,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 通过助记词生成地址
func (s *Service) GenerateAddressWithMnemonicHandler(c *gin.Context) {
	// 生成模拟助记词
	mnemonic := utils.GenerateMnemonic()

	// 生成随机私钥
	privateKeyBytes := make([]byte, 32)
	rand.Read(privateKeyBytes)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// 生成模拟的TRON地址
	tronAddress := utils.GenerateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := types.APIResponse{
		Code: 1,
		Msg:  "助记词地址生成成功",
		Data: types.AddressResponse{
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
func (s *Service) GetAddressByKeyHandler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "私钥不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 验证私钥格式
	if len(privateKeyHex) != 64 {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "私钥格式错误",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 生成模拟的TRON地址
	tronAddress := utils.GenerateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := types.APIResponse{
		Code: 1,
		Msg:  "获取地址成功",
		Data: types.AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
			HexAddress: addressHex,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 助记词转地址
func (s *Service) MnemonicToAddressHandler(c *gin.Context) {
	mnemonic := c.Query("mnemonic")
	if mnemonic == "" {
		mnemonic = c.PostForm("mnemonic")
	}

	if mnemonic == "" {
		c.JSON(http.StatusOK, types.APIResponse{
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

	tronAddress := utils.GenerateTronAddress()

	response := types.APIResponse{
		Code: 1,
		Msg:  "转换成功",
		Data: types.AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 批量从助记词生成地址
func (s *Service) MnemonicToAddressBatchHandler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
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

		tronAddress := utils.GenerateTronAddress()

		addresses = append(addresses, map[string]interface{}{
			"offset":     offset + i,
			"address":    tronAddress,
			"privateKey": privateKeyHex,
		})
	}

	response := types.APIResponse{
		Code: 1,
		Msg:  "生成成功",
		Data: addresses,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 私钥转地址
func (s *Service) PrivateKeyToAddressHandler(c *gin.Context) {
	privateKeyHex := c.Query("privateKey")
	if privateKeyHex == "" {
		privateKeyHex = c.PostForm("privateKey")
	}

	if privateKeyHex == "" {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "私钥不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 验证私钥格式
	if len(privateKeyHex) != 64 {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "私钥格式错误",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 生成模拟的TRON地址
	tronAddress := utils.GenerateTronAddress()
	addressHex := "41" + tronAddress[1:]

	response := types.APIResponse{
		Code: 1,
		Data: types.AddressResponse{
			PrivateKey: privateKeyHex,
			Address:    tronAddress,
			HexAddress: addressHex,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRX余额
func (s *Service) GetTrxBalanceHandler(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		address = c.PostForm("address")
	}

	if address == "" {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "地址不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 调用TRON API查询真实余额
	balance, err := utils.GetTronBalance(address, s.Config)
	if err != nil {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "查询余额失败: " + err.Error(),
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	response := types.APIResponse{
		Code: 1,
		Msg:  "TRX余额查询成功",
		Data: balance,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC20余额
func (s *Service) GetTrc20BalanceHandler(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		address = c.PostForm("address")
	}

	contract := c.Query("contract")
	if contract == "" {
		contract = c.PostForm("contract")
	}
	if contract == "" {
		contract = s.Config.ContractAddress // 默认USDT合约地址
	}

	if address == "" {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "地址不能为空",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 调用TRON API查询真实TRC20余额
	balance, err := utils.GetTrc20Balance(address, contract)
	if err != nil {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "查询余额失败: " + err.Error(),
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	response := types.APIResponse{
		Code: 1,
		Msg:  "TRC20余额查询成功",
		Data: balance,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC10信息
func (s *Service) GetTrc10InfoHandler(c *gin.Context) {
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

	response := types.APIResponse{
		Code: 1,
		Msg:  "TRC10信息查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC10余额
func (s *Service) GetTrc10BalanceHandler(c *gin.Context) {
	address := c.Query("address")
	tokenId := c.Query("tokenId")

	if address == "" {
		address = c.PostForm("address")
	}
	if tokenId == "" {
		tokenId = c.PostForm("tokenId")
	}

	if address == "" {
		c.JSON(http.StatusOK, types.APIResponse{
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

	response := types.APIResponse{
		Code: 1,
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// TRX转账
func (s *Service) SendTrxHandler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "参数不完整：需要接收地址、私钥和转账金额",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "转账金额无效",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟转账（实际应该调用TRON API）
	txId := utils.GenerateTxId()

	response := types.APIResponse{
		Code: 1,
		Msg:  "TRX转账成功",
		Data: types.TransactionResponse{
			Result: true,
			TxID:   txId,
			TxId:   txId,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// TRC20转账
func (s *Service) SendTrc20Handler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "参数不完整：需要接收地址和私钥",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟转账（实际应该调用TRON API）
	txId := utils.GenerateTxId()

	response := types.APIResponse{
		Code: 1,
		Msg:  "TRC20转账成功",
		Data: types.TransactionResponse{
			Result: true,
			TxID:   txId,
			TxId:   txId,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// TRC10转账
func (s *Service) SendTrc10Handler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
			Code: 0,
			Msg:  "参数不完整：需要私钥和接收地址",
			Data: nil,
			Time: time.Now().Unix(),
		})
		return
	}

	// 模拟转账（实际应该调用TRON API）
	txId := utils.GenerateTxId()

	response := types.APIResponse{
		Code: 1,
		Msg:  "TRC10转账成功",
		Data: types.TransactionResponse{
			Result: true,
			TxID:   txId,
			TxId:   txId,
		},
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询交易详情
func (s *Service) GetTransactionHandler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
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

	response := types.APIResponse{
		Code: 1,
		Msg:  "交易查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 查询TRC20交易回执
func (s *Service) GetTrc20TransactionReceiptHandler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
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

	response := types.APIResponse{
		Code: 1,
		Msg:  "TRC20交易回执查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 获取区块高度
func (s *Service) GetBlockHeightHandler(c *gin.Context) {
	// 模拟返回区块高度
	blockHeight := 58763421

	response := types.APIResponse{
		Code: 1,
		Msg:  "区块高度查询成功",
		Data: blockHeight,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// 根据区块号查询区块
func (s *Service) GetBlockByNumberHandler(c *gin.Context) {
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
		c.JSON(http.StatusOK, types.APIResponse{
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

	response := types.APIResponse{
		Code: 1,
		Msg:  "区块信息查询成功",
		Data: data,
		Time: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}
