### General
#### 1. testnet & mainnet
       Q: What is the difference?
       A: same but chain id
#### 2. delegators vs. validators
   * [Q: What are delegators?](https://cosmos.network/docs/cosmos-hub/validators/validator-faq.html#what-is-a-delegator)

      A:  The Cosmos Hub is based on Tendermint, which relies on a set of validators to secure the network. The role of validators is to run a full-node and participate in consensus by broadcasting votes which contain cryptographic signatures signed by their private key. Validators commit new blocks in the blockchain and receive revenue in exchange for their work. They must also participate in governance by voting on proposals. Validators are weighted according to their total stake.

   * [Q: What are validators?](https://cosmos.network/docs/cosmos-hub/validators/validator-faq.html#what-is-a-validator)

      A: 
Delegators are Atom holders who cannot, or do not want to run a validator themselves. Atom holders can delegate Atoms to a validator and obtain a part of their revenue in exchange. Because they share revenue with their validators, delegators also share risks. Should a validator misbehave, each of their delegators will be partially slashed in proportion to their delegated stake. This is why delegators should perform due diligence on validators before delegating, as well as spreading their stake over multiple validators. Delegators play a critical role in the system, as they are responsible for choosing validators. Being a delegator is not a passive role: Delegators should actively monitor the actions of their validators and participate in governance.
### Account 
#### 1. nonce management
       Q: What is accoutnumber, accountsquence?
       A: squence for nonce, number is what?
          account number is auto incremented id. every account created on one chain is identified by unique ID. 
#### 2. addresses
       Q: What sorts of addresses? What are they?
       A: 
### Transaction
#### 1. TxHash
       Q: Where & How Generated?
       A:
#### 2. 
### Voting
#### 1. top 100~300 validators on cosmos hub
#### 2. abnomality detection/penality system
       Q: How to detect double-sign ?
       A: [unique priv_validator.json](https://hub.cosmos.network/join-testnet.html#reset-data)
#### [3. Sentry Nodes & DDoS Protecton](https://cosmos.network/docs/cosmos-hub/validators/security.html#sentry-nodes-ddos-protection)
       Q: What are Sentry Nodes?
       A: They act like avatars of the validator nodes,
          similar to reverse proxy of a webapp,
          concealing the real ip addresses of the validators, where network is divided into external and internal.
       Q: Usage
       A:  1. To create a full node as a sentry node.
           2. To edit $HOME/.hsd/config/config.toml so that fill "private_peer_ids" with a validator node ID.
              > private_peer_ids = "5a533005d74b40ab954b33029a5682ec8794d014"
              > hsd tendermint show_node_id
           3. On the validator node, disallow all incoming connections in the firewall.
          Only allow incoming connection on 46656/tcp from the internal IP of the sentry node.
           4. On the validator node, edit $HOME/.hsd/config/config.toml.
          Remove all existing peer information from persistent_peers and put the node information of
          your sentry node in the format of ID@IP:PORT here.
              > persistent_peers = "1ebc5ca705b3ae1c06a0888ff1287ada82149dc3@10.10.0.2:46656"
          And the validator node should set pex to false.
              > pex = false
          Fifth, Restart both nodes.    

  * [Reference](https://medium.com/forbole/a-step-by-step-guide-to-join-cosmos-hub-testnet-e591a3d2cb41)
#### [4. voting power](https://cosmos.network/docs/cosmos-hub/validators/validator-faq.html#general-concepts)
       Q: voting power? 
       A: It depends on stake amount. 1 voting power can be gained when 1000,000 stake is staked.
### GasMetering
#### 1. Context & GasMetering
       Q: How Gas Consumed/Calculated?(How gasmetering works when every read,write,delete,has,iterator?)
       A: KVStore->context.KVStore->GasKV.NewStore(GasMetering enabled)
          Its own gasmeter is detached to every KVStore including stateDB's KVStores.
          So consumed gas will be calculated for every database operations(Read, Write, Delete, Has,...)
          Every operation includes consumeGas/consumeSeekGas. It will deduct the corresponding gas to every operation
          folling the setting in GasConfig

          keywords: ctx.KVStore
       Q: Where to initialze (Tx)GasMeter & BlockGasMeter?
       A: BlockGasMeter: baseapp.go/BeginBlock
          (Tx)GasMeter:  ante.go/AnteHanlder/SetGasMeter
#### 2. GasMeter vs. BlockGasMeter
       Q: Difference between InfiniteGasMeter vs. BasicGasMeter
       A: InfiniteGasMeter has no limit while BasicGasMeter has limit
       Q: What is BlockGasLimit?
       Q: How it works?
          --BlockGasMeter--
          * MaximumBlockGas From ConsensusParams's BlockParams(actually from genesis.json. default value is -1,that's,infiniteGasMeter)
          * ref: baseapp.go
          BeginBlock-->MaxGas Set
          runTx------->ConsumeGas

          --GasMeter
          runTx,runMsgs---->ConsumeGas for every kvstore handling
#### 3. cosmos vs. ethereum
       Q: structure analysis
       A:                 cosmos      ethereum
          block size     infinite    1,500,000
          avg tx limit  200,000      21,000, 53000+
          calc metrics  db handling  contract data size+content(evm)
### app vs. context
### version in cosmos
###  CheckTx vs. DeliverTx
####
