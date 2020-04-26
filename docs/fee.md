### cmds
    [min-gas-prices]
    hsd start --minimum-gas-price=20

    [fee & gasprice]
    hscli tx send [fromaddr] [toaddr] [amount] --gas-price=100
    hscli tx send [fromaddr] [toaddr] [amount] --gas-wanted=30000 --gas-price=100
### references
#### [client cmd: fee & gas](https://cosmos.network/docs/gaia/hscli.html#fees-gas)
#### [gaiad.toml: minimum-gas-prices](https://cosmos.network/docs/gaia/join-mainnet.html#set-minimum-gas-prices)