# Greenfield Cosmos-SDK

This repo is forked from [cosmos-sdk](https://github.com/cosmos/cosmos-sdk).

The Greenfield Block Chain leverages cosmos-sdk to fast build a dApp running with tendermint. Due to the many 
requirements of Greenfield blockchain that cannot be fully satisfied by cosmos-sdk at present, we have decided to fork 
the cosmos-sdk repo and add modules and features based on it.

## Key Features

1. **auth**. The address format of the Greenfield blockchain is fully compatible with BSC (and Ethereum). It accepts EIP712 transaction signing and verification. These enable the existing wallet infrastructure to interact with Greenfield at the beginning naturally.
2. **crosschain**. Cross-chain communication is the key foundation to allow the community to take advantage of the Greenfield and BNB Smart Chain dual chain structure..
3. **gashub**. As an application specific chain, Greenfield defines the gas fee of each transaction type instead of calculating gas according to the CPU and storage consumption.
4. **gov**. There are many system parameters to control the behavior of the Greenfield and its smart contract on BSC, e.g. gas price, cross-chain transfer fees. All these parameters will be determined by Greenfield Validator Set together through a proposal-vote process based on their staking. Such the process will be carried on cosmos sdk.
5. **oracle**. The bottom layer of cross-chain mechanism, which focuses on primitive communication package handling and verification.
6. **upgrade**. Seamless upgrade on Greenfield enable a client to sync blocks from genesis to the latest state.

## Quick Start
*Note*: Requires [Go 1.18+](https://go.dev/dl/)

```shell
## proto-all
make proto-all

## build from source
make build

## test
make test

## lint check 
make lint
```

See the [Cosmos Docs](https://cosmos.network/docs/) and [Getting started with the SDK](https://tutorials.cosmos.network/academy/1-what-is-cosmos/).

## Related Projects
- [Greenfield](https://github.com/bnb-chain/greenfield): the official greenfield blockchain client.
- [Greenfield-Tendermint](https://github.com/bnb-chain/greenfield-tendermint): the consensus layer of Greenfield blockchain.
- [Greenfield-Storage-Provider](https://github.com/bnb-chain/greenfield-storage-provider): the storage service infrastructures provided by either organizations or individuals.
- [Greenfield-Relayer](https://github.com/bnb-chain/greenfield-relayer): the service that relay cross chain package to both chains.
- [Greenfield-Cmd](https://github.com/bnb-chain/greenfield-cmd): the most powerful command line to interact with Greenfield system.
- [Awesome Cosmos](https://github.com/cosmos/awesome-cosmos): Collection of Cosmos related resources which also fits Greenfield.


## Contribution
Thank you for considering helping with the source code! We appreciate contributions from anyone on the internet, no 
matter how small the fix may be.

If you would like to contribute to Greenfield, please follow these steps: fork the project, make your changes, commit them, 
and send a pull request to the maintainers for review and merge into the main codebase. However, if you plan on submitting 
more complex changes, we recommend checking with the core developers first via GitHub issues (we will soon have a Discord channel)
to ensure that your changes align with the project's general philosophy. This can also help reduce the workload of both 
parties and streamline the review and merge process.

## Licence (pending)