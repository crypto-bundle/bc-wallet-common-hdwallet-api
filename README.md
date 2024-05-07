# Bc-wallet-common-hdwallet-api

## Description

Application for manage wallet sessions and execute wallet request, like: 
 - Get account address
 - Generate new mnemonic wallet
 - Validate mnemonic seed phrase
 - Sign data by private key by internal mnemonic account

This application is an environment for checking requests and transferring them for processing to the plugin, 
which makes cryptographic work for a specific blockchain.

HdWallet-api application is a second part of hdwallet applications bundle.
Another third part of bundle - target blockchain plugin.

## Api

### gRPC-API
* API documentation [bc-wallet-common-hdwallet-controller/docs/api/hdwallet_proto.md](https://github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/blob/develop/docs/api/hdwallet_proto.md)
* Golang proto generated code [bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet](https://github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/tree/develop/pkg/grpc/hdwallet)
* Protobuf descriptions [bc-wallet-common-hdwallet-controller/pkg/proto/hdwallet_api](https://github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/tree/develop/pkg/proto/hdwallet_api)

## Infrastructure dependencies

* **Hashicorp Vault** as service provider of secrets management and as provider of encrypt/decrypt sensitive information flow
* Plugin with implementation Hierarchical Deterministic Wallet for target blockchain.

### Hashicorp Vault

Application required two encryption keys:
* Common for whole crypto-bundle project transit key - crypto-bundle-bc-wallet-common-transit-key
* Target encryption key for hdwallet-controller and hdwallet-api - crypto-bundle-bc-wallet-tron-hdwallet

### HdWallet plugin

Implementation of HdWallet plugin must contains next exported functions:
* ```NewPoolUnitfunc(walletUUID string, mnemonicDecryptedData string) (interface{}, error)```
* ```GenerateMnemonic func() (string, error)```
* ```ValidateMnemonic func(mnemonic string) bool```
* ```GetPluginName func() string```
* ```GetPluginReleaseTag func() string```
* ```GetPluginCommitID func() string```
* ```GetPluginShortCommitID func() string```
* ```GetPluginBuildNumber func() string```
* ```GetPluginBuildDateTS func() string```

Function NewPoolUnit must return implementation for interface
```Go
import (
    pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

    "google.golang.org/protobuf/types/known/anypb"
)

type WalletPoolUnitService interface {
	UnloadWallet() error

	GetWalletUUID() string
	LoadAccount(ctx context.Context,
		accountParameters *anypb.Any,
	) (accountAddress *string, err  error)
	GetAccountAddress(ctx context.Context,
		accountParameters *anypb.Any,
	) (accountAddress *string, err error)
	GetMultipleAccounts(ctx context.Context,
		multipleAccountsParameters *anypb.Any,
	) (accountCount uint, accountsList []*pbCommon.AccountIdentity, err error)
	SignData(ctx context.Context,
		accountParameters *anypb.Any,
		dataForSign []byte,
	) (accountAddress *string, signedData []byte, err error)
}
```

Examples of plugin implementations:
* [bc-wallet-tron-hdwallet](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet)
* [bc-wallet-ethereum-hdwallet](https://github.com/crypto-bundle/bc-wallet-ethereum-hdwallet)
* [bc-wallet-bitcoin-hdwallet](https://github.com/crypto-bundle/bc-wallet-bitcoin-hdwallet)

### Bc-wallet-common-hdwallet-controller

Repository of hdwallet-controller - [bc-wallet-common-hdwallet-controller](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet-controller)

Application must work in pair with instance bc-wallet-common-hdwallet-controller.
For example in case of Tron blockchain:
Instance of bc-wallet-common-hdwallet-controller - **bc-wallet-tron-hdwallet-controller** must work with instance
of bc-wallet-common-hdwallet-api - **bc-wallet-tron-hdwallet-api**.

Communication between hdwallet-controller and hdwallet-api works via gRPC unix-file socket connection.
You can control socket file path by change `HDWALLET_UNIX_SOCKET_PATH` environment variable.
Also, you can set target blockchain of hdwallet-controller and hdwallet-api via `PROCESSING_NETWORK` env variable.

## Environment variables

Full example of env variables you can see in  [env-api-example.env](./env-controller-example.env) file.

## Deployment

Currently, support only kubernetes deployment flow via Helm

### Kubernetes
Application must be deployed as part of bc-wallet-<BLOCKCHAIN_NAME>-hdwallet bundle.
Application must be started as single container in Kubernetes Pod with shared volume.

You can see example of HELM-chart deployment application in next repositories:
* [bc-wallet-tron-hdwallet-api/deploy/helm/hdwallet](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet/tree/develop/deploy/helm/hdwallet)
* [bc-wallet-ethereum-hdwallet-api/deploy/helm/hdwallet](https://github.com/crypto-bundle/bc-wallet-ethereum-hdwallet/tree/develop/deploy/helm/hdwallet)

## Licence

**bc-wallet-common-hdwallet-api** is licensed under the [MIT](./LICENSE) License.