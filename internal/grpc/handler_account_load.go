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
	MethodNameLoadAccount = "LoadAccount"
)

type loadAccountHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *loadAccountHandler) Handle(ctx context.Context,
	req *pbApi.LoadAccountRequest,
) (*pbApi.LoadAccountsResponse, error) {
	var err error
	tCtx, _, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	vf := &AccountForm{}
	valid, err := vf.LoadAndValidateLoadAddrReq(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	addr, err := h.walletPoolSvc.LoadAccount(tCtx, vf.WalletUUIDRaw,
		vf.AccountParameters)
	if err != nil {
		h.l.Error("unable to get address by path", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if addr == nil {
		return nil, status.Error(codes.ResourceExhausted, "wallet not loaded or session already expired")
	}

	req.AccountIdentifier.Address = *addr

	return &pbApi.LoadAccountsResponse{
		WalletIdentifier:  req.WalletIdentifier,
		AccountIdentifier: req.AccountIdentifier,
	}, nil
}

func MakeLoadDerivationAddressHandlerHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *loadAccountHandler {
	return &loadAccountHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameLoadAccount)),
		walletPoolSvc: walletPoolSvc,
	}
}
