package grpc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/anypb"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
)

type SignDataForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	AccountParameters *anypb.Any `valid:"required"`

	DataForSign []byte `valid:"required"`
}

func (f *SignDataForm) LoadAndValidate(ctx context.Context,
	req *pbApi.SignDataRequest,
) (valid bool, err error) {
	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.WalletIdentifier.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(req.WalletIdentifier.WalletUUID)
	if err != nil {
		return false, err
	}

	if req.AccountIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}
	if req.AccountIdentifier.Parameters == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity parameters")
	}

	f.AccountParameters = req.AccountIdentifier.Parameters
	f.DataForSign = req.DataForSign

	return true, nil
}
