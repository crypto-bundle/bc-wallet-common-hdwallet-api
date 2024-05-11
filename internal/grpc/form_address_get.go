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
