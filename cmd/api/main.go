package main

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/plugin"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonProfiler "github.com/crypto-bundle/bc-wallet-common-lib-profiler/pkg/profiler"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/config"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/grpc"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/wallet_manager"

	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"

	"go.uber.org/zap"
)

// DO NOT EDIT THESE VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
var (
	// ReleaseTag - release tag in TAG.SHORT_COMMIT_ID.BUILD_NUMBER.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ReleaseTag = "v0.0.0-00000000-100500"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	CommitID = "0000000000000000000000000000000000000000"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ShortCommitID = "00000000"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildNumber string = "100500"

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildDateTS string = "1713280105"
)

func main() {
	var err error
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	appCfg, vaultSvc, err := config.Prepare(ctx, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerSrv, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}
	loggerEntry := loggerSrv.NewLoggerEntry("main").
		With(zap.String(app.BlockChainNameTag, appCfg.GetNetworkName()))

	transitSvc := commonVault.NewEncryptService(vaultSvc, appCfg.GetVaultCommonTransit())
	encryptorSvc := commonVault.NewEncryptService(vaultSvc, appCfg.GetVaultCommonTransit())

	pluginWrapper := plugin.NewPlugin(appCfg.GetHdWalletPluginPath())
	err = pluginWrapper.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init plugin", zap.Error(err))
	}
	loggerEntry.Info("plugin successfully loaded",
		zap.String(app.PluginNameTag, pluginWrapper.GetPluginName()),
		zap.String(app.PluginReleaseTag, pluginWrapper.GetReleaseTag()),
		zap.Uint64(app.PluginBuildNumberTag, pluginWrapper.GetBuildNumber()),
		zap.Int64(app.PluginBuildDateTag, pluginWrapper.GetBuildDateTS()),
		zap.String(app.PluginCommitIDTag, pluginWrapper.GetCommitID()),
		zap.String(app.PluginShortCommitIDTag, pluginWrapper.GetShortCommitID()))

	walletsPoolSvc := wallet_manager.NewWalletPool(ctx, loggerEntry, appCfg,
		pluginWrapper.GetMakeWalletCallback(), encryptorSvc)
	walletsPoolSvc.Run()

	profiler := commonProfiler.NewHTTPServer(loggerEntry, appCfg.ProfilerConfig)

	apiHandlers := grpc.NewHandlers(loggerEntry,
		pluginWrapper.GetMnemonicGeneratorFunc(),
		pluginWrapper.GetMnemonicValidatorFunc(),
		transitSvc, encryptorSvc, walletsPoolSvc)
	GRPCSrv, err := grpc.NewServer(ctx, loggerEntry, appCfg, apiHandlers)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String(app.GRPCBindPathTag, appCfg.GetConnectionPath()))
	}

	err = GRPCSrv.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to listen init grpc server instance", zap.Error(err),
			zap.String(app.GRPCBindPathTag, appCfg.GetConnectionPath()))
	}

	err = profiler.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init profiler", zap.Error(err))
	}
	loggerEntry.Info("profiler successfully initiated")

	// TODO: add healthcheck flow
	commonHealthcheck.NewHTTPHealthChecker(loggerEntry, appCfg)
	//checker.AddStartupProbeUnit(vaultSvc)
	//checker.AddStartupProbeUnit(redisConn)
	//checker.AddStartupProbeUnit(pgConn)
	//checker.AddStartupProbeUnit(natsConnSvc)

	err = GRPCSrv.ListenAndServe(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to start grpc", zap.Error(err),
			zap.String(app.GRPCBindPathTag, appCfg.GetConnectionPath()))
	}

	err = profiler.ListenAndServe(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init profiler", zap.Error(err))
	}

	loggerEntry.Info("application started successfully",
		zap.String(app.GRPCBindPathTag, appCfg.GetConnectionPath()))

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	loggerEntry.Warn("shutdown application")
	cancelCtxFunc()

	syncErr := loggerEntry.Sync()
	if syncErr != nil {
		log.Print(syncErr.Error(), syncErr)
	}

	log.Print("stopped")
}
