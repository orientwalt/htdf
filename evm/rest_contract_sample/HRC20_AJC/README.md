# access smart contract via RESTful interface 


## get byte code


```
$go run byte_code_sample.go ../../../tests/evm/erc20/AJC_sol_AJCToken.abi  ../../../tests/evm/erc20/AJC_sol_AJCToken.bin htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3


contractCode, create contract|Code=60606040526000600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550341561005157600080fd5b6aa49be39dc14cb8270000006003819055506003546000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610d61806100b76000396000f3006060604052600436106100af576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100b4578063095ea7b31461014257806318160ddd1461019c57806323b872dd146101c5578063313ce5671461023e5780634d853ee51461026d57806370a08231146102c257806393c32e061461030f57806395d89b4114610348578063a9059cbb146103d6578063dd62ed3e14610430575b600080fd5b34156100bf57600080fd5b6100c761049c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101075780820151818401526020810190506100ec565b50505050905090810190601f1680156101345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561014d57600080fd5b610182600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506104d5565b604051808215151515815260200191505060405180910390f35b34156101a757600080fd5b6101af61065c565b6040518082815260200191505060405180910390f35b34156101d057600080fd5b610224600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610662565b604051808215151515815260200191505060405180910390f35b341561024957600080fd5b610251610959565b604051808260ff1660ff16815260200191505060405180910390f35b341561027857600080fd5b61028061095e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156102cd57600080fd5b6102f9600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610984565b6040518082815260200191505060405180910390f35b341561031a57600080fd5b610346600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919050506109cc565b005b341561035357600080fd5b61035b610a6c565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561039b578082015181840152602081019050610380565b50505050905090810190601f1680156103c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156103e157600080fd5b610416600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610aa5565b604051808215151515815260200191505060405180910390f35b341561043b57600080fd5b610486600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610c77565b6040518082815260200191505060405180910390f35b6040805190810160405280600981526020017f414a4320636861696e000000000000000000000000000000000000000000000081525081565b60008082148061056157506000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054145b151561056c57600080fd5b81600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b60035481565b600080600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415151561072057600080fd5b80831115151561072f57600080fd5b610780836000808873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610813836000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506108688382610cfe90919063ffffffff16565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040518082815260200191505060405180910390a360019150509392505050565b601281565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610a2857600080fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6040805190810160405280600381526020017f414a43000000000000000000000000000000000000000000000000000000000081525081565b60008073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610ae257600080fd5b610b33826000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610bc6826000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a36001905092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000828211151515610d0c57fe5b818303905092915050565b6000808284019050838110151515610d2b57fe5b80915050929150505600a165627a7a7230582043a3cd97586e182885676a8c6e6413be040c6f728b9763d794ecdbfff9a4b7c90029
contractCode, transfer|accTestContractToAddress=htdf1vms0n5t80acapjnvr4t9xeelucujq58zml4kg2|Code=a9059cbb00000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2000000000000000000000000000000000000000000000000000000000000001e
contractCode, get balance|accTestContractToAddress=htdf1vms0n5t80acapjnvr4t9xeelucujq58zml4kg2|Code=70a0823100000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2
contractCode, get balance|strMinterAddress=htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3|Code=70a0823100000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9

```

## use curl to access smart contract
use REST api /hs/send to access smart contract

- /hs/send has three type of MOD
>classic transicion

>>  field "data" must be nil( "")
>>  field "amount" must be positive( amount >0)

>create smart contract

>>  field "data" must not be nil("")
>>  field "amount" must be zero( amount == 0)

>open smart contract  

>> fields same like `create smart contract`


## caution

**when input param change(minter address, from address ,to address , send amount ... ) , must  `get byte code`  and fill out the folling "data" field again**

## create contract
 use curl to create
