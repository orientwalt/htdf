
## Replay last block


解决`Wrong Block.Header.AppHash`错误我们可以使用`--replay-last-block`命令对最新的区块进行`replay`.注意:要确认使用正确版本的`hsd`.

```
hsd start --replay-last-block
```

## Reset chain state to specified block height 

> 警告: 请先备份节点数据


使用正确版本的`hsd`,使用命令 `hsd reset`会将链的状态重置到指定高度然后进行replay这个区间内的所有区块.这个命令主要用于修复链的状态.不会删除任何已有区块. 即它不是所谓的"回滚".

例如: 重置链的状态到高度`10000`

```shell
hsd reset --height 10000
```
