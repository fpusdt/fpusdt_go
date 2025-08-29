@echo off
chcp 65001 >nul
title FPUSDT - TRON API Go版本 环境配置

echo 🚀 开始配置 FPUSDT 开发环境...
echo.

REM 检查 Go 是否已安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未检测到 Go 环境
    echo 请先安装 Go 1.19+ 版本: https://golang.org/dl/
    pause
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo ✅ 检测到 Go 版本: %GO_VERSION%
echo.

echo 🔧 配置 Go 模块代理...

REM 设置环境变量（临时）
set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=sum.golang.google.cn

echo ✅ 临时环境变量已设置

REM 设置环境变量（永久）
echo.
echo 📝 设置永久环境变量...

setx GO111MODULE on >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  需要管理员权限设置永久环境变量
    echo 请以管理员身份运行此脚本，或手动设置环境变量
) else (
    echo ✅ GO111MODULE 设置成功
)

setx GOPROXY https://goproxy.cn,direct >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  GOPROXY 设置失败
) else (
    echo ✅ GOPROXY 设置成功
)

setx GOSUMDB sum.golang.google.cn >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  GOSUMDB 设置失败
) else (
    echo ✅ GOSUMDB 设置成功
)

REM 验证配置
echo.
echo 🔍 验证当前配置...
echo GO111MODULE: %GO111MODULE%
echo GOPROXY: %GOPROXY%
echo GOSUMDB: %GOSUMDB%

REM 测试网络连接
echo.
echo 🌐 测试代理连接...
go list -m golang.org/x/tools >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ 代理连接正常
) else (
    echo ⚠️  代理连接测试失败，可能需要检查网络
)

REM 安装依赖
echo.
echo 📦 开始安装项目依赖...
go mod tidy
if %errorlevel% equ 0 (
    echo ✅ 依赖安装成功
) else (
    echo ❌ 依赖安装失败，请检查网络连接
    pause
    exit /b 1
)

echo.
echo 🎉 环境配置完成！
echo.
echo 📝 下一步操作：
echo    1. 重新打开命令提示符（使永久环境变量生效）
echo    2. 启动服务: go run main.go
echo    3. 访问主页: http://localhost:9527
echo.
echo 🆘 如有问题，请联系: https://t.me/king_orz
echo.
pause
