package grpc

import (
	"context"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameUnLoadMultipleMnemonics = "UnLoadMultipleMnemonics"
)

type unLoadMultipleMnemonicsHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *unLoadMultipleMnemonicsHandler) Handle(ctx context.Context,
	req *pbApi.UnLoadMultipleMnemonicsRequest,
) (*pbApi.UnLoadMultipleMnemonicsResponse, error) {
	var err error
	tCtx, _, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	vf := &UnLoadMultipleMnemonicForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	err = h.walletPoolSvc.UnloadMultipleWalletUnit(tCtx, vf.MnemonicWalletsUUIDs)
	if err != nil {
		h.l.Error("unable unload multiple mnemonic wallets", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return nil, nil
}

func MakeUnLoadMultipleMnemonicsHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *unLoadMultipleMnemonicsHandler {
	return &unLoadMultipleMnemonicsHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameUnLoadMultipleMnemonics)),

		walletPoolSvc: walletPoolSvc,
	}
}
