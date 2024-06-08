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
	"crypto/sha256"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameEncryptMnemonic = "EncryptMnemonic"
)

type encryptMnemonicHandler struct {
	l *zap.Logger

	transitEncryptorSvc encryptService
	appEncryptorSvc     encryptService

	mnemonicValidatorFunc validateMnemonicFunc
}

// nolint:funlen // fixme
func (h *encryptMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.EncryptMnemonicRequest,
) (*pbApi.EncryptMnemonicResponse, error) {
	var err error
	tCtx, _, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	vf := &EncryptMnemonicForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	decryptedData, err := h.transitEncryptorSvc.Decrypt(vf.TransitEncryptedMnemonicData)
	if err != nil {
		return nil, err
	}
	defer func() {
		for i := range decryptedData {
			decryptedData[i] = 0
		}
	}()

	isValidMnemoPhrase := h.mnemonicValidatorFunc(string(decryptedData))
	if !isValidMnemoPhrase {
		return nil, status.Error(codes.InvalidArgument, "mnemonic phrase is not valid")
	}

	encryptedMnemonicData, err := h.appEncryptorSvc.Encrypt(decryptedData)
	if err != nil {
		return nil, err
	}

	mnemonicHash := fmt.Sprintf("%x", sha256.Sum256(decryptedData))
	req.WalletIdentifier.WalletHash = mnemonicHash

	return &pbApi.EncryptMnemonicResponse{
		WalletIdentifier:      req.WalletIdentifier,
		EncryptedMnemonicData: encryptedMnemonicData,
	}, nil
}

func MakeEncryptMnemonicHandler(loggerEntry *zap.Logger,
	transitEncryptorSvc encryptService,
	appEncryptorSvc encryptService,
) *encryptMnemonicHandler {
	return &encryptMnemonicHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameEncryptMnemonic)),

		transitEncryptorSvc: transitEncryptorSvc,
		appEncryptorSvc:     appEncryptorSvc,
	}
}
