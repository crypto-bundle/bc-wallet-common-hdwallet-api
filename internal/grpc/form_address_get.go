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
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"google.golang.org/protobuf/types/known/anypb"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
)

type AccountForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	AccountParameters *anypb.Any `valid:"required"`
}

func (f *AccountForm) LoadAndValidateLoadAddrReq(ctx context.Context,
	req *pbApi.LoadAccountRequest,
) (valid bool, err error) {
	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}

	if req.AccountIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}
	if req.AccountIdentifier.Parameters == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity parameters")
	}

	return f.validate(req.WalletIdentifier, req.AccountIdentifier)
}

func (f *AccountForm) LoadAndValidateGetAddrReq(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (valid bool, err error) {
	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}

	if req.AccountIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}
	if req.AccountIdentifier.Parameters == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity parameters")
	}

	return f.validate(req.WalletIdentifier, req.AccountIdentifier)
}

func (f *AccountForm) validate(mnemoIdentifier *pbCommon.MnemonicWalletIdentity,
	accIdentifier *pbCommon.AccountIdentity,
) (valid bool, err error) {

	f.WalletUUID = mnemoIdentifier.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(mnemoIdentifier.WalletUUID)
	if err != nil {
		return false, err
	}

	f.AccountParameters = accIdentifier.Parameters

	return true, nil
}
