[![](https://godoc.org/github.com/orientwalt/htdf?status.svg)](http://godoc.org/github.com/orientwalt/htdf) [![Go Report Card](https://goreportcard.com/badge/github.com/orientwalt/htdf)](https://goreportcard.com/report/github.com/orientwalt/htdf)
[![version](https://img.shields.io/github/tag/orientwalt/htdf.svg)](https://github.com/orientwalt/htdf/releases/latest)
[![Go version](https://img.shields.io/badge/go-1.12.9-blue.svg)](https://github.com/moovweb/gvm)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)

# HTDF

- [English](./README.md)
- [中文](./docs/chinese/README.md)

## Introduction
HTDF is a high-performance public blockchain developed by HTDF Foundation. It is based on [ethereum](https://github.com/ethereum/go-ethereum) & [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) on [tendermint](https://github.com/tendermint/tendermint)  . We merged tendermint's consensus, cosmos-sdk's application logic, and ethereum's account system & smart contract architecture into a brand new architecture - htdf blockchain. 

HTDF main chain had been released at 2020-03-01.This project is now under active and continuous development.
   
**Note**: Requires Go 12.9+

## Executables

```
hsd
hscli
```

## [Quick Start](https://github.com/orientwalt/htdf/blob/master/docs/build%20%26%20run.md)

> Install compiler and tools
> - Install `go` (requires go12.9+): https://golang.google.cn/doc/install
> - Install `make` and `gcc` by `sudo apt install make gcc -y` or  `yum install make gcc -y`

You can follow the below steps. You should type password for your genesis account after runing `make new`.

```
git clone https://github.com/orientwalt/htdf.git
cd htdf
make new
make start
tail -f ~/.hsd/app.log
```

##  Resources

- Official Website: https://www.orientwalt.com/
- Chinese API documents: https://gitee.com/orientwalt/apidoc_2020 
- Blockchain explorer: https://www.htdfscan.com/
- Build and run on Windows: [build_run_on_windows.md](./docs/build_run_on_windows.md)


## Contributions
We always welcome any issues and Pull Requests.
