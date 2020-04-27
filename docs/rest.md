### account rest
```bash
[newaccount]
curl -X POST "http://localhost:1317/accounts/newaccount" -H "accept: application/json" -d "{\"password\": \"12345678\"}"
```
```json
{"address": "htdf1h290f6kfjwjexudqtp7hujm52c86mf8q5vush5"}
```
[get accountlist]
```bash
curl -X GET "http://localhost:1317/accounts/list" -H "accept: application/json"
Account #0: {htdf1yt0q9rdypm6zw83tm7r58etglsvgxzz6rymz0w}
Account #1: {htdf1yysruaystrfuxuxfdsqjxa0shvzts27p8l2r2x}
```

[get account information]
```bash
curl -X GET "http://localhost:1317/auth/accounts/htdf14yz30q7fqkvkks3wnmdm3sxkafugguevul44ga" -H "accept: application/json"
```
```json
{
	"type": "auth/Account",
	"value": {
		"address": "htdf14yz30q7fqkvkks3wnmdm3sxkafugguevul44ga",
		"coins": [
			{
				"denom": "htdf",
				"amount": "1000"
			}
		],
		"public_key": null,
		"account_number": "0",
		"sequence": "0"
	}
}
```
[getbalance]
```bash
curl -X GET "http://localhost:1317/bank/balances/htdf14yz30q7fqkvkks3wnmdm3sxkafugguevul44ga" -H "accept: application/json"
```
```json
[
	{
		"denom": "htdf",
		"amount": "1000"
	}
]
```
### transaction rest
```bash
[send transaction]
curl -X POST "http://localhost:1317/hs/send" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"htdf1njv34aldy8nn90jjqursjvvyvgmk38ez6hwpne\", \"memo\": \"Sent via Cosmos Voyager \",\"password\": \"12345678\", \"chain_id\": \"testchain\", \"account_number\": \"0\", \"sequence\": \"0\", \"gas_wanted\": \"200000\", \"gas_price\": \"100\", \"simulate\": false }, \"amount\": [ { \"denom\": \"htdf\", \"amount\": \"0.1\" } ],\"to\": \"htdf1xxe7xd28zf4njuszp6m5hlut5mvlyna8pvdwf6\"}"
```
```json
{
	"height": "119",
	"txhash": "02A61744D89A14E9C01C9B08B74EFADD6FE9DB9A625EBF0D4D936D1D765B7684",
	"log": "[{\"msg_index\":\"0\",\"success\":true,\"log\":\"\"}]",
	"gas_wanted": "200000",
	"gas_used": "28327",
	"tags": [
		{
			"key": "action",
			"value": "send"
		},
		{
			"key": "sender",
			"value": "htdf14yz30q7fqkvkks3wnmdm3sxkafugguevul44ga"
		},
		{
			"key": "recipient",
			"value": "htdf1h290f6kfjwjexudqtp7hujm52c86mf8q5vush5"
		}
	]
}
```
[create contraction]
```bash
curl -X POST "http://localhost:1317/hs/send" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"htdf188ptmpj3rvthmtd5af2ajvyxg9qarkdf69kmzr\", \"memo\": \"Sent via Cosmos Voyager \",\"password\": \"12345678\", \"chain_id\": \"testchain\", \"account_number\": \"0\", \"sequence\": \"2\", \"gas_wanted\": \"1200000\", \"gas_price\": \"100\", \"gas_adjustment\": \"1.2\", \"simulate\": false }, \"amount\": [ { \"denom\": \"htdf\", \"amount\": \"0.1\" } ],\"to\": \"\",\"data\": \"6060604052341561000f57600080fd5b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061042d8061005e6000396000f300606060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063075461721461006757806327e235e3146100bc57806340c10f1914610109578063d0679d341461014b575b600080fd5b341561007257600080fd5b61007a61018d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156100c757600080fd5b6100f3600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919050506101b2565b6040518082815260200191505060405180910390f35b341561011457600080fd5b610149600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506101ca565b005b341561015657600080fd5b61018b600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610277565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60016020528060005260406000206000915090505481565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561022557610273565b80600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b5050565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156102c3576103fd565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555080600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055507f3990db2d31862302a685e8086b5755072a6e2b5b780af1ee81ece35ee3cd3345338383604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a15b50505600a165627a7a7230582025e341a800f5478ed9b8aa0ee7a05d1165c779df9fd2479f9efaabdd937329b50029\",\"encode\":true}"
```
```json
{
	"height": "0",
	"txhash": "D24001686BC12B31CAB5A87AF470266E056045FF0D1C013A556955BA6D885EDE"
}	
```
[create raw transaction]
```bash
curl -X POST "http://localhost:1317/hs/create" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"htdf103x7taejyqwxrvyadu2yxd7u04wdqs5stq5a40\", \"memo\": \"Sent via Cosmos Voyager \",\"password\": \"\", \"chain_id\": \"testchain\", \"account_number\": \"3\", \"sequence\": \"3\", \"gas_wanted\": \"30000\", \"gas_price\": \"100\", \"gas_adjustment\": \"1.2\", \"simulate\": false }, \"amount\": [ { \"denom\": \"htdf\", \"amount\": \"0.1\" } ],\"to\": \"htdf1ec5yff9km0tlaemmuz6lk5zftkjv44hztjtfnc\",\"encode\":true}"
7b2274797065223a22617574682f5374645478222c2276616c7565223a7b226d7367223a5b7b2274797065223a2268746466736572766963652f73656e64222c2276616c7565223a7b2246726f6d223a22757364703134797a3330713766716b766b6b7333776e6d646d3373786b6166756767756576756c34346761222c22546f223a2275736470316832393066366b666a776a657875647174703768756a6d35326338366d663871357675736835222c22416d6f756e74223a5b7b2264656e6f6d223a2268746466222c22616d6f756e74223a223130227d5d7d7d5d2c22666565223a7b22616d6f756e74223a5b7b2264656e6f6d223a2268746466222c22616d6f756e74223a223130227d5d2c22676173223a22323030303030227d2c227369676e617475726573223a6e756c6c2c226d656d6f223a2253656e742076696120436f736d6f7320566f7961676572227d7d
```
[sign raw transaction]
```bash
curl -X POST "http://localhost:1317/hs/sign" -H "accept: application/json" -H "Content-Type: application/json" -d "{\"tx\":\"xxx\", \"passphrase\":\"12345678\",\"offline\":false,\"encode\":true}"
7b0a20202274797065223a2022617574682f5374645478222c0a20202276616c7565223a207b0a20202020226d7367223a205b0a2020202020207b0a20202020202020202274797065223a202268746466736572766963652f73656e64222c0a20202020202020202276616c7565223a207b0a202020202020202020202246726f6d223a2022757364703134797a3330713766716b766b6b7333776e6d646d3373786b6166756767756576756c34346761222c0a2020202020202020202022546f223a202275736470316832393066366b666a776a657875647174703768756a6d35326338366d663871357675736835222c0a2020202020202020202022416d6f756e74223a205b0a2020202020202020202020207b0a20202020202020202020202020202264656e6f6d223a202268746466222c0a202020202020202020202020202022616d6f756e74223a20223130220a2020202020202020202020207d0a202020202020202020205d0a20202020202020207d0a2020202020207d0a202020205d2c0a2020202022666565223a207b0a20202020202022616d6f756e74223a205b0a20202020202020207b0a202020202020202020202264656e6f6d223a202268746466222c0a2020202020202020202022616d6f756e74223a20223130220a20202020202020207d0a2020202020205d2c0a20202020202022676173223a2022323030303030220a202020207d2c0a20202020227369676e617475726573223a205b0a2020202020207b0a2020202020202020227075625f6b6579223a207b0a202020202020202020202274797065223a202274656e6465726d696e742f5075624b6579536563703235366b31222c0a202020202020202020202276616c7565223a2022413759754c354c316571624d77704c65374e364f6d667a485169773257583344323478467765434150793349220a20202020202020207d2c0a2020202020202020227369676e6174757265223a202277504a394666612b467a4f4b3769736c335a55354f6764317852766f7061455361594b42666230706748746858474251756541515a764b637877425967635438767a72454c326754667335417a7259305849357a49773d3d220a2020202020207d0a202020205d2c0a20202020226d656d6f223a202253656e742076696120436f736d6f7320566f7961676572220a20207d0a7d
```
[broadcast raw transaction]
```bash
curl -X POST "http://localhost:1317/hs/broadcast" -H "accept: aplication/json" -H "Content-Type: application/json" -d "{\"tx\":\"xxxx\"}"
```
```json
{
	"height": "454",
	"txhash": "860A8B5D919C36F52437339D7424C4ECF40B3B62D12A2A2B9129DF4B7698D511",
	"log": "[{\"msg_index\":\"0\",\"success\":true,\"log\":\"\"}]",
	"gas_wanted": "200000",
	"gas_used": "25660",
	"tags": [
		{
			"key": "action",
			"value": "send"
		},
		{
			"key": "sender",
			"value": "htdf14yz30q7fqkvkks3wnmdm3sxkafugguevul44ga"
		},
		{
			"key": "recipient",
			"value": "htdf1h290f6kfjwjexudqtp7hujm52c86mf8q5vush5"
		}
	]
}
```
### query rest(block & transaction status check)

