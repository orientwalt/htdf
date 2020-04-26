### account
    hscli accounts new 12345678
### query
    hscli query account $(hscli accounts list| sed -n '2p')
### bank
    hscli tx send $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 20000000satoshi --gas-price=100
    hscli tx create $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 1000satoshi --gas-price=100 --encode=false
    hscli tx sign [rawcode] --encode=false
### staking
    [validating]
    hscli tx staking create-validator $(hscli accounts list| sed -n '1p') \
                                       --pubkey=$(hsd tendermint show-validator)\
                                       --amount=100000000satoshi \
                                       --moniker=client \
                                       --commission-rate=0.10 \
                                       --commission-max-rate=0.20 \
                                       --commission-max-change-rate=0.01 \
                                       --min-self-delegation=1 \
                                       --gas-price=100
    hscli tx staking edit-validator $(hscli accounts list| sed -n '2p') --gas-price=100

    [delegating]
    hscli tx staking delegate $(hscli accounts list| sed -n '2p') <validatorAddress> 10000000satoshi --gas-adjustment=1.5 --gas-price=100