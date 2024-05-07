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
	MethodNameSignMnemonic = "Sign"
)

type signDataHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *signDataHandler) Handle(ctx context.Context,
	req *pbApi.SignDataRequest,
) (*pbApi.SignDataResponse, error) {
	var err error
	tCtx, _, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	vf := &SignDataForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	addr, signedData, err := h.walletPoolSvc.SignData(tCtx, vf.WalletUUIDRaw,
		vf.AccountParameters,
		vf.DataForSign)
	if err != nil {
		h.l.Error("unable to sign data", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if addr == nil || signedData == nil {
		return nil, status.Error(codes.ResourceExhausted, "wallet not loaded")
	}

	req.AccountIdentifier.Address = *addr

	return &pbApi.SignDataResponse{
		WalletIdentifier:  req.WalletIdentifier,
		AccountIdentifier: req.AccountIdentifier,
		SignedData:        signedData,
	}, nil
}

func MakeSignDataHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *signDataHandler {
	return &signDataHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameSignMnemonic)),

		walletPoolSvc: walletPoolSvc,
	}
}
