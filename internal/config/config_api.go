package config

import (
	"fmt"
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonProfiler "github.com/crypto-bundle/bc-wallet-common-lib-profiler/pkg/profiler"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
)

// HdWalletConfig for application
type HdWalletConfig struct {
	// -------------------
	// External common configs
	// -------------------
	*commonConfig.BaseConfig
	*commonLogger.LoggerConfig
	*commonHealthcheck.HealthcheckHTTPConfig
	*commonProfiler.ProfilerConfig
	*pbApi.HdWalletClientConfig // yes, client config for listen on unix file socket
	*ProcessionEnvironmentConfig
	*VaultWrappedConfig
	// -------------------
	// Internal configs
	// -------------------
	*MnemonicConfig
	// VaultCommonTransitKey - common vault transit key for whole processing cluster,
	// must be saved in common vault kv bucket, for example: crypto-bundle/bc-wallet-common/transit
	VaultCommonTransitKey string `envconfig:"VAULT_COMMON_TRANSIT_KEY" secret:"true"`
	// VaultApplicationEncryptionKey - vault encryption key for hd-wallet-controller and hd-wallet-api application,
	// must be saved in bc-wallet-<blockchain_name>-hdwallet vault kv bucket,
	// for example: crypto-bundle/bc-wallet-tron-hdwallet/common
	VaultApplicationEncryptionKey string `envconfig:"VAULT_APP_ENCRYPTION_KEY" secret:"true"`
	// ----------------------------
	// Dependencies
	baseAppCfgSrv baseConfigService
}

func (c *HdWalletConfig) GetVaultCommonTransit() string {
	return c.VaultCommonTransitKey
}

func (c *HdWalletConfig) GetVaultAppEncryptionKey() string {
	return c.VaultApplicationEncryptionKey
}

// Prepare variables to static configuration
func (c *HdWalletConfig) Prepare() error {
	appName := fmt.Sprintf(ApplicationManagerNameTpl, c.ProcessionEnvironmentConfig.GetNetworkName())

	c.baseAppCfgSrv.SetApplicationName(appName)

	return nil
}

func (c *HdWalletConfig) PrepareWith(cfgSvcList ...interface{}) error {
	for _, cfgSrv := range cfgSvcList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSrv = castedCfg
		default:
			continue
		}
	}

	return nil
}