[getblock]
```bash
curl -X GET "http://localhost:1317/blocks/latest" -H "accept: application/json"
```
```json
{
	"block_meta": {
		"block_id": {
			"hash": "3A78C849E9FE2A617B334D1626B3B2542428F1322900C30C1872B37C8187C7BE",
			"parts": {
	"total": "1",
	"hash": "288E37E4A836241C16A405C4F44668DE882FCDA17F77FC5899CE44D704439BA4"
			}
		},
		"header": {
			"version": {
	"block": "10",
	"app": "0"
			},
			"chain_id": "testchain",
			"height": "988",
			"time": "2019-04-02T06:15:50.977129339Z",
			"num_txs": "0",
			"total_txs": "0",
			"last_block_id": {
	"hash": "6B52F57908929CCD9BF69887C53237D31D611078C68DA7F06BBDC4A14BDC35A2",
	"parts": {
		"total": "1",
		"hash": "59EC13C9AA319B7B4489F72BBC0EA9CC449441F54AA0AF9E35D03F4FB8C960F9"
	}
			},
			"last_commit_hash": "6BDA956C84C070DCFA1F06C46547CF2F5BF532AD45796BBFFD8B0F76854C897A",
			"data_hash": "",
			"validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"next_validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
			"app_hash": "F7936381CFB337C7828C8409EAFB131D647C399B36D8B617B2E72CB93118175B",
			"last_results_hash": "",
			"evidence_hash": "",
			"proposer_address": "BA9667E11EC53439B20E4D3F03D5677BE2218BA3"
		}
	},
	"block": {
		"header": {
			"version": {
	"block": "10",
	"app": "0"
			},
			"chain_id": "testchain",
			"height": "988",
			"time": "2019-04-02T06:15:50.977129339Z",
			"num_txs": "0",
			"total_txs": "0",
			"last_block_id": {
	"hash": "6B52F57908929CCD9BF69887C53237D31D611078C68DA7F06BBDC4A14BDC35A2",
	"parts": {
		"total": "1",
		"hash": "59EC13C9AA319B7B4489F72BBC0EA9CC449441F54AA0AF9E35D03F4FB8C960F9"
	}
			},
			"last_commit_hash": "6BDA956C84C070DCFA1F06C46547CF2F5BF532AD45796BBFFD8B0F76854C897A",
			"data_hash": "",
			"validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"next_validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
			"app_hash": "F7936381CFB337C7828C8409EAFB131D647C399B36D8B617B2E72CB93118175B",
			"last_results_hash": "",
			"evidence_hash": "",
			"proposer_address": "BA9667E11EC53439B20E4D3F03D5677BE2218BA3"
		},
		"data": {
			"txs": null
		},
		"evidence": {
			"evidence": null
		},
		"last_commit": {
			"block_id": {
	"hash": "6B52F57908929CCD9BF69887C53237D31D611078C68DA7F06BBDC4A14BDC35A2",
	"parts": {
		"total": "1",
		"hash": "59EC13C9AA319B7B4489F72BBC0EA9CC449441F54AA0AF9E35D03F4FB8C960F9"
	}
			},
			"precommits": [
	{
		"type": 2,
		"height": "987",
		"round": "0",
		"block_id": {
			"hash": "6B52F57908929CCD9BF69887C53237D31D611078C68DA7F06BBDC4A14BDC35A2",
			"parts": {
				"total": "1",
				"hash": "59EC13C9AA319B7B4489F72BBC0EA9CC449441F54AA0AF9E35D03F4FB8C960F9"
			}
		},
		"timestamp": "2019-04-02T06:15:50.977129339Z",
		"validator_address": "BA9667E11EC53439B20E4D3F03D5677BE2218BA3",
		"validator_index": "0",
		"signature": "oyLY9YFTKcphIyj7qGB7KIfpDTauD6opNH7HOM5i3snhaIP6ttSecVNJtEMp6miBD0Z0al8l57KTG44ZSwLjAw=="
	}
			]
		}
	}
}
```
[getblock at a certain height]
```bash
curl -X GET "http://localhost:1317/blocks/5" -H "accept: application/json"
```
```json
{
	"block_meta": {
		"block_id": {
			"hash": "560512C063D68301E35DE32DD23765F3224F8F3FD1449713212994F1065DFC5E",
			"parts": {
	"total": "1",
	"hash": "2A193505571DD11212777460B7DE99479621A24ADE5907DC0B4B71E2C54EBA15"
			}
		},
		"header": {
			"version": {
	"block": "10",
	"app": "0"
			},
			"chain_id": "testchain",
			"height": "5",
			"time": "2019-04-02T04:49:06.852977978Z",
			"num_txs": "0",
			"total_txs": "0",
			"last_block_id": {
	"hash": "16ACFCA4AA5BE81E225C3BE64118C2C4DBA063B31A3DD91CE7E0EAC6165D6432",
	"parts": {
		"total": "1",
		"hash": "909612ED3B9518B983E934A0E346E90F94F67B34E54CCD3269999F82707A32BC"
	}
			},
			"last_commit_hash": "F85E7AB07BF4D6DC6DFFA36DDD15B21CF41BC294C0007B9D24C418ACB7EA6648",
			"data_hash": "",
			"validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"next_validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
			"app_hash": "77C8556C0D3A8BD4AD0A423CF072F0CAEE6C0648AEC5CE944B9BC1A5DFF24CA2",
			"last_results_hash": "",
			"evidence_hash": "",
			"proposer_address": "BA9667E11EC53439B20E4D3F03D5677BE2218BA3"
		}
	},
	"block": {
		"header": {
			"version": {
	"block": "10",
	"app": "0"
			},
			"chain_id": "testchain",
			"height": "5",
			"time": "2019-04-02T04:49:06.852977978Z",
			"num_txs": "0",
			"total_txs": "0",
			"last_block_id": {
	"hash": "16ACFCA4AA5BE81E225C3BE64118C2C4DBA063B31A3DD91CE7E0EAC6165D6432",
	"parts": {
		"total": "1",
		"hash": "909612ED3B9518B983E934A0E346E90F94F67B34E54CCD3269999F82707A32BC"
	}
			},
			"last_commit_hash": "F85E7AB07BF4D6DC6DFFA36DDD15B21CF41BC294C0007B9D24C418ACB7EA6648",
			"data_hash": "",
			"validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"next_validators_hash": "86B3FB7CEBCFB09A511A365609D5BBDC62127085B5219BC01FFF99BD7FD541BB",
			"consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
			"app_hash": "77C8556C0D3A8BD4AD0A423CF072F0CAEE6C0648AEC5CE944B9BC1A5DFF24CA2",
			"last_results_hash": "",
			"evidence_hash": "",
			"proposer_address": "BA9667E11EC53439B20E4D3F03D5677BE2218BA3"
		},
		"data": {
			"txs": null
		},
		"evidence": {
			"evidence": null
		},
		"last_commit": {
			"block_id": {
	"hash": "16ACFCA4AA5BE81E225C3BE64118C2C4DBA063B31A3DD91CE7E0EAC6165D6432",
	"parts": {
		"total": "1",
		"hash": "909612ED3B9518B983E934A0E346E90F94F67B34E54CCD3269999F82707A32BC"
	}
			},
			"precommits": [
	{
		"type": 2,
		"height": "4",
		"round": "0",
		"block_id": {
			"hash": "16ACFCA4AA5BE81E225C3BE64118C2C4DBA063B31A3DD91CE7E0EAC6165D6432",
			"parts": {
				"total": "1",
				"hash": "909612ED3B9518B983E934A0E346E90F94F67B34E54CCD3269999F82707A32BC"
			}
		},
		"timestamp": "2019-04-02T04:49:06.852977978Z",
		"validator_address": "BA9667E11EC53439B20E4D3F03D5677BE2218BA3",
		"validator_index": "0",
		"signature": "M1/C2iwsLCR73SclF1kSUGqsufdAvpy/eKMyyvUPyWCFYn02YwrmrsRgdKHYCewsbrlkGCmkEt4SIRMLoZMICA=="
	}
			]
		}
	}
}
```
[block detail]
```bash
curl -X GET "http://localhost:1317/block_detail/315817" -H "accept: application/json"
```
```json
{
  "block_meta": {
    "block_id": {
      "hash": "6B2189D5E3DD1A119A69A547FBEB7807D03831A19ADB900EFCE0DDE0FAF7E652",
      "parts": {
        "total": "1",
        "hash": "AB1676A62DC6EABCBA2152C124B4DE3ACD3AA10EB6607CBE510E9522380633C5"
      }
    },
    "header": {
      "version": {
        "block": "10",
        "app": "0"
      },
      "chain_id": "testchain",
      "height": "315817",
      "time": "2020-04-26T09:06:56.189176818Z",
      "num_txs": "1",
      "total_txs": "511",
      "last_block_id": {
        "hash": "56F1B23B69D9AB55140BC349D97F86DAA8160AC95D783ADE6CD54167DB5C52E6",
        "parts": {
          "total": "1",
          "hash": "9FA7F151E2EDB5E4CE4AB642B3983C317FEAB43B2B36B54E5C451ADE5549F703"
        }
      },
      "last_commit_hash": "B6B47F5400F124E65763D23BBB241D0224717C636E9E11B0EBA4516DA15193D4",
      "data_hash": "9B742472E3A231E68F1BCBFB6CFF99F5DEE982821D4260A93FF3C1C2320E4DB3",
      "validators_hash": "9EC932179412D65965494FB1E48C06C231F96EE53069F5F35FD20E87232B4F28",
      "next_validators_hash": "9EC932179412D65965494FB1E48C06C231F96EE53069F5F35FD20E87232B4F28",
      "consensus_hash": "4A3999F2F89E2E3AB1198825F730FCAE48D21406F7B9E36813A76B9A9B5ED509",
      "app_hash": "CE06F70423841E087275D6C6D377209DBFBFF2759ECF819E82DACAD8DFA475C3",
      "last_results_hash": "",
      "evidence_hash": "",
      "proposer_address": "443E6FC1949174C568E6C5103044EB40307E424F"
    }
  },
  "block": {
    "txs": [
      {
        "From": "htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3",
        "To": "",
        "Amount": [
          {
            "denom": "htdf",
            "amount": "0"
          }
        ],
        "Hash": "b275bf0b0f8e379cb1f7149c19485ae54fa3965d1c47e7e03560c4d90c042459",
        "Memo": "contract test",
        "Data": "60606040526000600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550341561005157600080fd5b6a52b7d2dcc80cd2e40000006003819055506003546000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610d61806100b76000396000f3006060604052600436106100af576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100b4578063095ea7b31461014257806318160ddd1461019c57806323b872dd146101c5578063313ce5671461023e5780634d853ee51461026d57806370a08231146102c257806393c32e061461030f57806395d89b4114610348578063a9059cbb146103d6578063dd62ed3e14610430575b600080fd5b34156100bf57600080fd5b6100c761049c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101075780820151818401526020810190506100ec565b50505050905090810190601f1680156101345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561014d57600080fd5b610182600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506104d5565b604051808215151515815260200191505060405180910390f35b34156101a757600080fd5b6101af61065c565b6040518082815260200191505060405180910390f35b34156101d057600080fd5b610224600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610662565b604051808215151515815260200191505060405180910390f35b341561024957600080fd5b610251610959565b604051808260ff1660ff16815260200191505060405180910390f35b341561027857600080fd5b61028061095e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156102cd57600080fd5b6102f9600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610984565b6040518082815260200191505060405180910390f35b341561031a57600080fd5b610346600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919050506109cc565b005b341561035357600080fd5b61035b610a6c565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561039b578082015181840152602081019050610380565b50505050905090810190601f1680156103c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156103e157600080fd5b610416600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610aa5565b604051808215151515815260200191505060405180910390f35b341561043b57600080fd5b610486600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610c77565b6040518082815260200191505060405180910390f35b6040805190810160405280600981526020017f42574320636861696e000000000000000000000000000000000000000000000081525081565b60008082148061056157506000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054145b151561056c57600080fd5b81600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b60035481565b600080600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415151561072057600080fd5b80831115151561072f57600080fd5b610780836000808873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610813836000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506108688382610cfe90919063ffffffff16565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040518082815260200191505060405180910390a360019150509392505050565b601281565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610a2857600080fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6040805190810160405280600381526020017f425743000000000000000000000000000000000000000000000000000000000081525081565b60008073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610ae257600080fd5b610b33826000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610bc6826000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a36001905092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000828211151515610d0c57fe5b818303905092915050565b6000808284019050838110151515610d2b57fe5b80915050929150505600a165627a7a72305820773594eba6917aa10c59fc8d6d7525346b2bf44ab4d0a38883a251b322a6abb90029",
        "TxClassify": 1,
        "TypeName": "send"
      }
    ],
    "evidence": {
      "evidence": null
    },
    "last_commit": {
      "block_id": {
        "hash": "56F1B23B69D9AB55140BC349D97F86DAA8160AC95D783ADE6CD54167DB5C52E6",
        "parts": {
          "total": "1",
          "hash": "9FA7F151E2EDB5E4CE4AB642B3983C317FEAB43B2B36B54E5C451ADE5549F703"
        }
      },
      "precommits": [
        {
          "type": 2,
          "height": "315816",
          "round": "0",
          "block_id": {
            "hash": "",
            "parts": {
              "total": "0",
              "hash": ""
            }
          },
          "timestamp": "2020-04-26T09:06:56.289601681Z",
          "validator_address": "1892BA91342D1CE62A69E0B6FAD2F8D504D59684",
          "validator_index": "0",
          "signature": "AF9GabnTiuhx+ymLozUtpAA/muWpIEaez9JhrPx65JSEMNzdnOmhIMBTunYDvFcmemG5GQ98mL62TQJTCffiAw=="
        },
        {
          "type": 2,
          "height": "315816",
          "round": "0",
          "block_id": {
            "hash": "",
            "parts": {
              "total": "0",
              "hash": ""
            }
          },
          "timestamp": "2020-04-26T09:06:56.295354797Z",
          "validator_address": "381C7EE64848C3BCF81E7D25C33B06272D9845C3",
          "validator_index": "1",
          "signature": "erxEvx7eVRu63buP3TYqbMf4LvIML7iwwMp0rEZPCWUn3DxFJiQujBEWniXvq4qWY7hZXlz4wIY23B2j27UDDw=="
        },
        {
          "type": 2,
          "height": "315816",
          "round": "0",
          "block_id": {
            "hash": "56F1B23B69D9AB55140BC349D97F86DAA8160AC95D783ADE6CD54167DB5C52E6",
            "parts": {
              "total": "1",
              "hash": "9FA7F151E2EDB5E4CE4AB642B3983C317FEAB43B2B36B54E5C451ADE5549F703"
            }
          },
          "timestamp": "2020-04-26T09:06:56.189176818Z",
          "validator_address": "443E6FC1949174C568E6C5103044EB40307E424F",
          "validator_index": "2",
          "signature": "CiaoI7Zq2prj0ZcacuASV+cncb7UPzbihSFVjwx7z+xjKC3WDh37WJUP8bRICnRjWmhJJurFrY0dfriOQBPIBA=="
        },
        {
          "type": 2,
          "height": "315816",
          "round": "0",
          "block_id": {
            "hash": "",
            "parts": {
              "total": "0",
              "hash": ""
            }
          },
          "timestamp": "2020-04-26T09:06:56.303408536Z",
          "validator_address": "62AF07AD3B9D90E870747959A299A702D63909AD",
          "validator_index": "3",
          "signature": "qhYhJDeqDVtf+51oMmb0isDSGjaIwehbBv+GGwzYYFHQvaXY9t/TLuohYv3nXsoJXdoIhib/CatNShqDyXhTDQ=="
        },
        {
          "type": 2,
          "height": "315816",
          "round": "0",
          "block_id": {
            "hash": "",
            "parts": {
              "total": "0",
              "hash": ""
            }
          },
          "timestamp": "2020-04-26T09:06:56.290126384Z",
          "validator_address": "8A36A1A3231B08F8B5D4911F314DD10428738755",
          "validator_index": "4",
          "signature": "Sys7ewxdy4LxccG75BcSlZ4X58+xPYZitbB9spu1HEAfkcQI944I24AWxKu3pyvXMRSgl0cb+F1EcIXkKjhWDw=="
        }
      ]
    }
  },
  "time": "2020-04-26 17:06:56"
}
```
[get tx by hash]
```bash
curl -X GET "http://localhost:1317/txs/02A61744D89A14E9C01C9B08B74EFADD6FE9DB9A625EBF0D4D936D1D765B7684" -H "accept: application/json"
```
```json
{
	"height": "119",
	"txhash": "02A61744D89A14E9C01C9B08B74EFADD6FE9DB9A625EBF0D4D936D1D765B7684",
	"log": "[{\"msg_index\":\"0\",\"success\":true,\"log\":\"\"}]",
	"gas_wanted": "200000",
	"gas_used": "28327",
	"tags": [
		{
			"key": "action",
			"value": "send"
		},
		{
			"key": "sender",
			"value": "htdf14yz30q7fqkvkks3wnmdm3sxkafugguevul44ga"
		},
		{
			"key": "recipient",
			"value": "htdf1h290f6kfjwjexudqtp7hujm52c86mf8q5vush5"
		}
	],
	"tx": {
		"type": "auth/StdTx",
		"value": {
			"msg": [
				{
					"type": "htdfservice/send",
					"value": {
						"From": "htdf14yz30q7fqkvkks3wnmdm3sxkafugguevul44ga",
						"To": "htdf1h290f6kfjwjexudqtp7hujm52c86mf8q5vush5",
						"Amount": [
							{
								"denom": "htdf",
								"amount": "10"
							}
						]
					}
				}
			],
			"fee": {
				"amount": [
					{
						"denom": "htdf",
						"amount": "10"
					}
				],
				"gas": "200000"
			},
			"signatures": [
				{
					"pub_key": {
						"type": "tendermint/PubKeySecp256k1",
						"value": "A7YuL5L1eqbMwpLe7N6OmfzHQiw2WX3D24xFweCAPy3I"
					},
					"signature": "or4Mb15p1WuENrVsNJ9KszlAqbiUMgPGM3UWoTrpPi5FU+p9xL5zCq57HLPjI0GTRJW345KJyOBrsQkHw/tS1w=="
				}
			],
			"memo": "Sent via Cosmos Voyager"
		}
	}
}
```
### node rest
```bash
[node_info]
curl -X GET "http://localhost:1317/node_info" -H "accept: application/json"
```
```json
{
	"protocol_version": {
		"p2p": "7",
		"block": "10",
		"app": "0"
	},
	"id": "92d61e4bc16fd5962351e215d759ae318639d90e",
	"listen_addr": "tcp://0.0.0.0:26656",
	"network": "testchain",
	"version": "0.31.1",
	"channels": "4020212223303800",
	"moniker": "suva-ubuntu",
	"other": {
		"tx_index": "on",
		"rpc_address": "tcp://0.0.0.0:26657"
	}
}
```
### validator rest
```bash
[get latest validator]
curl -X GET "http://localhost:1317/validatorsets/latest" -H "accept: application/json"
```
```json
{
		"block_height": "541",
		"validators": [
			{
				"address": "htdfvalcons1ad7rzcehm76c6zn0d9e9wrdjymlylas8mgjer3",
				"pub_key": "htdfvalconspub1zcjduepqwvxldrg9ftnwwskst7n5u8p8lny8mffuxql3yrscf75pynwt505s0rt3xl",
				"proposer_priority": "0",
				"voting_power": "10"
			}
		]
	}
```

