## steps
```
[1. stop hsd]
[2. backup .hsd]
[3. update hsd]
[4. start hsd]

[5. submit proposal]
hscli tx gov submit-proposal htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 --gas-price=100  --switch-height=120 --description="first proposal"  --title="test0" --type="software_upgrade" --deposit="1000000000satoshi" --version="1"

[6. vote]
hscli tx send htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 htdf1xvktg68uwrtkml7m4yqul2sjm6fhtglm6y20eg 1000000000satoshi --gas-price=100
hscli tx send htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 htdf1p97p8vckpvkzx34se7eansm3rp223r3trmn63h 1000000000satoshi --gas-price=100
hscli tx send htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 htdf18efa8m8yudlqc765s9vkdrev5475jyqalqm5a0 1000000000satoshi --gas-price=100
hscli tx send htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 htdf1nuxf4amphaajuwg0ph3se6kmsda9cs6v0sja7r 1000000000satoshi --gas-price=100
hscli tx gov vote  htdf1xvktg68uwrtkml7m4yqul2sjm6fhtglm6y20eg 1  yes --gas-price=100 
hscli tx gov vote  htdf1p97p8vckpvkzx34se7eansm3rp223r3trmn63h 1  yes --gas-price=100 
hscli tx gov vote  htdf18efa8m8yudlqc765s9vkdrev5475jyqalqm5a0 1  yes --gas-price=100 
hscli tx gov vote  htdf1nuxf4amphaajuwg0ph3se6kmsda9cs6v0sja7r 1  yes --gas-price=100 

[7. check]
hscli query staking params
unbonding, unslashing test

```