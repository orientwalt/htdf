# HTDF 2.0

| 功能 |  HTDF 1.x  | HTDF 2.0 |
|----- | ----| ------------|
| 普通转账 | ✓ | ✓ |
| 委托挖矿 | ✓ | ✓|
| HRC20代币 | ✓| ✓|
| solidity版本 | 0.4.20 | 0.8.x |
| 智能合约Event/Log | ✗ |✓ |
| 事件订阅 | ✗ |✓  |
| Webscoket API| ✗ |✓  |
| 轻节点 |✗  | ✓|
| DApp生态支持 |✗  |✓ |
| 其他新功能| --| TODO:参考以太坊并配合web3.js |



## HTDF 1.x 链上数据迁移方案

### 为什么要迁移?

在升级底层Tendermint和Cosmos SDK的版本之后,集成了Ethermint的Event等功能, 必然存在和HTDF 1.x版本兼容性问题. 如果能够解决兼容性问题最好, 如果解决不了兼容问题或者根本不可能兼容,那么就需要一个备用方案. 这个备用方案就是将HTDF 1.x版本的链上数据迁移到 HTDF 2.0链上.

### 怎么迁移?

首先需要解决的是"迁移什么?"的问题.

迁移什么?

- 普通账户的HTDF余额
- 合约账户的代币余额
- 验证节点的委托金额


#### 普通账户的HTDF余额

- 将所有账户地址及余额导出
- 写入创世区块 (或者通过空投的方式进行,这需要考虑交易所等场景会不会出现重复充币的问题)


#### 合约账户的代币余额

- 从区块链浏览器数据库获取所有的HRC20合约相关的交易的地址
- 编写脚本获取其余额
- 空投代币

#### 验证节点的委托金额

需要解决的问题: 委托金额提取之后会导致某些节点的权重比例降低, 可能造成安全问题,

理想的状态: 如果HTDF2.0升级之后出现问题, HTDF1.x仍旧可以正常运行.


**方案A**

| HTDF 1.x     |    HTDF 2.0                           |
|--------------|-------------------------------------|
|  发布通知(交易所, 区块链浏览器, 超级节点)        |    ---                           |
|  备份(冻结)链数据 |    --- |
|  导出链数据  |    --- |
|   ---           | 将HTDF1.x的数据映射到2.0的创世区块   |
|   ---           | 启动HTDF 2.0                      |
|   ---           | 进入观察阶段                        |
|   ---           | 结束观察阶段                        |
|   ---           | 创建HRC20代币合约并进行空投 |
|   ---           | 正式开放                           |
|   ---           | 区块浏览器, 华特钱包进行调整恢复正常使用  |
|   ---           | 交易所进行调整(注意: 不能将空投交易当成充币交易), 重新开放充币提币  |
|   ---           | HTDF 2.0 升级完成        |





**方案B(弃用)**

- 通知用户自行提取委托金和收益(限时)
- 2.0升级之后需要用户重新委托

**方案C(弃用)**

- 通知用户提取委托收益
- 只对委托金额进行快照, 在HTDF2.0 直接加入账户的余额并写入创世区块,
- 2.0升级之后需要用户重新委托



## 关于"方案A"的具体操作的研究


关键点:
- 委托收益补发: 如何计算委托收益?


使用 `hsd export`导出链的状态是否可行?

导出的内容不包含智能合约的状态, 如果智能合约有HTDF余额,但是又没有被导出, 是否会导致链的状态不对(总金额)?




### 导出内容的分析


