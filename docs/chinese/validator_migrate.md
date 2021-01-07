# 验证节点服务器迁移

> 注意: 两个使用相同key的验证节点**不能**同时在线!

为了方便说明, 将已有的验证节点记为 node1 , 将新建的服务器上的验证节点记为 node2

1. 停止节点 node1
2. 复制 node1的下的2个目录: `~/.hscli` 和  `~/.hsd/config` 到 node2的`$HOME`下
3. 将node1中的可执行文件`hsd`和`hscli`复制到node2上, 在node2上运行命令 `hsd unsafe-reset-all`
4. 在node2上使用命令 `hsd start` 启动
5. 观察node2节点日志输出
6. 在node2上使用命令`hscli query staking validators` 查看node2节点是否被`jailed`; 如果被jailed则需要使用命令`hscli tx slashing unjail $VARLIDATOR`进行`unjail`
7. 在node2上运行命令 `curl -s http://192.168.0.171:1317/blocks/latest | grep validator_address` 观察node2验证节点地址是否存在
8. (可选)修改`$HOME/.hsd/config/config.toml`中的 `persistent_peers`, 将node1的IP为node2的IP. 方便起见, 可以待所有验证节点迁移成功后再修改, 所有验证节点可以使用同一份`config.toml` (前提是没有复杂的网络拓扑!). 最后, 将修改后的`config.toml`上传至gitee以便其他用户搭建新的节点.