```

# 发交易  send;           新建合约
$curl --location --request POST 'http://htdf2020-test01.orientwalt.cn:1317/hs/send' \
 --header 'Content-Type: application/x-www-form-urlencoded' \
 --data-raw '    { "base_req": 
       { "from": "htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3", 
         "memo": "",
         "password": "12345678", 
         "chain_id": "testchain", 
         "account_number": "0", 
         "sequence": "0",
         "gas_wanted": "5000000", 
         "gas_price": "100", 
         "simulate": false
       },          
       "amount": [ 
               { "denom": "htdf", 
                 "amount": "0" } ],
       "to": "",
       "data": "60606040526000600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550341561005157600080fd5b6aa49be39dc14cb8270000006003819055506003546000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610d61806100b76000396000f3006060604052600436106100af576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100b4578063095ea7b31461014257806318160ddd1461019c57806323b872dd146101c5578063313ce5671461023e5780634d853ee51461026d57806370a08231146102c257806393c32e061461030f57806395d89b4114610348578063a9059cbb146103d6578063dd62ed3e14610430575b600080fd5b34156100bf57600080fd5b6100c761049c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101075780820151818401526020810190506100ec565b50505050905090810190601f1680156101345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561014d57600080fd5b610182600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506104d5565b604051808215151515815260200191505060405180910390f35b34156101a757600080fd5b6101af61065c565b6040518082815260200191505060405180910390f35b34156101d057600080fd5b610224600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610662565b604051808215151515815260200191505060405180910390f35b341561024957600080fd5b610251610959565b604051808260ff1660ff16815260200191505060405180910390f35b341561027857600080fd5b61028061095e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156102cd57600080fd5b6102f9600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610984565b6040518082815260200191505060405180910390f35b341561031a57600080fd5b610346600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919050506109cc565b005b341561035357600080fd5b61035b610a6c565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561039b578082015181840152602081019050610380565b50505050905090810190601f1680156103c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156103e157600080fd5b610416600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610aa5565b604051808215151515815260200191505060405180910390f35b341561043b57600080fd5b610486600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610c77565b6040518082815260200191505060405180910390f35b6040805190810160405280600981526020017f414a4320636861696e000000000000000000000000000000000000000000000081525081565b60008082148061056157506000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054145b151561056c57600080fd5b81600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b60035481565b600080600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415151561072057600080fd5b80831115151561072f57600080fd5b610780836000808873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610813836000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506108688382610cfe90919063ffffffff16565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040518082815260200191505060405180910390a360019150509392505050565b601281565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610a2857600080fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6040805190810160405280600381526020017f414a43000000000000000000000000000000000000000000000000000000000081525081565b60008073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610ae257600080fd5b610b33826000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610bc6826000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a36001905092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000828211151515610d0c57fe5b818303905092915050565b6000808284019050838110151515610d2b57fe5b80915050929150505600a165627a7a7230582043a3cd97586e182885676a8c6e6413be040c6f728b9763d794ecdbfff9a4b7c90029"
       
     }' 


```
if /send success , will return txHash;
query tx by txHash ( REST api /txs/{hash}), check the evm call return code ("code") and contract_address ("contract_address") in field "log"
 

```
    "logs": [
            {
                "msg_index": "0",
                "success": true,
                "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt\",\"evm_output\":\"\"}"
            }
        ],
````


## getBalance(free fee)

Ref [HRC-20 query for free fee](https://gitee.com/orientwalt/apidoc_2020/blob/master/%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3/HRC-20%E5%B8%B8%E7%94%A8%E6%8E%A5%E5%8F%A3.md)


get Minter balance
```
$balanceOfHex=`curl http://htdf2020-test01.orientwalt.cn:1317/hs/contract/htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt/70a0823100000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9`

$echo $balanceOfHex                            ### (hex)
"000000000000000000000000000000000000000000a49be39dc14cb827000000"

$python -c "print int(${balanceOfHex}, 16)"  ###（decimal）
199000000000000000000000000

```

get reveiver balance
```

$balanceOfHex=`curl http://htdf2020-test01.orientwalt.cn:1317/hs/contract/htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt/70a0823100000000000000000000000066E0F9D1677F71D0CA6C1D5653673FE6392050E2`

$echo $balanceOfHex                            ### (hex)
"0000000000000000000000000000000000000000000000000000000000000000"