```json
{
  "genesis_time": "2021-03-08T03:11:11.739124578Z", // 创世时间
  "chain_id": "testchain",
  "consensus_params": {  // 共识参数
    "block": {
      "max_bytes": "4194304",
      "max_gas": "15000000",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age": "100000"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    }
  },
  "validators": [ // 验证节点
    {
      "address": "",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "nZ20J8AA4XOJe3QX9gBiO/6hZ3HRYH40iEQ7ebfiwKI="
      },
      "power": "100",
      "name": "mynode"
    }
  ],
  "app_hash": "",
  "app_state": { // 链的状态
    "accounts": [ // 账户详情
      {
        "address": "htdf1rg46k2m9fde2d2zg25qdctuxxxzzmj0zpgwevy", // 地址
        "coins": [
          "1999989997000000satoshi" // 余额 (格式已改为 sdk.Coins)
        ],
        "sequence_number": "2",
        "account_number": "0"
      },
      {
        "address": "htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml",
        "coins": [
          "1000000000000000satoshi"
        ],
        "sequence_number": "0",
        "account_number": "2"
      },
      {
        "address": "htdf1794ehp0atyezw64e4765qeu5f4cp2k78fa0kl3",
        "coins": [
          "3000000000000000satoshi"
        ],
        "sequence_number": "0",
        "account_number": "1"
      }
    ],
    "auth": { // auth模块相关参数
      "collected_fees": null,
      "params": {
        "gas_price_threshold": "6000000000000",
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "4000"
      }
    },
    "staking": { // staking相关参数
      "pool": { // 收益池
        "not_bonded_tokens": "5999990087446565",
        "bonded_tokens": "10000000000",
        "last_zero_block_height": "2",
        "amplitude_sine_function": "11998162",
        "cycle_sine_function": "1176"
      },
      "params": { // 参数
        "unbonding_time": "259200000000000",
        "max_validators": 50,
        "max_entries": 7,
        "bond_denom": "satoshi"
      },
      "last_total_power": "100",
      "last_validator_powers": [
        {
          "Address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "Power": "100"
        }
      ],
      "validators": [ // 验证节点
        {
          "operator_address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "consensus_pubkey": "htdfvalconspub1zcjduepqnkwmgf7qqrsh8ztmwstlvqrz80l2zem369s8udyggsahndlzcz3qxtucr5",
          "jailed": false,
          "status": 2,
          "tokens": "10000000000",
          "delegator_shares": "10000000000.000000000000000000",
          "description": {
            "moniker": "mynode",
            "identity": "",
            "website": "",
            "details": ""
          },
          "unbonding_height": "0",
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "rate": "0.100000000000000000", // 佣金率
            "max_rate": "0.200000000000000000",
            "max_change_rate": "0.010000000000000000",
            "update_time": "2021-03-08T03:11:11.739124578Z"
          },
          "min_self_delegation": "1"
        }
      ],
      "delegations": [ // 委托
        {
          "delegator_address": "htdf1rg46k2m9fde2d2zg25qdctuxxxzzmj0zpgwevy",
          "validator_address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "shares": "10000000000.000000000000000000",
          "status": false
        }
      ],
      "unbonding_delegations": null,
      "redelegations": null,
      "exported": true
    },
    "mint": {
      "minter": {
        "inflation": "0.014999999818193565", // 当前通胀
        "annual_provisions": "90000000000000.000000000000000000" // 准备金
      },
      "params": {
        "mint_denom": "satoshi",
        "inflation_rate_change": "0.130000000000000000", // 通胀率
        "inflation_max": "0.200000000000000000", // 最大通胀率
        "inflation_min": "0.070000000000000000", // 最小通胀率
        "goal_bonded": "0.670000000000000000", // 总奖励
        "blocks_per_year": "6311520" // 每年出块数
      }
    },
    "distr": { // 发行相关
      "fee_pool": { // 手续费
        "community_pool": [
          {
            "denom": "satoshi",
            "amount": "1808931.300000000000000000"
          }
        ]
      },
      "community_tax": "0.020000000000000000", // 社区税
      "base_proposer_reward": "0.010000000000000000",
      "bonus_proposer_reward": "0.040000000000000000",
      "withdraw_addr_enabled": true,
      "delegator_withdraw_infos": [],
      "previous_proposer": "htdfvalcons15ln6fpgth4y3gsvu39a52km29vs2vrgg347r5d",
      "outstanding_rewards": [ // 出块奖励
        {
          "validator_address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "outstanding_rewards": [
            {
              "denom": "satoshi",
              "amount": "88637633.700000000000000000"
            }
          ]
        }
      ],
      "validator_accumulated_commissions": [ // 累积佣金
        {
          "validator_address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "accumulated": [
            {
              "denom": "satoshi",
              "amount": "8863763.370000000000000000"
            }
          ]
        }
      ],
      "validator_historical_rewards": [ // 验证节点历史奖励
        {
          "validator_address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "period": "1",
          "rewards": {
            "cumulative_reward_ratio": null,
            "reference_count": 2
          }
        }
      ],
      "validator_current_rewards": [ // 验证节点当前奖励
        {
          "validator_address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "rewards": {
            "rewards": [
              {
                "denom": "satoshi",
                "amount": "79773870.330000000000000000"
              }
            ],
            "period": "2"
          }
        }
      ],
      "delegator_starting_infos": [ // 委托开始信息
        {
          "delegator_address": "htdf1rg46k2m9fde2d2zg25qdctuxxxzzmj0zpgwevy",
          "validator_address": "htdfvaloper1rg46k2m9fde2d2zg25qdctuxxxzzmj0zt3kqm0",
          "starting_info": {
            "previous_period": "1",
            "stake": "10000000000.000000000000000000",
            "height": "0"
          }
        }
      ],
      "validator_slash_events": [] // 验证节点罚没事件
    },
    "gov": { // 链上治理相关
      "starting_proposal_id": "1",
      "deposits": null,
      "votes": null,
      "proposals": [],
      "deposit_params": {
        "min_deposit": [
          {
            "denom": "satoshi",
            "amount": "1000000000"
          }
        ],
        "max_deposit_period": "172800000000000"
      },
      "voting_params": {
        "voting_period": "172800000000000"
      },
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto": "0.334000000000000000",
        "penalty": "0.000000000000000000"
      }
    },
    "upgrade": { // upgrade 升级模块相关
      "GenesisVersion": {
        "UpgradeInfo": {
          "ProposalID": "0",
          "Protocol": {
            "version": "2",
            "software": "https://github.com/orientwalt/htdf/releases/tag/v1.3.1",
            "height": "1",
            "threshold": "0.900000000000000000"
          }
        },
        "Success": true
      }
    },
    "cirsis": { // 惩罚
      "constant_fee": {
        "denom": "satoshi",
        "amount": "1000"
      }
    },
    "slashing": { // 罚没参数
      "params": {
        "max_evidence_age": "7200",
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "double_sign_jail_duration": "0",
        "censorship_jail_duration": "0",
        "downtime_jail_duration": "600000000000",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000",
        "slash_fraction_censorship": "0.000000000000000000"
      },
      "signing_infos": {
        "htdfvalcons15ln6fpgth4y3gsvu39a52km29vs2vrgg347r5d": {
          "start_height": "0",
          "index_offset": "5",
          "jailed_until": "1970-01-01T00:00:00Z",
          "tombstoned": false,
          "missed_blocks_counter": "0"
        }
      },
      "missed_blocks": {
        "htdfvalcons15ln6fpgth4y3gsvu39a52km29vs2vrgg347r5d": []
      }
    },
    "service": { // 服务参数
      "params": {
        "max_request_timeout": "100",
        "min_deposit_multiple": "1000",
        "service_fee_tax": "0.010000000000000000",
        "slash_fraction": "0.001000000000000000",
        "complaint_retrospect": "1296000000000000",
        "arbitration_time_limit": "432000000000000",
        "tx_size_limit": "4000"
      }
    },
    "guardian": { // 超级管理员
      "profilers": [
        {
          "description": "genesis",
          "type": "Genesis",
          "address": "htdf1rg46k2m9fde2d2zg25qdctuxxxzzmj0zpgwevy",
          "added_by": "htdf1rg46k2m9fde2d2zg25qdctuxxxzzmj0zpgwevy"
        }
      ],
      "trustees": [
        {
          "description": "genesis",
          "type": "Genesis",
          "address": "htdf1rg46k2m9fde2d2zg25qdctuxxxzzmj0zpgwevy",
          "added_by": "htdf1rg46k2m9fde2d2zg25qdctuxxxzzmj0zpgwevy"
        }
      ]
    },
    "gentxs": null
  }
}
```



