# access smart contract via RESTful interface 


## get byte code


```
$go run byte_code_sample.go ../../../tests/evm/erc20/AJC_sol_AJCToken.abi  ../../../tests/evm/erc20/AJC_sol_AJCToken.bin htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3


contractCode, create contract|Code=6060604052341561000f57600080fd5b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061042d8061005e6000396000f300606060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063075461721461006757806327e235e3146100bc57806340c10f1914610109578063d0679d341461014b575b600080fd5b341561007257600080fd5b61007a61018d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156100c757600080fd5b6100f3600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919050506101b2565b6040518082815260200191505060405180910390f35b341561011457600080fd5b610149600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506101ca565b005b341561015657600080fd5b61018b600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610277565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60016020528060005260406000206000915090505481565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561022557610273565b80600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b5050565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156102c3576103fd565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555080600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055507f3990db2d31862302a685e8086b5755072a6e2b5b780af1ee81ece35ee3cd3345338383604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a15b50505600a165627a7a72305820f3c54d8cf0c62d5295ef69e3fc795fa1886b4de4d3d58f50f83c70ed26b99d890029
contractCode, minter|Code=07546172
contractCode, mint|minterAddress=htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk|Code=40c10f19000000000000000000000000ffa01831ff418f30a2266012c84ce168c335b28d00000000000000000000000000000000000000000000000000000000000f4240
contractCode, send|testContractToAddress=htdf1vms0n5t80acapjnvr4t9xeelucujq58zml4kg2|Code=d0679d3400000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2000000000000000000000000000000000000000000000000000000000000001e
contractCode, get balance|testContractToAddress=htdf1vms0n5t80acapjnvr4t9xeelucujq58zml4kg2|Code=27e235e300000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2
contractCode, get balance|strMinterAddress=htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk|Code=27e235e3000000000000000000000000ffa01831ff418f30a2266012c84ce168c335b28d

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
$ curl http://127.0.0.1:1317/hs/send \
    -H 'Content-Type: application/json' \
    -X POST \
    --data '{
                "base_req": {
                    "from": "htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk",
                    "memo": "",
                    "password": "12345678",
                    "chain_id": "testchain",
                    "account_number": "0",
                    "sequence": "0",
                    "gas_wanted": "500000",
                    "gas_price": "100",
                    "simulate": false
                },
                "amount": [{
                    "denom": "htdf",
                    "amount": "0"
                }],
                "to": "",
                "data": "6060604052341561000f57600080fd5b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061042d8061005e6000396000f300606060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063075461721461006757806327e235e3146100bc57806340c10f1914610109578063d0679d341461014b575b600080fd5b341561007257600080fd5b61007a61018d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156100c757600080fd5b6100f3600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919050506101b2565b6040518082815260200191505060405180910390f35b341561011457600080fd5b610149600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506101ca565b005b341561015657600080fd5b61018b600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610277565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60016020528060005260406000206000915090505481565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561022557610273565b80600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b5050565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156102c3576103fd565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555080600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055507f3990db2d31862302a685e8086b5755072a6e2b5b780af1ee81ece35ee3cd3345338383604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a15b50505600a165627a7a72305820f3c54d8cf0c62d5295ef69e3fc795fa1886b4de4d3d58f50f83c70ed26b99d890029"
            }'
```
if /send success , will return txHash;
query tx by txHash ( REST api /txs/{hash}), check the evm call return code ("code") and contract_address ("contract_address") in field "log"
 

```
    "log": [
        {
            "msg_index": "0",
            "success": true,
            "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"htdf1nzekkd4ax38rma33023rytan0letpaf9km9p50\",\"evm_output\":\"\"}"
        }
    ],
````


## contract method: minter
  send tx to contractAddr

```

# 发交易  send;           打开合约
curl http://127.0.0.1:1317/hs/send \
    -H 'Content-Type: application/json' \
    -X POST \
    --data '{
                "base_req": {
                    "from": "htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk",
                    "memo": "",
                    "password": "12345678",
                    "chain_id": "testchain",
                    "account_number": "0",
                    "sequence": "0",
                    "gas_wanted": "500000",
                    "gas_price": "100",
                    "simulate": false
                },
                "amount": [{
                    "denom": "htdf",
                    "amount": "0"
                }],
                "to": "htdf1mwv9agmm9f2vy68av0hd52lgkqjflltl2tggf7",
                "data": "07546172"
            }'
```


if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "log": [
            {
                "msg_index": "0",
                "success": true,
                "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000ffa01831ff418f30a2266012c84ce168c335b28d\"}"
            }
        ],
````

## contract method:mint

```

