# Change Log

## [v0.0.31] 09.06.2024
### Added
* Added support of new plugin versions - added chainID support:
  * Lookup of GetChainID function
  * Lookup of GetSupportedChainIDs function
* Added support config environment variable - HDWALLET_CHAIN_ID

## [v0.0.30] 08.06.2024
### Changed
* Changed MIT License to NON-AI MIT
* Changed version of hdwallet-controller dependency to v0.0.29

## [v0.0.29] 08.06.2024
### Changed
* Changed flow of reset unload interval - added compare of unload time value before timer reset

## [v0.0.28] 12.05.2024
### Added
* Added usage of lib-profiler:
  * Added lib-profiler config ENV to example.env file
  * Added lib-profiler config to application config
### Changed
* Added support of new version bc-wallet-common-hdwallet-controller v0.0.27
* Changed unix-socket listen flow:
  * Implemented usage of new ENV
  * Changed listen flow
  * Fixed bugs in shutdown gRPC server process
* Changed build of common-hdwallet-api docker container:
  * golang:1.22.2-bookworm is stage 1 image
  * busybox:glibc is stage 2 image
* Small fixes to prevent "Race Condition" state
* Changed validation flow - removed usage of govalidator

## [v0.0.27] 05.05.2024
### Changed
* HdWallet api application move to another repository - https://github.com/crypto-bundle/bc-wallet-common-hdwallet-api
* Changed go-name space
* Changed docker container name
* Removed Tron plugin extension, now repository [bc-wallet-tron-hdwallet](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet) is repository with implementation Tron blockhain cryptographic processes
* New content of README.md file
* Bump version of bc-wallet-common-lib-vault v0.0.14
* Changed README.md content

## [v0.0.26] 05.05.2024
### Added
* Added plugin path ENV variable
* Added ldflags support in build application flow
* Added unit-tests to all mnemonicWalletUnit methods
* Added plugin wrapper package
* Added ldflags support in build plugin flow
* Added support of new version bc-wallet-common-hdwallet-controller v0.0.25
### Changed
* Refactored hd-wallet service-component
  * All struct and variables is un-exportable
  * Moved to tron plugin directory
* Changed gRPC methods:
  * GetDerivationAddress replaced by GetAccount
  * GetDerivationAddressByRange replaced by GetMultipleAccounts
* Refactored wallet pool unit service-component
  * Changed code for support plugin flow
  * Moved to tron plugin directory
  * Now plugin must support next exported functions:
    * NewPoolUnit
    * GenerateMnemonic
    * ValidateMnemonic
    * GetPluginName
    * GetPluginReleaseTag
    * GetPluginCommitID
    * GetPluginShortCommitID
    * GetPluginBuildNumber
    * GetPluginBuildDateTS

## [v0.0.25] 24.04.2024
### Added
* Implemented new gRPC methods:
  * GenerateMnemonic
  * LoadMnemonic
  * UnLoadMnemonic
  * EncryptMnemonic
  * ValidateMnemonic
  * UnLoadMultipleMnemonics
  * LoadDerivationAddress
  * SignData
### Changed
* Removed all types of store - postgres, redis, nats. Hdwallet-api app now storage less application
* Removed multiple mnemonics per wallet flow - now one wallet - one mnemonic
* Changed mnemonic management flow:
  * Now all wallet only in memory
  * Logic of load/unload wallet in gRPC method requests
  * Added vault encrypt/decrypt flow for passing mnemonic data between hdwallet-controller and hdwallet-api
* Bump go version 1.19 -> 1.22
* Integrated new version of hdwallet-controller dependency - v0.0.24
* Bump common-lib version:
  * bc-wallet-common-lib-config v0.0.5
  * bc-wallet-common-lib-grpc v0.0.4
  * bc-wallet-common-lib-healthcheck v0.0.4
  * bc-wallet-common-lib-logger v0.0.4
  * bc-wallet-common-lib-nats-queue v0.1.12
  * bc-wallet-common-lib-postgres v0.0.8
  * bc-wallet-common-lib-redis v0.0.7
  * bc-wallet-common-lib-tracer v0.0.4
  * bc-wallet-common-lib-vault v0.0.13

## [v0.0.24] 14.02.2024
### Info
Start of big application refactoring
### Added
* Added wallet sessions entities for storing in persistent and cache stores
### Changed
* Separated application on two parts
  * bc-wallet-common-hdwallet-controller
  * bc-wallet-tron-hdwallet
* Changed GetDerivationAddressByRange gRPC method - now support get addresses by multiple ranges
* Added HdWallet API proto description
  * new gRPC method - GenerateMnemonic
  * new gRPC method - LoadMnemonic
  * new gRPC method - UnLoadMnemonic
* Added Controller API proto description
  * new gRPC method - StartWalletSession
  * new gRPC method - GetWalletSession
* Removed go-tron-sdk dependency

## [v0.0.23 - v0.0.6] 28.04.2023 - 05.04.2023
### Added
* Added gRPC client wrapper
* Small security improvements:
  * Filling private keys with zeroes - private key clearing
* Added data cache flow for storing wallet in redis and nats
* Added new gRPC-handler - GetWalletInfo
### Changed
* Changed deployment flow
  * Added helm-chart option for docker container repository
  * Fixed helm-chart template for VAULT_DATA_PATH variable
* Optimization in get addresses by range flow
* Clone private key in sign transaction flow
### Fixed
* Fixed bug in sign transaction flow
* Fixed migrations - wrong rollback SQL-code, missing drop index and drop table

## [v0.0.5] 05.04.2023
### Added
* Encryption private data via hashicorp vault
* Added gRPC client config
### Changed
* Cleaned up repository:
  * Removed ansible database deployment script
  * Removed vault polices
  * Removed private data from helm-chart
* Updated common-libs:
  * removed old bc-wallet-common dependency
  * integrated lib-common dependencies:
    * lib-postgres
    * lib-config
    * lib-grpc
    * lib-tracer
    * lib-logger
    * lib-vault
### Fixed
* Fixed bug in wallet init stage
* Fixed crash in wallet pool init stage
* Fixed bugs in flow in new wallet creation

## [v0.0.4] 16.03.2023
### Changed
* Refactoring wallet pool service-component:
  * Added wallet pool unit
  * Added unit maker
  * Added support of multiple and single mnemonic wallet
  * Added timer for mnemonic unloading flow

## [initial] - 16.03.2023 - 06.03.2022
* Created go module as bc-wallet-eth-hdwallet
* Added proto files for gRPC API
* Integrated common dependencies
* Moved to crypto-bundle namespace
* Added wallet persistent store
* Added functionality for gRPC handlers
  * AddNewWallet
  * GetDerivationAddress
  * GetEnabledWallets
  * GetDerivationAddressByRange
* Added mnemonic encryption via rsa-keys
* Added MIT licence
* Refactoring service for supporting TRON blockchain
* Created Helm chart


