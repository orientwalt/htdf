## steps
```
[1. stop hsd]
[2. backup .hsd]
[3. update hsd]
[4. start hsd]

[5. submit proposal]
hscli tx gov submit-proposal htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 --gas-price=100  --switch-height=400 --description="first proposal"  --title="test0" --type="software_upgrade" --deposit="1000000000satoshi" --version="1"

[6. vote]
hscli tx send htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 htdf1rgdsa5kjyulwzy9a56qxnsfvkfgxvntcmqnkqr 1000000000satoshi --gas-price=100
hscli tx gov vote  htdf1td6ak7uygf6zamyaaq3vrfrg705hx8ye2lnvy9 1  yes --gas-price=100 

[7. check]
hscli query staking params
unbonding, unslashing test

```