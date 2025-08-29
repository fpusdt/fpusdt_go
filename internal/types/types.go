package types

import (
	"math/big"
)

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
