# Changelog

## v0.0.4
This release includes features and bug fixes, mainly:
1. Introduced `crosschain` and `oracle` module;
2. Updated EIP712 related functions;
3. Enabled public delegation upgrade

* [\#42](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/42) feat: implement cross chain and oracle modules
* [\#78](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/78) fix: update EIP712 related functions (#78)
* [\#87](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/87) feat: introduce enable public delegation upgrade

## v0.0.3

This release includes features and bug fixes, mainly:
1. Gashub module;
2. Customized upgrade module;
3. Customized Tendermint with vote pool;
4. Disable create validator operation;
5. EIP712 bug fix;

* [\#72](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/72) feat: add gashub module
* [\#79](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/79) fix: disable create validator operation
* [\#80](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/80) fix: EIP712 issue with uint8[] type
* [\#81](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/81) feat: update tendermint to enable validator updates and votepool
* [\#82](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/82) feat: custom upgrade module
* [\#83](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/83) ci: fix release flow

## v0.0.2
This is the first release of inscription-cosmos-SDK.

It includes two key features:

1. EIP721 signing schema support
2. New staking mechanism

* [\#38](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/38) ci: fix failed ci jobs
* [\#47](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/47) docs: change pull request template
* [\#36](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/36) feat: add support for eth address standard
* [\#68](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/68) fix: errors with EIP712 signature
* [\#71](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/71) test: fix the unstable UT bugs and remove useless testcase
* [\#70](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/70) fix: errors with EIP712 signature
* [\#46](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/46) feat: customize staking module for inscription
* [\#73](https://github.com/bnb-chain/inscription-cosmos-sdk/pull/73) feat: add bls key types into keyring management tool

## v0.0.1
Fork from cosmos-sdk 0.46.4
