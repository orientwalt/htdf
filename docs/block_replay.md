
## Replay last block

To resolve error `Wrong Block.Header.AppHash` we should have a approach to fix then chaindb. Firstly, you should use the right version `hsd`. Then, replay the bad block (latest block) by command `--replay-last-block`.

example:

```
hsd start --replay-last-block
```

## Reset(repair) app state to specified block height 

> Warning: Backup your node home before you reset.

Use the right version `hsd`, `hsd reset` reset the app state to specified height. This command only be used for repairring the chain db. It doesn't delete the any block. It's NOT ***rollback***.

example:

```shell
hsd reset --height 10000
```
