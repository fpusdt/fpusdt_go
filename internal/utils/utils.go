package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"tron-api-go/internal/types"

	"github.com/gin-gonic/gin"
)

// 辅助函数：生成模拟TRON地址
func GenerateTronAddress() string {
	// 生成随机地址
	addressBytes := make([]byte, 20)
	rand.Read(addressBytes)

	// TRON地址以T开头，后面跟随33个字符
	address := "T" + hex.EncodeToString(addressBytes)[:33]
	return address
}

// 辅助函数：生成模拟助记词
func GenerateMnemonic() string {
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
func GenerateTxId() string {
	// 生成随机的64位十六进制字符串作为交易ID
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 辅助函数：获取基础URL
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

// 获取本地IP地址
func GetLocalIPs() []string {
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
func PrintStartupInfo(config *types.Config) {
	PrintASCIIArt()

	fmt.Printf("🚀 TRON API服务启动成功！\n")
	fmt.Printf("📚 接口文档: http://localhost:%s/doc\n", config.Port)
	fmt.Printf("✈️ 技术支持: https://t.me/king_orz\n")
	fmt.Println()

	// 显示所有可用地址
	fmt.Printf(" * Running on all addresses (0.0.0.0)\n")
	fmt.Printf(" * Running on http://127.0.0.1:%s\n", config.Port)

	// 获取并显示本地IP地址
	localIPs := GetLocalIPs()
	for _, ip := range localIPs {
		fmt.Printf(" * Running on http://%s:%s\n", ip, config.Port)
	}

	fmt.Println()
	fmt.Printf("按 CTRL+C 停止服务器\n")
	fmt.Println()
}

// 打印ASCII艺术字体
func PrintASCIIArt() {
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
func CheckServerReady(url string) bool {
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
func OpenBrowser(url string) {
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

// 查询TRX真实余额
func GetTronBalance(address string, config *types.Config) (string, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s", config.TronAPIURL, address)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求TRON API失败: %v", err)
	}
	defer resp.Body.Close()

	var apiResp types.TronAPIResponse
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
func GetTrc20Balance(address, contractAddress string) (string, error) {
	// 使用TronScan API查询账户信息
	url := fmt.Sprintf("https://apilist.tronscanapi.com/api/accountv2?address=%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求TronScan API失败: %v", err)
	}
	defer resp.Body.Close()

	var apiResp types.TronScanAPIResponse
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

// CORS中间件
func CorsMiddleware() gin.HandlerFunc {
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