[get validator set certain height]
```bash
curl -X GET "http://localhost:1317/validatorsets/5" -H "accept: application/json"
```
```json
{
		"block_height": "5",
		"validators": [
			{
				"address": "htdfvalcons1ad7rzcehm76c6zn0d9e9wrdjymlylas8mgjer3",
				"pub_key": "htdfvalconspub1zcjduepqwvxldrg9ftnwwskst7n5u8p8lny8mffuxql3yrscf75pynwt505s0rt3xl",
				"proposer_priority": "0",
				"voting_power": "10"
			}
		]
	}
```
### reward
#### - block rewards
```bash
curl -X GET "http://localhost:1317/minting/rewards/2" -H "accept: application/json"
```
```json
{
  "Reward": "14467592"
}
```
#### - validator rewards
```bash
curl -X GET "http://localhost:1317/distribution/validators/htdfvaloper1hv3hgjq9qadlnsf38qsrmnjwua92gm80ql45g7" -H "accept: application/json"
```
```json
{"operator_address":"htdf1hv3hgjq9qadlnsf38qsrmnjwua92gm802xddl4","self_bond_rewards":[{"denom":"satoshi","amount":"193328.804218410200000000"}],"val_commission":[{"denom":"satoshi","amount":"25890302.770065504352646621"}]}
```
#### - delegator rewards
```bash
curl -X GET "http://localhost:1317/distribution/delegators/htdf10fjsnx05ewesqjlmy5pesxzwa2t7z4e6vvqxvj/rewards/htdfvaloper10fjsnx05ewesqjlmy5pesxzwa2t7z4e6x4clme" -H "accept: application/json"
```
```json
[
  {
    "denom": "satoshi",
    "amount": "12264151227.930000000000000000"
  }
]
```
```bash
curl -X GET "http://localhost:1317/distribution/validators/htdfvaloper1hv3hgjq9qadlnsf38qsrmnjwua92gm80ql45g7/rewards" -H "accept: application/json"
```
```json
[
  {
    "denom": "satoshi",
    "amount": "12478630044.132000000000000000"
  }
]
```
```bash
curl -X GET "http://localhost:1317/distribution/delegators/htdf10fjsnx05ewesqjlmy5pesxzwa2t7z4e6vvqxvj/rewards" -H "accept: application/json"
```
```json
[
  {
    "denom": "satoshi",
    "amount": "12564381357.234000000000000000"
  }
]

```
#### outstanding
```bash
curl -X GET "http://localhost:1317/distribution/validators/htdfvaloper1hv3hgjq9qadlnsf38qsrmnjwua92gm80ql45g7/outstanding_rewards" -H "accept: application/json"
```
```json
[
  {
    "denom": "satoshi",
    "amount": "480004718330.997907965555643778"
  }
]
```
### total provisions
```bash
curl -X GET "http://localhost:1317/minting/total-provisions" -H "accept: application/json"
```
```json
{
  "Provision": "6000005748108018"
}
```
### contract call
[call contract]
```bash
curl -X GET "http://localhost:1317/hs/contract/htdf12dvguqedrvgfrdl35hcgfmz4fz6rm6chrvf96g/70a0823100000000000000000000000027681ceb7de9bae3c5f7f10f81ff5106e2ca48a7" -H "accept: application/json"
```
[get contract code]
```bash
curl -X GET "http://localhost:1317/hs/contract/htdf1ks6vgnp25r2eaa9k70dmsp448wmrma8mnucrsz/0000" -H "accept: application/json"
``` 
