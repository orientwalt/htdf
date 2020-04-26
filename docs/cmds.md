### configuration
    hscli config chain-id [chain-id]

### accounts cmds
    hscli accounts newaccount
    hscli accounts listaccounts
    hscli accounts genprivkey [addr]
    hscli accounts getbalance [addr]

### transaction cmds
    hscli tx send [fromaddr] [toaddr] [amount]
    hscli tx create [fromaddr] [toaddr] [amount]
    hscli tx sign [rawdata]
    hscli tx broadcast [rawdata]

### query cmds
```
hscli query accounts [addr]
hscli query block
hscli query txs
hscli query tx

[additional]
hscli query rewards [block-height]
hscli query total-provisions

[contractcall]
contract-addr: htdf1l03rqalmg58wgw9ya39wwc3466lyy20xnpeaee
callcode: 27e235e300000000000000000000000027681ceb7de9bae3c5f7f10f81ff5106e2ca48a7
hscli query contract [contract-addr] [callcode]
```
### check
    hscli query staking pool
    hscli query staking params
    hscli query distr params

### [staking cmds](https://github.com/orientwalt/htdf/blob/master/x/staking/client/cli/tx.go)
    delegator-addr: htdf1zf07fyt2an2ral8zve0u4y7lzqa6x4lqfeyl8m
    validator-addr: htdfvaloper1zf07fyt2an2ral8zve0u4y7lzqa6x4lqrquxss
    amount: 100000stake
    
    [unbound]
    hscli tx staking unbond [delegator-addr] [validator-addr] [amount] --gas-adjustment 1.5 --gas-price=100

    [delegate]
    hscli tx staking delegate [delegator-addr] [validator-addr] [amount] --gas-adjustment=1.5 --gas-price=100
### [rewards](https://github.com/orientwalt/htdf/blob/master/x/distribution/client/cli/tx.go)
    [query]
    hscli query distr rewards [delegator-addr]
    hscli query distr rewards <delegator_address> <validator_address>
    hscli query distr commission <validator_address>
    hscli query distr community-pool
    hscli query rewards 1

    [withdraw]
    hscli tx distr withdraw-rewards [delegator-addr] [validator-addr] --gas-adjustment 1.5 --gas-price=100
    hscli tx distr withdraw-rewards [delegator-addr] [validator-addr] --commission --gas-adjustment 1.5 --gas-price=100