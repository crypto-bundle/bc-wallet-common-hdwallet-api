package handlers

import (
	"bc-wallet-eth-hdwallet/internal/entities"
	"context"
)

const (
	MethodNameTag = "method_name"
)

type walleter interface {
	GetAddressByPath(ctx context.Context,
		walletUUID string,
		account, change, index uint32,
	) (string, error)

	GetEnabledWalletsUUID(ctx context.Context) ([]string, error)

	CreateNewMnemonicWallet(ctx context.Context,
		title string,
		purpose string,
		isHot bool,
	) (*entities.MnemonicWallet, error)
}