package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
)

type LoadMnemonicForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	TimeToLive            uint64 `valid:"type(uint64),numeric,required"`
	EncryptedMnemonicData []byte `valid:"required"`
}

func (f *LoadMnemonicForm) LoadAndValidate(ctx context.Context,
	req *pbApi.LoadMnemonicRequest,
) (valid bool, err error) {
	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.WalletIdentifier.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(req.WalletIdentifier.WalletUUID)
	if err != nil {
		return false, err
	}

	f.TimeToLive = req.TimeToLive
	f.EncryptedMnemonicData = req.EncryptedMnemonicData

	return true, nil
}
