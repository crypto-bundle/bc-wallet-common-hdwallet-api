package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
)

type EncryptMnemonicForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	TransitEncryptedMnemonicData []byte `valid:"type([]byte]),required"`
}

func (f *EncryptMnemonicForm) LoadAndValidate(ctx context.Context,
	req *pbApi.EncryptMnemonicRequest,
) (valid bool, err error) {
	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.WalletIdentifier.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(req.WalletIdentifier.WalletUUID)
	if err != nil {
		return false, err
	}

	f.TransitEncryptedMnemonicData = req.MnemonicData

	return true, nil
}
