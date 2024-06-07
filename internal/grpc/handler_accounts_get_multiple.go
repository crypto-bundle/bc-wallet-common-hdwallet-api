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
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameGetDerivationsAddresses = "GetDerivationsAddresses"
)

type getDerivationsAddressesHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *getDerivationsAddressesHandler) Handle(ctx context.Context,
	req *pbApi.GetMultipleAccountRequest,
) (*pbApi.GetMultipleAccountResponse, error) {
	var err error
	tCtx, _, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	vf := &derivationAddressByRangeForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	count, list, err := h.walletPoolSvc.GetMultipleAccounts(tCtx, vf.MnemonicWalletUUIDRaw,
		vf.AccountsParameters)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if count == 0 {
		return nil, status.Error(codes.ResourceExhausted, "wallet not loaded or session already expired")
	}

	return &pbApi.GetMultipleAccountResponse{
		WalletIdentifier:       req.WalletIdentifier,
		AccountIdentitiesCount: uint64(count),
		AccountIdentifier:      list,
	}, nil
}

func MakeGetDerivationsAddressesHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *getDerivationsAddressesHandler {
	return &getDerivationsAddressesHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameGetDerivationsAddresses)),

		walletPoolSvc: walletPoolSvc,
	}
}
