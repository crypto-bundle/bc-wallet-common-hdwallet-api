/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package config

import (
	"fmt"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonProfiler "github.com/crypto-bundle/bc-wallet-common-lib-profiler/pkg/profiler"
)

// HdWalletConfig for application
type HdWalletConfig struct {
	// -------------------
	// External common configs
	// -------------------
	*commonHealthcheck.HealthcheckHTTPConfig
	*commonProfiler.ProfilerConfig
	*pbApi.HdWalletClientConfig // yes, client config for listen on unix file socket
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
	baseAppCfgSvc       baseConfigService
	loggerCfgSvc        loggerCfgService
	processingEnvCfgSvc processingEnvironmentConfigService
}

func (c *HdWalletConfig) GetVaultCommonTransit() string {
	return c.VaultCommonTransitKey
}

func (c *HdWalletConfig) GetVaultAppEncryptionKey() string {
	return c.VaultApplicationEncryptionKey
}

// Prepare variables to static configuration
func (c *HdWalletConfig) Prepare() error {
	appName := fmt.Sprintf(ApplicationManagerNameTpl, c.processingEnvCfgSvc.GetNetworkName())

	c.baseAppCfgSvc.SetApplicationName(appName)

	return nil
}

func (c *HdWalletConfig) PrepareWith(cfgSvcList ...interface{}) error {
	for _, cfgSrv := range cfgSvcList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSvc = castedCfg
		case processingEnvironmentConfigService:
			c.processingEnvCfgSvc = castedCfg
		case loggerCfgService:
			c.loggerCfgSvc = castedCfg
		default:
			continue
		}
	}

	return nil
}
