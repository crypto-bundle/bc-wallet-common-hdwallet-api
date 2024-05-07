package grpc

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/app"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameUnLoadMnemonic = "UnLoadMnemonic"
)

type unLoadMnemonicHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *unLoadMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.UnLoadMnemonicRequest,
) (*pbApi.UnLoadMnemonicResponse, error) {
	var err error
	tCtx, _, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	vf := &UnLoadMnemonicForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletUUID, err := h.walletPoolSvc.UnloadWalletUnit(tCtx, vf.WalletUUIDRaw)
	if err != nil {
		h.l.Error("unable to unload mnemonic wallet", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if walletUUID == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	if walletUUID.String() != req.WalletIdentifier.WalletUUID {
		h.l.Error("requested wallet uuid not equal with wallet pool unit uuid",
			zap.Error(ErrWalletUUIDMismatched),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &pbApi.UnLoadMnemonicResponse{
		WalletIdentifier: req.WalletIdentifier,
	}, nil
}

func MakeUnLoadMnemonicHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *unLoadMnemonicHandler {
	return &unLoadMnemonicHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameUnLoadMnemonic)),
		walletPoolSvc: walletPoolSvc,
	}
}