$python -c "print int(${balanceOfHex}, 16)"  ###（decimal）
0


```


## contract method:balances
#### get the minter balance


```

# 发交易  send;           打开合约
curl --location --request POST 'http://htdf2020-test01.orientwalt.cn:1317/hs/send' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-raw '    { "base_req": 
      { "from": "htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3", 
        "memo": "",
        "password": "12345678", 
        "chain_id": "testchain", 
        "account_number": "0", 
        "sequence": "0", 
        "gas_wanted": "500000", 
        "gas_price": "100", 
        "simulate": false
      },          
      "amount": [ 
              { "denom": "htdf", 
                "amount": "0" } ],
      "to": "htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt",
      "data": "70a0823100000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9"
    }'
```

if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "logs": [
            {
                "msg_index": "0",
                "success": true,
                "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000000000000000000000a49be39dc14cb827000000\"}"
            }
        ],
````

#### get the receiver balance

```

# 发交易  send;           打开合约
curl --location --request POST 'http://htdf2020-test01.orientwalt.cn:1317/hs/send' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-raw '    { "base_req": 
      { "from": "htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3", 
        "memo": "",
        "password": "12345678", 
        "chain_id": "testchain", 
        "account_number": "0", 
        "sequence": "0", 
        "gas_wanted": "500000", 
        "gas_price": "100", 
        "simulate": false
      },          
      "amount": [ 
              { "denom": "htdf", 
                "amount": "0" } ],
      "to": "htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt",
      "data": "70a0823100000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2"
    }'

```


if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "logs": [
            {
                "msg_index": "0",
                "success": true,
                "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"0000000000000000000000000000000000000000000000000000000000000000\"}"
            }
        ],
```



## contract method:transfer


```

# 发交易  send;           打开合约
curl --location --request POST 'http://htdf2020-test01.orientwalt.cn:1317/hs/send' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-raw '    { "base_req": 
      { "from": "htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3", 
        "memo": "",
        "password": "12345678", 
        "chain_id": "testchain", 
        "account_number": "0", 
        "sequence": "0", 
        "gas_wanted": "500000", 
        "gas_price": "100", 
        "simulate": false
      },          
      "amount": [ 
              { "denom": "htdf", 
                "amount": "0" } ],
      "to": "htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt",
      "data": "a9059cbb00000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2000000000000000000000000000000000000000000000000000000000000001e"
    }'
    
```


if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "logs": [
            {
                "msg_index": "0",
                "success": true,
                "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"0000000000000000000000000000000000000000000000000000000000000001\"}"
            }
        ],
```

## contract method:balances
after send , get the minter and the receiver balance, again
the get balance curl same like above

we can find that, the minter and the receiver balanc ,has change

HRC-20 query for free fee
```

$balanceOfHex=`curl http://htdf2020-test01.orientwalt.cn:1317/hs/contract/htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt/70a0823100000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9`

$echo $balanceOfHex                            ### (hex)
"000000000000000000000000000000000000000000a49be39dc14cb826ffffe2"

$python -c "print int(${balanceOfHex}, 16)"  ###（decimal）
198999999999999999999999970
```

```
$balanceOfHex=`curl http://htdf2020-test01.orientwalt.cn:1317/hs/contract/htdf1j0yrar88jgw0whzakxefy6t3ru3qxn4act3qvt/70a0823100000000000000000000000066E0F9D1677F71D0CA6C1D5653673FE6392050E2`

$echo $balanceOfHex                            ### (hex)
"000000000000000000000000000000000000000000000000000000000000001e"

$python -c "print int(${balanceOfHex}, 16)"  ###（decimal）
30


```


tx query

```
    "logs": [
            {
                "msg_index": "0",
                "success": true,
                "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000000000000000000000a49be39dc14cb826ffffe2\"}"
            }
        ],
````


```
    "logs": [
            {
                "msg_index": "0",
                "success": true,
                "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000000000000000000000000000000000000000001e\"}"
            }
        ],
````