### 问题1: 使用导出的内容重新构造新链, 委托的金额和收益是否还在?

截至2021.3.10, 已经测试过数据迁移, 已迁移成功, 并成功使用junying/event_filter分支代码成功启动新链
委托的收益还在, 但是, 因junying/evnet_filter代码可能存在bug, 委托不能产生新的收益

已经提交issue: https://github.com/orientwalt/htdf/issues/41



查询委托收益:
```
curl  localhost:1317/distribution/delegators/htdf1g2jndvq2afh9dcecglp7gwmqzq347zgal0yqv5/rewards/htdfvaloper1g2jndvq2afh9dcecglp7gwmqzq347zga4kueml
```


## 测试方案A: 使用 hsd export 功能进行导出,然后重建新链

> 注意: 使用 junying/event_filter代码, 导入genesis.json时会报错, 参考:https://github.com/orientwalt/htdf/issues/39


```
查询委托收益
curl  localhost:1317/distribution/delegators/htdf1g2jndvq2afh9dcecglp7gwmqzq347zgal0yqv5/rewards/htdfvaloper1g2jndvq2afh9dcecglp7gwmqzq347zga4kueml

hsd  export > genesis.json

cp -R ~/.hsd ~/.hsd_bak_03_09_15

hsd unsafe-reset-all

修改genesis.json 中的 consensus_params.evidence
修改为
"evidence": {
       "max_age_num_blocks": "100000",
        "max_age_duration": "172800000000000"
}

genesis.json 添加  "initial_height": "新链起始高度",

mv genesis.json ~/.hsd/config/


修改 ~/.hsd/config/config.toml中的 db_backend
修改为
db_backend = "goleveldb"

hsd start


再次查询委托收益
curl  localhost:1317/distribution/delegators/htdf1g2jndvq2afh9dcecglp7gwmqzq347zgal0yqv5/rewards/htdfvaloper1g2jndvq2afh9dcecglp7gwmqzq347zga4kueml

查询账户余额
hscli query account htdf1g2jndvq2afh9dcecglp7gwmqzq347zgal0yqv5

提取委托收益

hscli tx distr withdraw-rewards  htdf1g2jndvq2afh9dcecglp7gwmqzq347zgal0yqv5  htdfvaloper1g2jndvq2afh9dcecglp7gwmqzq347zga4kueml --gas-price=100

查询交易
hscli query tx 5CC3BE1B2A28980308FE38A9712738264B2C2C2D6CD8F0959BD615A7BCADCB54


curl  localhost:1317/distribution/delegators/htdf1g2jndvq2afh9dcecglp7gwmqzq347zgal0yqv5/rewards/htdfvaloper1g2jndvq2afh9dcecglp7gwmqzq347zga4kueml

查询账户余额
hscli query account htdf1g2jndvq2afh9dcecglp7gwmqzq347zgal0yqv5


```




