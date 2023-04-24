# Null-Blockchain

The Null-Blockchain is considered an alternative chain-provider for core that acts like a dummy tendermint. The idea is that by removing all the config and setup needed for tendermint, a single zeta node can be started on its own that can process transactions and create blocks without needing to provide any consensus. Block-time is also frozen until enough transactions fill a block, or a call to a "backdoor" API is used to move time time forward. The aim is that this can be used as an internal tool to run simulations and aid testing.

## Using the Null-blockchain

Providing the following options when running `zeta node`, or by setting them in `config.toml` will start zeta with a null-blockchain:

```
--blockchain.chain-provider=nullchain
--blockchain.nullchain.log-level=debug
--blockchain.nullchain.transactions-per-block=10
--blockchain.nullchain.block-duration=1s
--blockchain.nullchain.genesis-file=PATH_TO_TM_GENESIS_FILE
```

`PATH_TO_TM_GENESIS_FILE` is required and can be a normal `genesis.json` that would be used with Tendermint. The Null-blockchain requires it to be able to parse and send the `app_state` to `InitChain`. Also if `genesis_time` is set it will be used as the initial frozen time of the chain, otherwise it will be set to `time.Now()`. 


## Moving Time Forward

There are two ways in which time can be moved forward. The first is by submitting a number of transactions equal to `transactions-per-block`. Once this threshold is hit the submitted transactions will be processed, and `zetatime` will be incremented by
`block-duration`.

The other is by using an exposed HTTP endpoint to specify either a duration, or a future datetime:

```
# By duration
curl -X POST -d "{\"forward\": \"10s\"}" http://localhost:3101/api/v1/forwardtime


# By datetime
curl -X POST -d "{\"forward\": \"2021-11-25T14:14:00Z\"}" http://localhost:3101/api/v1/forwardtime
```

Moving time forward will create empty blocks until the target time is reached. Any pending transactions will be processed in the first block. If the target time is such that it does not move ahead by a multiple of `block-duration` then time will be snapped backwards to the block last ended, and `zetatime` could be less than the target time. 

## Depositing Funds and Staking

The Null-Blockchain is made to work with as few external dependencies as possible and so does not dial into the Ethereum chain. This means that all assets being used must be built-in assets, and not ERC20. Funds can be deposited into the system using the faucet (See `zeta faucet --help`).

To be able to be able to flex goverance the null-blockchain will need to be able to pretend that a party has staked to allow voting and proposals to work. This is done with a mock up a staking account that loops itself into the collateral engine. To be able to simulate staking the faucet can be used to deposit the built-in asset `VOTE` into a party's general account in the collateral engine. This general account balance for `VOTE` is then sneakily looped into governance as if it were a staked balance.


## Replaying the chain and restoring from a snapshot

The Null-Blockchain contains functionality that allows saving to a file all the per-block-transactions in a way that it can be replayed. The following options control this behaviour:
```
--blockchain.nullchain.replay-file=SOME_PATH/replay-file
--blockchain.nullchain.replay
--blockchain.nullchain.record
```

where `--blockchain.nullchain.record` will turn on writing transaction data into `SOME_PATH/replay-file`, and `--blockchain.nullchain.replay` will turn on replaying the chain from that file.

The Null-Blockchain also supports loading from a snapshot. This means that with replay turned on, a stopped chain that is restarted can quickly return to its previous state in exactly the same way as a chain using Tendermint and the blockchain provider.
