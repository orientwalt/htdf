[![CircleCI](https://circleci.com/gh/orientwalt/htdf/tree/master.svg?style=shield)](https://circleci.com/gh/orientwalt/htdf/tree/master)
[![](https://godoc.org/github.com/orientwalt/htdf?status.svg)](http://godoc.org/github.com/orientwalt/htdf) [![Go Report Card](https://goreportcard.com/badge/github.com/orientwalt/htdf)](https://goreportcard.com/report/github.com/orientwalt/htdf)
[![Travis](https://travis-ci.org/orientwalt/htdf.svg?branch=master)](https://travis-ci.org/orientwalt/htdf)
[![version](https://img.shields.io/github/tag/orientwalt/htdf.svg)](https://github.com/orientwalt/htdf/releases/latest)
[![Go version](https://img.shields.io/badge/go-1.12.9-blue.svg)](https://github.com/moovweb/gvm)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)

# HTDF
## Introduction
   HTDF is a high-performance public blockchain developed by HTDF Foundation. It is based on [ethereum](https://github.com/ethereum/go-ethereum) & [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) on [tendermint](https://github.com/tendermint/tendermint)  . We merged tendermint's consensus, cosmos-sdk's application logic, and ethereum's account system & smart contract architecture into a brand new architecture - htdf blockchain.  This project is now UNDER ACTIVE DEVELOPMENT.
   
   **Note**: Requires Go 12.9+
## Executables
```
hsd
hscli
```
## [Quick Start](https://github.com/orientwalt/htdf/blob/master/docs/build%20%26%20run.md)
Only one command is enough to set up a standalone blockchain on your local machine.
```
make new
make start
tail -f ~/.hsd/app.log
```
## Contributions
We always welcome any issue/advice/suggestion.
