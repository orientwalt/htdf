# upgrade升级步骤


## 准备工作

> 注意:因提交升级提案后, 验证节点只有5分钟的时间进行投票, 所以需要提前做好准备工作
- 在这里下载最新发布版: https://gitee.com/orientwalt/htdf/releases/
- 选择几个验证节点,其总委托的占比超过`80%`, 最好能超过`85%`
- 将上一步选择的验证节点的私钥文件准备好, 最好放在同一机器上, 方便后续操作
- 给超级管理员地址转账`11HTDF`, 给每个投票的验证节点地址分别转账`1HTDF` 用于后面的提案和投票
- 将超级管理员(guardian)的私钥文件准备好
- 参考下文中的[操作指南](#操作指南)编写好命令, 形成文档, 方便后续操作

## 升级步骤
1. 停掉`hsd`(和`hscli`进程)
2. 备份数据目录 `.hsd`
3. 用准备工作时下载的最新`hsd`和`hscli`替换旧的`hsd`和`hscli`
4. 启动`hsd`(和`hscli`), 并使用`hsd version`和 `hscli version`核对版本,确保时可执行文件时最新的
5. 启动后,观察一段时间,10分钟到30分钟不等
6. 依次运行准备好的命令,执行升级
7. 检查升级结果
8. 升级完成


## 操作指南

- 配置hscli, 方便后续步骤的操作

    ```shell
    hscli config chain-id mainchain
    hscli config output json
    hscli config indent true
    hscli config trust-node true
    ```

- 超级管理员的地址
    ```shell
    # 获取genesis.json中的guardian地址
    GUARDIAN=$(cat ~/.hsd/config/genesis.json | jq .app_state.guardian.profilers[0].address | sed 's/"//g')
    ```

- 验证节点的地址
    ```shell
    VARLIDATOR_1=htdf1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    VARLIDATOR_2=htdf1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    VARLIDATOR_3=htdf1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    VARLIDATOR_4=htdf1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

    # 其他验证节点地址.....
    ```

- 给超级管理员转账`11HTDF`, 给每个投票的验证节点转账`1HTDF`
    ```shell

    hscli tx send $FROM_ADDR $GUARDIAN 1100000000satoshi  --gas-price=100

    hscli tx send $FROM_ADDR $VARLIDATOR_1 100000000satoshi  --gas-price=100
    hscli tx send $FROM_ADDR $VARLIDATOR_2 100000000satoshi  --gas-price=100
    hscli tx send $FROM_ADDR $VARLIDATOR_3 100000000satoshi  --gas-price=100
    hscli tx send $FROM_ADDR $VARLIDATOR_4 100000000satoshi  --gas-price=100

    # 其他...
    ```


- 升级的协议版本号
    |版本号|协议版本号|
    |------|-------|
    |v1.1.x|0|
    |v1.2.x|1|
    |v1.3.x|2|
    |... |... |
    |v1.n.x| n-1 | 

    例如, 如果本次升级是版本是`v1.3.1`,则协议版本如下:

    ```shell
    HTDF_VERSION=v1.3.1
    PROTOCOL_VERSION=2
    ```


- 获取最大的提案编号, 本次提案将在此基础上加1

    ```shell
    PROPOSAL_ID=$(expr $(hscli query gov proposals | jq | grep proposal_id | awk '{print $2}' | sed 's/[",]//g' | sort -nr | head -1) + 1)
    ```


- 获取最新高度, 升级高度在最新高度上加上100即可

    ```shell
    SWITCH_HEIGHT=$(expr $(curl -s http://localhost:1317/blocks/latest | jq .block_meta.header.height | sed 's/"//g') + 1 )
    ```

- 提交提案

    ```shell
    hscli tx gov submit-proposal $GUARDIAN \
    --gas-price=100  \
    --switch-height=$SWITCH_HEIGHT \
    --description="$HTDF_VERSION upgrade"\
    --title="$HTDF_VERSION"\
    --type="software_upgrade"\
    --deposit="1000000000satoshi"\
    --version="$PROTOCOL_VERSION"\
    --software="https://github.com/orientwalt/htdf/releases/tag/$HTDF_VERSION"
    ```

- 投票

    ```shell
    hscli tx gov vote  $VARLIDATOR_1  $PROPOSAL_ID yes --gas-price=100
    hscli tx gov vote  $VARLIDATOR_2  $PROPOSAL_ID yes --gas-price=100
    hscli tx gov vote  $VARLIDATOR_3  $PROPOSAL_ID yes --gas-price=100
    hscli tx gov vote  $VARLIDATOR_4  $PROPOSAL_ID yes --gas-price=100

    # 如果还有其他继续添加即可
    ```

- 查看upgrade信息

    ```shell
    # 查看提案的状态
    hscli query gov proposal $PROPOSAL_ID

    # 查看升级状态
    hscli query upgrade |jq
    ```