## 结论

从目前的情况来看, 将HTDF 1.x的数据迁移到 HTDF 2.0上最大几个已经**可行的**解决方案:

- 普通账户的HTDF余额: 导出到genesis文件
- 验证节点的委托金额: 导出到genesis文件
- 合约账户的代币余额: 主要是HRC20代币合约, 从区块链浏览器导出合约地址, 使用脚本获取每个账户的余额, 然后再HTDF2.0上进行部署, 然后, 再进行空投




## 附录A: Staking收益计算公式

Staking 收益计算公式
下面的计算公式基于主网治理参数.

年收益（忽略手续费收益和出块奖励）

- 年通胀 = 基数 * 通胀率 (即 9600000 * 13% = 12480000 HTDF )
- 验证人收益 = (年通胀 / 抵押总量) * (1 - 社区基金率) * (验证人自抵押 + 委托人抵押 * 佣金率)
- 委托人收益 = (年通胀 / 抵押总量) * (1 - 社区基金率) * 委托人自抵押 * (1 - 佣金率)

区块收益
- 区块通胀 = 年通胀 / (365*24*60*12) (约为 1.9787HTDF)
- 出块人额外奖励 = (BaseProposerReward + BonusProposerReward * PrecommitPower/TotalVotingPower) * (区块通胀 + 区块手续费收入)
- 区块总收益 = (区块通胀 + 区块手续费收入) * (1 - 社区基金率) - 出块人额外奖励
- 验证人总收益 =
    - 非出块人：(区块总收益 / 抵押总量) * 该验证人抵押总量
    - 出块人：非出块人收益 + 出块人额外奖励
- 佣金 = 验证人总收益 * 佣金率
- 验证人收益 = 验证人总收益 * (该验证人自抵押 / 该验证人抵押总量) + 佣金
- 委托人收益 = (验证人总收益 - 佣金) * (委托人自抵押 / 该验证人抵押总量)
