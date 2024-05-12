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
	MethodNameGetDerivationsAddresses = "GetDerivationsAddresses"
)

type getDerivationsAddressesHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *getDerivationsAddressesHandler) Handle(ctx context.Context,
	req *pbApi.GetMultipleAccountRequest,
) (*pbApi.GetMultipleAccountResponse, error) {
	var err error
	tCtx, _, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	vf := &derivationAddressByRangeForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	count, list, err := h.walletPoolSvc.GetMultipleAccounts(tCtx, vf.MnemonicWalletUUIDRaw,
		vf.AccountsParameters)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if count == 0 {
		return nil, status.Error(codes.ResourceExhausted, "wallet not loaded or session already expired")
	}

	return &pbApi.GetMultipleAccountResponse{
		WalletIdentifier:       req.WalletIdentifier,
		AccountIdentitiesCount: uint64(count),
		AccountIdentifier:      list,
	}, nil
}

func MakeGetDerivationsAddressesHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *getDerivationsAddressesHandler {
	return &getDerivationsAddressesHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameGetDerivationsAddresses)),

		walletPoolSvc: walletPoolSvc,
	}
}
