<p align="center">
  <br> <a href="../../README.md">English</a>  | 中文 </br>
</p>

[![](https://godoc.org/github.com/orientwalt/htdf?status.svg)](http://godoc.org/github.com/orientwalt/htdf) [![Go Report Card](https://goreportcard.com/badge/github.com/orientwalt/htdf)](https://goreportcard.com/report/github.com/orientwalt/htdf)
[![version](https://img.shields.io/github/tag/orientwalt/htdf.svg)](https://github.com/orientwalt/htdf/releases/latest)
[![Go version](https://img.shields.io/badge/go-1.12.9-blue.svg)](https://github.com/moovweb/gvm)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)

# HTDF

## 简介

HTDF是由orientwalt(HTDF)基金组织开发的高性能的公链. HTDF基于Ethereum + Cosmos SDK + Tendermint框架开发.使用了Tendermint的共识 + Cosmos-SDK的应用逻辑 + Ethereum的账户体系和智能合约(EVM)架构.

HTDF主网已于2020年3月1日上线,目前该项目持续开发中.
   
**Note**: Requires Go 12.9+

## 可执行文件

```
hsd
hscli
```

## 快速启动

> 安装编译器
> - 安装最新版的 `go` (requires go12.9+): https://golang.google.cn/doc/install
> - 安装`make` 和 `gcc`, 通过命令 `sudo apt install make gcc -y` 或者  `yum install make gcc -y`

在htdf目录下,按照以下命令启动. 在执行 `make new` 命令之后需要输入设置的密码



```
# 国内可以使用gitee下载, git clone https://gitee.com/orientwalt/htdf.git
git clone https://github.com/orientwalt/htdf.git
cd htdf
make new
make start
tail -f ~/.hsd/app.log
```

## 更多资源

- 官网: https://www.orientwalt.com/
- 中文API文档: https://gitee.com/orientwalt/docs
- 区块链浏览器: https://www.htdfscan.com/
- 在Windows上构建和运行: [build_run_on_windows.md](./docs/build_run_on_windows.md)


## 贡献

我们欢迎提交issue和PR
