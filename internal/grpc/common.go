package grpc

import (
	"context"
	"time"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	MethodNameTag = "method_name"
)

type generateMnemonicFunc func() (string, error)
type validateMnemonicFunc func(mnemonic string) bool

type configService interface {
	IsDev() bool
	IsDebug() bool
	IsLocal() bool

	GetConnectionPath() string
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}

type walletPoolService interface {
	AddAndStartWalletUnit(_ context.Context,
		mnemonicWalletUUID uuid.UUID,
		timeToLive time.Duration,
		mnemonicEncryptedData []byte,
	) error
	LoadAccount(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		accountParameters *anypb.Any,
	) (*string, error)
	UnloadWalletUnit(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
	) (*uuid.UUID, error)
	UnloadMultipleWalletUnit(ctx context.Context,
		mnemonicWalletUUIDs []uuid.UUID,
	) error
	GetAccountAddress(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		accountParameters *anypb.Any,
	) (*string, error)
	GetMultipleAccounts(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		multipleAccountsParameters *anypb.Any,
	) (uint, []*pbCommon.AccountIdentity, error)
	SignData(ctx context.Context,
		mnemonicUUID uuid.UUID,
		accountParameters *anypb.Any,
		transactionData []byte,
	) (*string, []byte, error)
}

type generateMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.GenerateMnemonicRequest) (*pbApi.GenerateMnemonicResponse, error)
}

type validateMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.ValidateMnemonicRequest) (*pbApi.ValidateMnemonicResponse, error)
}

type loadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.LoadMnemonicRequest) (*pbApi.LoadMnemonicResponse, error)
}

type unLoadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.UnLoadMnemonicRequest) (*pbApi.UnLoadMnemonicResponse, error)
}

type unLoadMultipleMnemonicsHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.UnLoadMultipleMnemonicsRequest,
	) (*pbApi.UnLoadMultipleMnemonicsResponse, error)
}

type encryptMnemonicHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.EncryptMnemonicRequest,
	) (*pbApi.EncryptMnemonicResponse, error)
}

type getDerivationsAddressesHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.GetMultipleAccountRequest,
	) (*pbApi.GetMultipleAccountResponse, error)
}
type loadDerivationsAddressesHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.LoadAccountRequest,
	) (*pbApi.LoadAccountsResponse, error)
}

type signDataHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.SignDataRequest,
	) (*pbApi.SignDataResponse, error)
}

type getAccountHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.GetAccountRequest,
	) (*pbApi.GetAccountResponse, error)
}
