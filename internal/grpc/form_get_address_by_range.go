package grpc

import (
	"context"
	"fmt"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

type derivationAddressByRangeForm struct {
	MnemonicWalletUUID    string    `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID `valid:"-"`

	AccountsParameters *anypb.Any `valid:"required"`
}

func (f *derivationAddressByRangeForm) LoadAndValidate(ctx context.Context,
	req *pbApi.GetMultipleAccountRequest,
) (valid bool, err error) {
	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.MnemonicWalletUUID = req.WalletIdentifier.WalletUUID

	if req.Parameters == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Multiple accounts parameter")
	}
	f.AccountsParameters = req.Parameters

	mnemonicWalletUUIDRaw, err := uuid.Parse(f.MnemonicWalletUUID)
	if err != nil {
		return false, err
	}
	f.MnemonicWalletUUIDRaw = mnemonicWalletUUIDRaw

	return true, nil
}
