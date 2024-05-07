package grpc

import (
	"context"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
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
	validateMnemonicHandlerSvc        validateMnemonicHandlerService
	loadMnemonicHandlerSvc            loadMnemonicHandlerService
	unLoadMnemonicHandlerSvc          unLoadMnemonicHandlerService
	unLoadMultipleMnemonicsHandlerSvc unLoadMultipleMnemonicsHandlerService
	encryptMnemonicHandlerSvc         encryptMnemonicHandlerService
	getAccountHandlerSvc              getAccountHandlerService
	getAccountsSvc                    getDerivationsAddressesHandlerService
	loadDerivationAddressSvc          loadDerivationsAddressesHandlerService
	signDataSvc                       signDataHandlerService
}

func (h *grpcServerHandle) GenerateMnemonic(ctx context.Context,
	req *pbApi.GenerateMnemonicRequest,
) (*pbApi.GenerateMnemonicResponse, error) {
	return h.generateMnemonicHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) ValidateMnemonic(ctx context.Context,
	req *pbApi.ValidateMnemonicRequest,
) (*pbApi.ValidateMnemonicResponse, error) {
	return h.validateMnemonicHandlerSvc.Handle(ctx, req)
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

func (h *grpcServerHandle) GetAccount(ctx context.Context,
	req *pbApi.GetAccountRequest,
) (*pbApi.GetAccountResponse, error) {
	return h.getAccountHandlerSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) GetMultipleAccounts(ctx context.Context,
	req *pbApi.GetMultipleAccountRequest,
) (*pbApi.GetMultipleAccountResponse, error) {
	return h.getAccountsSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) LoadAccount(ctx context.Context,
	req *pbApi.LoadAccountRequest,
) (*pbApi.LoadAccountsResponse, error) {
	return h.loadDerivationAddressSvc.Handle(ctx, req)
}

func (h *grpcServerHandle) SignData(ctx context.Context,
	req *pbApi.SignDataRequest,
) (*pbApi.SignDataResponse, error) {
	return h.signDataSvc.Handle(ctx, req)
}

// NewHandlers - create instance of grpc-handler service
func NewHandlers(loggerSrv *zap.Logger,
	mnemoGenFunc generateMnemonicFunc,
	mnemoValidatorFunc validateMnemonicFunc,
	transitEncryptorSvc encryptService,
	appEncryptorSvc encryptService,
	walletPoolSvc walletPoolService,
) pbApi.HdWalletApiServer {

	l := loggerSrv.Named("grpc.server.handler")

	return &grpcServerHandle{
		UnimplementedHdWalletApiServer: &pbApi.UnimplementedHdWalletApiServer{},
		logger:                         l,

		generateMnemonicHandlerSvc:        MakeGenerateMnemonicHandler(l, mnemoGenFunc, appEncryptorSvc),
		validateMnemonicHandlerSvc:        MakeValidateMnemonicHandler(l, mnemoValidatorFunc, appEncryptorSvc),
		loadMnemonicHandlerSvc:            MakeLoadMnemonicHandler(l, walletPoolSvc),
		unLoadMnemonicHandlerSvc:          MakeUnLoadMnemonicHandler(l, walletPoolSvc),
		unLoadMultipleMnemonicsHandlerSvc: MakeUnLoadMultipleMnemonicsHandler(l, walletPoolSvc),
		encryptMnemonicHandlerSvc:         MakeEncryptMnemonicHandler(l, transitEncryptorSvc, appEncryptorSvc),
		loadDerivationAddressSvc:          MakeLoadDerivationAddressHandlerHandler(l, walletPoolSvc),
		getAccountHandlerSvc:              MakeGetDerivationAddressHandler(l, walletPoolSvc),
		getAccountsSvc:                    MakeGetDerivationsAddressesHandler(l, walletPoolSvc),
		signDataSvc:                       MakeSignDataHandler(l, walletPoolSvc),
	}
}
