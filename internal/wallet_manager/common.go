package wallet_manager

import (
	"context"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

	"google.golang.org/protobuf/types/known/anypb"
)

type configService interface {
	GetHdWalletPluginPath() string
}

type walletMakerFunc func(walletUUID string,
	mnemonicDecryptedData string,
) (interface{}, error)

type WalletPoolUnitService interface {
	UnloadWallet() error

	GetWalletUUID() string
	LoadAccount(ctx context.Context,
		accountParameters *anypb.Any,
	) (*string, error)
	GetAccountAddress(ctx context.Context,
		accountParameters *anypb.Any,
	) (*string, error)
	GetMultipleAccounts(ctx context.Context,
		multipleAccountsParameters *anypb.Any,
	) (uint, []*pbCommon.AccountIdentity, error)
	SignData(ctx context.Context,
		accountParameters *anypb.Any,
		dataForSign []byte,
	) (*string, []byte, error)
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}
