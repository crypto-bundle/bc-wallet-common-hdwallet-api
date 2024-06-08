/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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
	GetUnixFileNameTemplate() string
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
