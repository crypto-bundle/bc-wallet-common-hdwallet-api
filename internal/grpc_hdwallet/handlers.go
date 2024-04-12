package grpc_hdwallet

import (
	"context"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpcServerHandle is wrapper struct for implementation all grpc handlers
type grpcServerHandle struct {
	*pbApi.UnimplementedHdWalletApiServer

	logger *zap.Logger
	cfg    configService
	// all GRPC handlers
	generateMnemonicHandlerSvc        generateMnemonicHandlerService
	loadMnemonicHandlerSvc            loadMnemonicHandlerService
	unLoadMnemonicHandlerSvc          unLoadMnemonicHandlerService
	unLoadMultipleMnemonicsHandlerSvc unLoadMultipleMnemonicsHandlerService
	encryptMnemonicHandlerSvc         encryptMnemonicHandlerService
	getDerivationAddressSvc           getDerivationAddressHandlerService
	getDerivationsAddressesSvc        getDerivationsAddressesHandlerService
	loadDerivationAddressSvc          loadDerivationsAddressesHandlerService
	signDataSvc                       signDataHandlerService
}

func (h *grpcServerHandle) GenerateMnemonic(ctx context.Context,
	req *pbApi.GenerateMnemonicRequest,
) (*pbApi.GenerateMnemonicResponse, error) {
	return h.generateMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) LoadMnemonic(ctx context.Context,
	req *pbApi.LoadMnemonicRequest,
) (*pbApi.LoadMnemonicResponse, error) {
	return h.loadMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) UnLoadMnemonic(ctx context.Context,
	req *pbApi.UnLoadMnemonicRequest,
) (*pbApi.UnLoadMnemonicResponse, error) {
	return h.unLoadMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) UnLoadMultipleMnemonics(ctx context.Context,
	req *pbApi.UnLoadMultipleMnemonicsRequest,
) (*pbApi.UnLoadMultipleMnemonicsResponse, error) {
	return h.unLoadMultipleMnemonicsHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) EncryptMnemonic(context.Context,
	*pbApi.EncryptMnemonicRequest,
) (*pbApi.EncryptMnemonicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EncryptMnemonic not implemented")
}

func (h *grpcServerHandle) GetDerivationAddress(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	return h.getDerivationAddressSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetDerivationAddressByRange(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	return h.getDerivationsAddressesSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) LoadDerivationAddress(ctx context.Context,
	req *pbApi.LoadDerivationAddressRequest,
) (*pbApi.LoadDerivationAddressResponse, error) {
	return h.loadDerivationAddressSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) SignData(ctx context.Context,
	req *pbApi.SignDataRequest,
) (*pbApi.SignDataResponse, error) {
	return h.signDataSvc.Handle(ctx, req)
}

// New instance of service
func New(loggerSrv *zap.Logger) pbApi.HdWalletApiServer {

	l := loggerSrv.Named("grpc.server.handler").With(
		zap.String(app.BlockChainNameTag, app.BlockChainName))

	//addrRespPool := &sync.Pool{New: func() any {
	//	return new(pbApi.DerivationAddressIdentity)
	//}}

	return &grpcServerHandle{
		UnimplementedHdWalletApiServer: &pbApi.UnimplementedHdWalletApiServer{},
		logger:                         l,

		generateMnemonicHandlerSvc:        MakeGenerateMnemonicHandler(l),
		loadMnemonicHandlerSvc:            MakeLoadMnemonicHandler(l),
		unLoadMnemonicHandlerSvc:          MakeUnLoadMnemonicHandler(l),
		unLoadMultipleMnemonicsHandlerSvc: MakeUnLoadMultipleMnemonicsHandler(l),
		encryptMnemonicHandlerSvc:         MakeEncryptMnemonicHandler(l),
		getDerivationAddressSvc:           MakeGetDerivationAddressHandler(l),
		getDerivationsAddressesSvc:        MakeGetDerivationsAddressesHandler(l),
		signDataSvc:                       MakeSignDataHandler(l),
	}
}