# 发交易  send;           打开合约
curl http://127.0.0.1:1317/hs/send \
    -H 'Content-Type: application/json' \
    -X POST \
    --data '{
                "base_req": {
                    "from": "htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk",
                    "memo": "",
                    "password": "12345678",
                    "chain_id": "testchain",
                    "account_number": "0",
                    "sequence": "0",
                    "gas_wanted": "500000",
                    "gas_price": "100",
                    "simulate": false
                },
                "amount": [{
                    "denom": "htdf",
                    "amount": "0"
                }],
                "to": "htdf1mwv9agmm9f2vy68av0hd52lgkqjflltl2tggf7",
                "data": "40c10f19000000000000000000000000ffa01831ff418f30a2266012c84ce168c335b28d00000000000000000000000000000000000000000000000000000000000f4240"
            }'
```


if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "log": [
        {
            "msg_index": "0",
            "success": true,
            "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"\"}"
        }
    ],
````


## contract method:balances
#### get the minter balance



```

# 发交易  send;           打开合约
curl http://127.0.0.1:1317/hs/send \
    -H 'Content-Type: application/json' \
    -X POST \
    --data '{
                "base_req": {
                    "from": "htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk",
                    "memo": "",
                    "password": "12345678",
                    "chain_id": "testchain",
                    "account_number": "0",
                    "sequence": "0",
                    "gas_wanted": "500000",
                    "gas_price": "100",
                    "simulate": false
                },
                "amount": [{
                    "denom": "htdf",
                    "amount": "0"
                }],
                "to": "htdf1mwv9agmm9f2vy68av0hd52lgkqjflltl2tggf7",
                "data": "27e235e3000000000000000000000000ffa01831ff418f30a2266012c84ce168c335b28d"
            }'
```

if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "log": [
        {
            "msg_index": "0",
            "success": true,
            "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"00000000000000000000000000000000000000000000000000000000000f4240\"}"
        }
    ],
````

#### get the receiver balance

```

# 发交易  send;           打开合约
curl http://127.0.0.1:1317/hs/send \
    -H 'Content-Type: application/json' \
    -X POST \
    --data '{
                "base_req": {
                    "from": "htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk",
                    "memo": "",
                    "password": "12345678",
                    "chain_id": "testchain",
                    "account_number": "0",
                    "sequence": "0",
                    "gas_wanted": "500000",
                    "gas_price": "100",
                    "simulate": false
                },
                "amount": [{
                    "denom": "htdf",
                    "amount": "0"
                }],
                "to": "htdf1mwv9agmm9f2vy68av0hd52lgkqjflltl2tggf7",
                "data": "27e235e300000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2"
            }'   

```


if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "log": [
        {
            "msg_index": "0",
            "success": true,
            "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"0000000000000000000000000000000000000000000000000000000000000000\"}"
        }
    ],
```



## contract method:send


```

# 发交易  send;           打开合约
curl http://127.0.0.1:1317/hs/send \
    -H 'Content-Type: application/json' \
    -X POST \
    --data '{
                "base_req": {
                    "from": "htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk",
                    "memo": "",
                    "password": "12345678",
                    "chain_id": "testchain",
                    "account_number": "0",
                    "sequence": "0",
                    "gas_wanted": "900000",
                    "gas_price": "100",
                    "simulate": false
                },
                "amount": [{
                    "denom": "htdf",
                    "amount": "0"
                }],
                "to": "htdf1mwv9agmm9f2vy68av0hd52lgkqjflltl2tggf7",
                "data": "d0679d3400000000000000000000000066e0f9d1677f71d0ca6c1d5653673fe6392050e2000000000000000000000000000000000000000000000000000000000000001e"
            }'
    
```


if /send success , will return txHash;
query tx by txHash ( REST api /txs/{txHash}), check the evm call return code ("code") and evm_output ("evm_output") in field "log"
 

```
    "log": [
        {
            "msg_index": "0",
            "success": true,
            "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"\"}"
        }
    ],
```

## contract method:balances
after send , get the minter and the receiver balance, again
the get balance curl same like above  `contract method:balances`

we can find that, the minter and the receiver balanc ,has change


```
    "log": [
        {
            "msg_index": "0",
            "success": true,
            "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"00000000000000000000000000000000000000000000000000000000000f4222\"}"
        }
    ],
````


```
    "log": [
        {
            "msg_index": "0",
            "success": true,
            "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000000000000000000000000000000000000000001e\"}"
        }
    ],
````

