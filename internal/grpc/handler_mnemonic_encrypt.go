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
