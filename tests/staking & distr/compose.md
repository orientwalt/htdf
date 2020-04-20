# MainNet
### action
    make hstart
    make hstest
    make hstop
    sudo make clean
### hstest
    [alias]
    cp hs* /usr/local/bin
    alias hscli="hscli --home node4/.hscli"
    alias hsd="hsd --home node4/.hsd"

    [config/start]
    hscli config chain-id testchain
    hscli config trust-node true
    hscli config node http://192.168.10.2:26657
    # nohup hsd start & > /dev/null
    # clear

    [transactions]
    hscli accounts new 12345678
    hscli accounts list
    hscli query account $(hscli accounts list| sed -n '1p')
    hscli tx send $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 20000000satoshi --gas-price=100
    hscli tx send $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 20000000satoshi --gas-price=100

    [validators]
    - show yours
    hsd tendermint show-validator
    - check status
    hscli query staking validators
    hscli  query staking validator [htdfvaloper]
    - confirm running
    hscli query tendermint-validator-set
    hscli query tendermint-validator-set | grep [htdfvalcons/htdfvalconspub]
    hscli query tendermint-validator-set | grep "$(hsd tendermint show-validator)"
    - start yours(tip: 100,000,000 for voting power 1, 1,000,000,000 for 10, 10,000,000,000 for 100)
    hscli tx staking create-validator $(hscli accounts list| sed -n '2p') \
                                       --pubkey=$(hsd tendermint show-validator)\
                                       --amount=100000000satoshi \
                                       --moniker=client \
                                       --commission-rate=0.10 \
                                       --commission-max-rate=0.20 \
                                       --commission-max-change-rate=0.01 \
                                       --min-self-delegation=1 \
                                       --gas-price=100
    hscli tx staking edit-validator $(hscli accounts list| sed -n '2p') --gas-price=100
    or
    hscli tx staking edit-validator $(hscli accounts list| sed -n '2p')\
                --moniker=client \
                --chain-id=testchain \
                --website="https://htdf.network" \
                --identity=23870f5bb12ba2c4967c46db \
                --details="To infinity and beyond!" \
                --gas-price=100 \
                --commission-rate=0.10
    - unjail
    hscli tx slashing unjail $(hscli accounts list| sed -n '2p') --gas-price=100
    - log
    hscli query slashing signing-info [htdfvalconspub]
    - check

    [delegators]
    hscli tx staking delegate $(hscli accounts list| sed -n '2p') \
                               $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2)\
                               100000000satoshi --gas-adjustment=1.5 --gas-price=100
    hscli tx staking redelegate $(hscli accounts list| sed -n '2p') \
                                 $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '3p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 100000000satoshi --gas-adjustment=1.5 --gas-price=100
    hscli query distr rewards $(hscli accounts list| sed -n '2p')
    hscli query account $(hscli accounts list| sed -n '2p')
    hscli tx staking unbond $(hscli accounts list| sed -n '1p') \
                             $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2) \
                             100000000satoshi \
                             --gas-adjustment 1.5 --gas-price=100
    hscli tx distr withdraw-rewards $(hscli accounts list| sed -n '1p') \
                                     $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2) \ --gas-adjustment 1.5 --gas-price=100