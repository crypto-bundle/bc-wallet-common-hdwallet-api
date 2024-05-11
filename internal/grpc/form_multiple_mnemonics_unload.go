package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
)

type UnLoadMultipleMnemonicForm struct {
	MnemonicsList        []*UnLoadMnemonicForm `valid:"required"`
	MnemonicWalletsUUIDs []uuid.UUID           `valid:"-"`

	WalletsCount uint `valid:"-"`
}

func (f *UnLoadMultipleMnemonicForm) LoadAndValidate(ctx context.Context,
	req *pbApi.UnLoadMultipleMnemonicsRequest,
) (valid bool, err error) {
	if req.WalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identities")
	}

	f.MnemonicsList = make([]*UnLoadMnemonicForm, len(req.WalletIdentifier))
	f.MnemonicWalletsUUIDs = make([]uuid.UUID, len(req.WalletIdentifier))

	for i, _ := range req.WalletIdentifier {
		walletIdentityForm := &UnLoadMnemonicForm{
			WalletUUID: req.WalletIdentifier[i].WalletUUID,
		}

		walletIdentityForm.WalletUUIDRaw, _ = uuid.Parse(walletIdentityForm.WalletUUID)
		f.MnemonicsList[i] = walletIdentityForm
		f.MnemonicWalletsUUIDs[i] = walletIdentityForm.WalletUUIDRaw
	}

	f.WalletsCount = uint(len(req.WalletIdentifier))

	return true, nil
}
