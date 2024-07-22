# Changelog

## V1.9.0
This release introduces the Mongolian upgrade

## V1.8.0
This release introduces the Veld upgrade

## V1.7.0
This release introduces the Erdos upgrade

Features:
* [#417](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/417) feat: feat: add multi message support for greenfield crosschain app

Fixes:
* [#422](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/422) fix: fix the deps security


## V1.6.0
This release introduces the Serengeti upgrade.

Features:
* [#400](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/400) feat: add upgrade Serengeti


## V1.5.0
This release introduces the Pawnee upgrade.

Features:
* [#396](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/396) feat: introduce Pawnee upgrade

BUGFIX
* [#401](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/401) fix: sync failed from genesis


## V1.4.0
This release introduces the Ural upgrade.

Features:
* [#388](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/388) feat: introduce the Ural upgrade

BUGFIX
* [#391](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/391) fix: module consensus version is not correctly set when upgrade after state sync


## v1.3.0
This release introduces the Hulunbeier upgrade.

Features:
* [#368](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/368) feat: introduce the Hulunbeier upgrade
* [#376](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/376) feat: introduces hulunbeierPatch upgrade

## v1.2.2
This release introduce the Manchurian upgrade to the mainnet.

Features:
* [#385](https://github.com/bnb-chain/greenfield/pull/385) feat: introduce the Manchurian upgrade to the mainnet

## v1.2.1
This release changes the Manchurian upgrade height of testnet.

Features:
* [#374](https://github.com/bnb-chain/greenfield/pull/374) feat: modify Manchurian hardfork height of testnet

## V1.2.0
This release introduce the Manchurian upgrade to the testnet.

Chores:
* [#365](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/365) chore: add xml marshal for Int type
* [#364](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/364) chore: add xml marshal for UInt type

## v1.1.1
This release introduces the Pampas upgrade to the mainnet.  

Features:
* [#360](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/360) feat: enable Pampas hardfork to mainnet

## v1.1.0
This release introduces the Pampas upgrade and contains 8 new features.  

Features:
* [#316](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/316) feat: add validation with context information
* [#321](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/321) feat: add tool to migrate stores for fastnode
* [#336](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/336) feat: introduce the pampas upgrade
* [#341](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/341) feat: add support for some json rpc queries
* [#353](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/353) feat: distinguish inturn relayer
* [#320](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/320) feat: add MsgRejectMigrateBucket
* [#357](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/357) feat: support Secp256k1 format private keys import
* [#358](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/358) feat: enable Pampas hardfork to testnet

## v1.0.1
This release includes 1 bug fix.

* [#338](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/338) fix: count total when pagination request is empty

## v1.0.0
This release includes 2 new features.

* [#317](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/317) feat: add Nagqu fork to mainnet
* [#323](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/323) feat: support multi msg for `PrintEIP712MsgType` flag

## v0.2.6
This release caps the pagination limit for queries at 100 records if it exceeds the default pagination limit.

* [#315](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/315) fix: restrict pagination limit  

## v0.2.6-alpha.1
This release includes 1 new feature.

Features:
* [#310](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/310) feat: add FlagPrintEIP712MsgType  

Chores:
* [#312](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/312) chore: improve the validations of parameters

## v0.2.5
This release includes all the changes in the v0.2.5 alpha versions and 1 new feature.

Features:
* [#306](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/306) feat: enable Nagqu hardfork to testnet  
* [#277](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/277) feat: restrict token transfers to payment accounts

Chores:
* [#300](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/300) chore: add hardfork logic for Nagqu

## v0.2.5-alpha.1
This release includes 1 feature and 1 chore.

Features:
* [#277](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/277) feat: restrict token transfers to payment accounts

Chores:
* [#300](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/300) chore: add hardfork logic for Nagqu

## v0.2.4
This release includes all the changes in the v0.2.4 alpha versions and 2 new bugfixes.

Bugfixes:
* [#281](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/281) fix: disable pre deliver when raw db store is used
* [#287](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/287) fix: fix dependency security issues

Chores:
* [#289](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/289) chore: modify code comments in VerifySignature

## v0.2.4-alpha.2
This release includes 5 new features.

Features:
* [#266](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/266) feat: implement cross-chain mechanism between op and greenfield
* [#267](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/267) feat: add method to access check state in app
* [#269](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/269) feat: track store r/w consume for enabling plain store
* [#270](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/270) feat: make cross-chain token mintable
* [#273](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/273) feat: skip sig verification on genesis block

Chores:
* [#268](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/268) chore: modify default gas

## v0.2.4-alpha.1  
This release includes 4 features and 1 bugfix.

Features:
* [#256](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/256) feat: add Nagqu fork name for upcoming upgrading
* [#258](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/258) feat: add flag to enable or disable heavy queries
* [#263](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/263) feat: add MsgRenewGroupMember to renew storage group member  
* [#262](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/262) feat: add MsgUpdateStorageProviderStatus to update sp status  

Fix:  
* [#257](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/257) fix: add bls proof in proposal description

## v0.2.3
This is a official release for v0.2.3, includes all the changes since v0.2.2.

## v0.2.3-alpha.5
This is a maintenance release.

Features:
* [#253](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/253) feat: add option for disabling event emit

Fix:
* [#252](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/252) fix: limit pagination to protect the node would not be Query DoS

## v0.2.3-alpha.4
This is a maintenance release.

Features:
* [#247](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/247) feat: add UpdateChannelPermissions tx for crosschain module

Chores:
* [#246](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/246) chore: remove unused tools
* [#248](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/248) chore: implement base64 encoding in EIP712
* [#249](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/249) chore: adjust MsgSealObject gas to 120

## v0.2.3-alpha.2
This is a maintenance release.  

* [#231](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/231) feat: enable diff on iavl store 
* [#233](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/233) chore: bnb wording change 
* [#232](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/232) feat: using abi.encode for update param tx 
* [#222](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/222) feat: performance improvement 
* [#234](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/234) staking: create validator in one transaction 
* [#237](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/237) chore: update greenfield-cometbft-db version 
* [#238](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/238) feat: show proposal failed reason 
* [#239](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/239) feat: add bls verification 
* [#242](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/242) feat: support cross chain for multiple blockchains 

## v0.2.3-alpha.1
This release upgrades the reference cosmos-sdk to v0.47.3.
Please refer to the [changelogs of cosmos-sdk v0.47.3](https://github.com/cosmos/cosmos-sdk/blob/v0.47.3/CHANGELOG.md) for more details regarding the changes.

* [#220](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/220) feat: upgrade cosmos-sdk to v0.47.3
* [#214](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/214) chore: fix typo and update swagger
* [#219](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/219) fix: fix the security issues 
* [#218](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/218) fix: add sorting of EIP712 msg types
* [#224](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/224) feat: allows for setting a custom http client when NewClientFromNode 
* [#228](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/228) feat: optimize NewCustomClientFromNode 

## v0.2.2
This is a maintenance release. The changelog includes all the changes since v0.2.1.

* [#214](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/214) chore: fix typo and update swagger
* [\#210](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/210) feat: add msg in gashub
* [\#211](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/211) fix: fix blockchain stop to produce blocks

## v0.2.2-alpha.1
This is a maintenance release.  

* [#214](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/214) chore: fix typo and update swagger 

## v0.2.1-alpha.2
This is a maintenance release.

* [\#210](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/210) feat: add msg in gashub
* [\#211](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/211) fix: fix blockchain stop to produce blocks

## v0.2.1-alpha.1
This release enable the support of multiple messages for EIP712.

* [\#205](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/205) feat: support multiple messages for EIP712
* [\#206](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/206) fix: fix potential panic in tx simulation

## v0.2.1
This release is a hot fix release for v0.2.0.

* [\#203](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/203) fix: update DefaultMaxTxSize and gas simulation logic
* [\#204](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/204) fix: allow GasParams fixedtype's gas is zero

## v0.2.0
This release upgrades the reference cosmos-sdk to v0.47.2. As the cosmos-sdk v0.47.2 is a huge breaking upgrade, 
we decide to cherry-pick the recent contributed commits and apply to the v0.47.2. The commit history previous 
releases are archived and the branch is backed up in the `master_back` branch.

Notable breaking changes:
1. The previous keyring is replaced by the new keyring. Please regenerate them, otherwise you will get error.

Other notable changes, please refer to the [changelogs of cosmos-sdk v0.47.1](https://github.com/cosmos/cosmos-sdk/blob/v0.47.1/CHANGELOG.md)
and [changelogs of cosmos-sdk v0.47.2](https://github.com/cosmos/cosmos-sdk/blob/v0.47.2/CHANGELOG.md).

## v0.1.1
This is a maintenance release.

* [\#187](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/187) fix: fix validator update logic

## v0.1.0
This is a maintenance release.

* [\#163](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/163) fix: update DefaultMaxTxSize and gas simulation logic
* [\#155](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/155) feat: add gas config for discontinue object message

## v0.0.14
This is a hotfix release.

* [\#153](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/153) fix: fee calculation in oracle

## v0.0.13
This is a hotfix release.

* [\#150](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/150) fix: revert KeyUpgradedClient key in upgrade module

## v0.0.12
This is a maintenance release.

* [\#147](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/147) chore: refine the default gas
* [\#146](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/146) chore: refine crosschain module
* [\#145](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/145) chore: refine oracle module
* [\#143](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/143) docs: add licence and disclaim
* [\#141](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/141) refine the storage tx fee
* [\#140](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/140) chore: refine the code of gashub module
* [\#142](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/142) refactor: refine cross-chain governance
* [\#144](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/144) chore: refine upgrade module for code quality
* [\#129](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/129) feat: add challenger address to validators

## v0.0.11
This is a maintenance release.

* [\#135](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/135) fix: data race issue
* [\#134](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/134) feat: add gas params for new messages in storage module
* [\#133](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/133) feat: a relayer can relay cross chain tx in batch
* [\#136](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/136) fix: refine the code of crosschain and oracle module


## v0.0.10
This release reverts the unneeded changes.

* [\#128](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/128) revert: revert the changes of the callbackGasPrice

## v0.0.9
This release fix the v0.0.8 dependencies.

* [\#126](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/126) feat: register new storage message to gashub

## v0.0.8
This release updates some default module params.

* [\#117](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/117) feat: add Bytes/SetBytes for Uint
* [\#120](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/120) feat: update the initial balance for the crosschain module
* [\#119](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/119) fix: keep address format the same with ethereum
* [\#121](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/121) feat(cross-chain): add callbackGasPrice to cross-chain package
* [\#124](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/124) fix: fix the crosschain keeper in params module
* [\#123](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/123) feat: add gas config for challenge module
* [\#118](https://github.com/bnb-chain/greenfield-cosmos-sdk/pull/118) fix: update gashub default params

## v0.0.7
This release add the support of cross chain governance.

* [\#110](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/110) feat: use proposal for cross chain parameter governance
* [\#114](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/114) feat: update min gas price in GasInfo
* [\#111](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/111) docs: add the readme of gnfd-cosmso-sdk
* [\#112](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/112) dep: update tendermint version
* [\#113](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/113) fix: remove the std print
* [\#105](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/105) feat: refactor gashub module
* [\#108](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/108) feat: split commands for withdrawal of rewards and commission
* [\#86](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/86) test: integration tests for creating and impeaching validator
* [\#106](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/106) fix: implemented fix for proto-gen-swagger error
* [\#107](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/107) fix: jail until a proper time

## v0.0.6
This release is a maintenance release.

* [\#102](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/102) feat: Add gas calculator for all the messages of storage module
* [\#101](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/101) fix: add missing registration for gashub query server
* [\#100](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/100) fix: fix impeach validator
* [\#99](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/99) fix: fix the params query for oracle and crosschain module
* [\#98](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/98) feat: add comments for the events of oracle and crosschain modules
* [\#97](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/97)  fix: changed validator jail timestamp to grpc compliant unit
* [\#96](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/96)  feat: add a default gas calculator
* [\#95](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/95)  fix: update balance query eth jsonrpc method
* [\#94](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/94)  fix: add gas calculator for msg create validator
* [\#90](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/90)  feat: add support for EVM jsonrpc

## v0.0.5
This release is for rebranding from inscription to greenfield, renaming is applied to all packages, files.
* [\#91](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/91) feat: rebrand from inscription to greenfield

## v0.0.4
This release includes features and bug fixes, mainly:
1. Introduced `crosschain` and `oracle` module;
2. Updated EIP712 related functions;
3. Enabled public delegation upgrade

* [\#42](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/42) feat: implement cross chain and oracle modules
* [\#78](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/78) fix: update EIP712 related functions (#78)
* [\#87](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/87) feat: introduce enable public delegation upgrade

## v0.0.3

This release includes features and bug fixes, mainly:
1. Gashub module;
2. Customized upgrade module;
3. Customized Tendermint with vote pool;
4. Disable create validator operation;
5. EIP712 bug fix;

* [\#72](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/72) feat: add gashub module
* [\#79](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/79) fix: disable create validator operation
* [\#80](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/80) fix: EIP712 issue with uint8[] type
* [\#81](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/81) feat: update tendermint to enable validator updates and votepool
* [\#82](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/82) feat: custom upgrade module
* [\#83](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/83) ci: fix release flow

## v0.0.2
This is the first release of gnfd-cosmos-SDK.

It includes two key features:

1. EIP721 signing schema support
2. New staking mechanism

* [\#38](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/38) ci: fix failed ci jobs
* [\#47](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/47) docs: change pull request template
* [\#36](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/36) feat: add support for eth address standard
* [\#68](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/68) fix: errors with EIP712 signature
* [\#71](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/71) test: fix the unstable UT bugs and remove useless testcase
* [\#70](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/70) fix: errors with EIP712 signature
* [\#46](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/46) feat: customize staking module for gnfd
* [\#73](https://github.com/bnb-chain/gnfd-cosmos-sdk/pull/73) feat: add bls key types into keyring management tool

## v0.0.1
Fork from cosmos-sdk 0.46.4
