```
使用 master分支最新代码, 编译hsd, 将数据导出


hsd  export > genesis.json

cp -R ~/.hsd ~/.hsd_bak_03_09_15

新链只需要 .hsd/config/* (需要删除addrbook.json) 和 data/priv_validator_state.json

hsd unsafe-reset-all



修改genesis.json 中的 consensus_params.evidence
修改为
"evidence": {
       "max_age_num_blocks": "100000",
        "max_age_duration": "172800000000000"
}

genesis.json 添加  "initial_height": "新链起始高度",

在auth, service 下添加 initial_height , 最好和 initial_height保持一致

修改 upgrade 下的信息, version必须为0, height和 initial_height保持一致

修改accounts 下面的 original_vesting , 将null改为 []
修改accounts 下面的 delegate_free, 将null改为 []
修改accounts 下面的 delegate_vesting, 将null改为 []

对比 validators的佣金率 commission 和节点信息

mv genesis.json ~/.hsd/config/


修改 ~/.hsd/config/config.toml中的 db_backend
修改为
db_backend = "goleveldb"



启动, 如果没有panic就ok了
hsd start

```