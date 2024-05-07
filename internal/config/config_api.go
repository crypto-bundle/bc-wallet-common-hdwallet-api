package config

import (
	"fmt"
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"

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
	*pbApi.HdWalletClientConfig // yes, client config for listen on unix file socket
	*ProcessionEnvironmentConfig
	*VaultWrappedConfig
	// -------------------
	// Internal configs
	// -------------------
	*MnemonicConfig
	// VaultCommonTransitKey - common vault transit key for whole processing cluster
	VaultCommonTransitKey string `envconfig:"VAULT_COMMON_TRANSIT_KEY" default:"-"`
	// VaultApplicationEncryptionKey - vault encryption key for hd-wallet-controller and hd-wallet-api application
	VaultApplicationEncryptionKey string `envconfig:"VAULT_APP_ENCRYPTION_KEY" default:"-"`
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
