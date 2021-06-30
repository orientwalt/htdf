<p align="center">
  <br> English | <a href="docs/chinese/README.md">中文</a>
</p>


<p align="center">
  <a href="http://godoc.org/github.com/orientwalt/htdf"><img src="https://godoc.org/github.com/orientwalt/htdf?status.svg" alt=""></a>
  <a href="https://goreportcard.com/report/github.com/orientwalt/htdf"><img src="https://goreportcard.com/badge/github.com/orientwalt/htdf" alt=""></a>
  <a href="https://github.com/orientwalt/htdf/releases/latest"><img src="https://img.shields.io/github/tag/orientwalt/htdf.svg" alt=""></a>
  <a href="https://github.com/moovweb/gvm"><img src="https://img.shields.io/badge/go-1.14.1-blue.svg" alt=""></a>
  <a href="https://opensource.org/licenses/Apache-2.0"><img src="https://img.shields.io/badge/License-Apache%202.0-green.svg" alt=""></a>
</p>


# HTDF


## Introduction
HTDF is a high-performance public blockchain developed by HTDF Foundation. It is based on [ethereum](https://github.com/ethereum/go-ethereum) & [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) on [tendermint](https://github.com/tendermint/tendermint)  . We merged tendermint's consensus, cosmos-sdk's application logic, and ethereum's account system & smart contract architecture into a brand new architecture - htdf blockchain.

HTDF main chain had been released at 2020-03-01.This project is now under active and continuous development.

**Note**: Requires Go 14.1+

## Executables

```
hsd
hscli
```

## Quick Start

Install compiler , tools and dependencies:
- Install latest `golang` (requires go14.1+): https://golang.org/dl/
- Install `make` and `gcc` by `sudo apt install make gcc -y` or  `yum install make gcc -y`
- Install `libleveldb` by `sudo apt-get install libleveldb-dev -y` or `yum install leveldb-devel -y`

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
- Chinese API documents: https://github.com/orientwalt/docs
- Blockchain explorer: https://www.htdfscan.com/
- Build and run on Windows: [build_run_on_windows.md](./docs/build_run_on_windows.md)
- Android Wallet : https://github.com/orientwalt/walt-android-public

## Contributions
We always welcome issues and PRs.
