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
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletApiServer

	logger *zap.Logger
	cfg    configService
	// all GRPC handlers
	generateMnemonicHandlerSvc        generateMnemonicHandlerService
	validateMnemonicHandlerSvc        validateMnemonicHandlerService
	loadMnemonicHandlerSvc            loadMnemonicHandlerService
	unLoadMnemonicHandlerSvc          unLoadMnemonicHandlerService
	unLoadMultipleMnemonicsHandlerSvc unLoadMultipleMnemonicsHandlerService
	encryptMnemonicHandlerSvc         encryptMnemonicHandlerService
	getAccountHandlerSvc              getAccountHandlerService
	getAccountsSvc                    getDerivationsAddressesHandlerService
	loadDerivationAddressSvc          loadDerivationsAddressesHandlerService
	signDataSvc                       signDataHandlerService
}

func (h *grpcServerHandle) GenerateMnemonic(ctx context.Context,
	req *pbApi.GenerateMnemonicRequest,
) (*pbApi.GenerateMnemonicResponse, error) {
	return h.generateMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) ValidateMnemonic(ctx context.Context,
	req *pbApi.ValidateMnemonicRequest,
) (*pbApi.ValidateMnemonicResponse, error) {
	return h.validateMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) LoadMnemonic(ctx context.Context,
	req *pbApi.LoadMnemonicRequest,
) (*pbApi.LoadMnemonicResponse, error) {
	return h.loadMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) UnLoadMnemonic(ctx context.Context,
	req *pbApi.UnLoadMnemonicRequest,
) (*pbApi.UnLoadMnemonicResponse, error) {
	return h.unLoadMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) UnLoadMultipleMnemonics(ctx context.Context,
	req *pbApi.UnLoadMultipleMnemonicsRequest,
) (*pbApi.UnLoadMultipleMnemonicsResponse, error) {
	return h.unLoadMultipleMnemonicsHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) EncryptMnemonic(context.Context,
	*pbApi.EncryptMnemonicRequest,
) (*pbApi.EncryptMnemonicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EncryptMnemonic not implemented")
}

func (h *grpcServerHandle) GetAccount(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (*pbApi.GetAccountResponse, error) {
	return h.getAccountHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetMultipleAccounts(ctx context.Context,
	req *pbApi.GetMultipleAccountRequest,
) (*pbApi.GetMultipleAccountResponse, error) {
	return h.getAccountsSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) LoadAccount(ctx context.Context,
	req *pbApi.LoadAccountRequest,
) (*pbApi.LoadAccountsResponse, error) {
	return h.loadDerivationAddressSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) SignData(ctx context.Context,
	req *pbApi.SignDataRequest,
) (*pbApi.SignDataResponse, error) {
	return h.signDataSvc.Handle(ctx, req)
}

// NewHandlers - create instance of grpc-handler service
func NewHandlers(loggerSrv *zap.Logger,
	mnemoGenFunc generateMnemonicFunc,
	mnemoValidatorFunc validateMnemonicFunc,
	transitEncryptorSvc encryptService,
	appEncryptorSvc encryptService,
	walletPoolSvc walletPoolService,
) pbApi.HdWalletApiServer {

	l := loggerSrv.Named("grpc.server.handler")

	return &grpcServerHandle{
		UnimplementedHdWalletApiServer: &pbApi.UnimplementedHdWalletApiServer{},
		logger:                         l,

		generateMnemonicHandlerSvc:        MakeGenerateMnemonicHandler(l, mnemoGenFunc, appEncryptorSvc),
		validateMnemonicHandlerSvc:        MakeValidateMnemonicHandler(l, mnemoValidatorFunc, appEncryptorSvc),
		loadMnemonicHandlerSvc:            MakeLoadMnemonicHandler(l, walletPoolSvc),
		unLoadMnemonicHandlerSvc:          MakeUnLoadMnemonicHandler(l, walletPoolSvc),
		unLoadMultipleMnemonicsHandlerSvc: MakeUnLoadMultipleMnemonicsHandler(l, walletPoolSvc),
		encryptMnemonicHandlerSvc:         MakeEncryptMnemonicHandler(l, transitEncryptorSvc, appEncryptorSvc),
		loadDerivationAddressSvc:          MakeLoadDerivationAddressHandlerHandler(l, walletPoolSvc),
		getAccountHandlerSvc:              MakeGetDerivationAddressHandler(l, walletPoolSvc),
		getAccountsSvc:                    MakeGetDerivationsAddressesHandler(l, walletPoolSvc),
		signDataSvc:                       MakeSignDataHandler(l, walletPoolSvc),
	}
}
