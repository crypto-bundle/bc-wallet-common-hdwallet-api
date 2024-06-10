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

import "errors"

var (
	ErrWrongMnemonicWordsCount = errors.New("wrong mnemonic words count.allowed: 15, 18, 21, 24")
	ErrBaseAppConfigIsEmpty    = errors.New("base app config is missing or empty")
	ErrMinWordsCount           = errors.New("minimal words count for environment is 21")
)

type MnemonicConfig struct {
	MnemonicWordCount uint8  `envconfig:"HDWALLET_WORDS_COUNT" default:"24"`
	ChainID           int    `envconfig:"HDWALLET_CHAIN_ID" default:"-1"`
	PluginPath        string `envconfig:"HDWALLET_PLUGIN_PATH" default:"/usr/local/bin/hdwallet_plugin.so"`

	baseAppCfgSrv baseConfigService
}

func (c *MnemonicConfig) GetDefaultMnemonicWordsCount() uint8 {
	return c.MnemonicWordCount
}
func (c *MnemonicConfig) GetHdWalletChainID() int {
	return c.ChainID
}
func (c *MnemonicConfig) GetHdWalletPluginPath() string {
	return c.PluginPath
}

//nolint:funlen // its ok
func (c *MnemonicConfig) Prepare() error {
	if c.baseAppCfgSrv == nil {
		return ErrBaseAppConfigIsEmpty
	}

	if c.baseAppCfgSrv.IsProd() && c.MnemonicWordCount <= 18 {
		return ErrMinWordsCount
	}

	switch c.MnemonicWordCount {
	case 15, 18, 21, 24:
		return nil
	default:
		return ErrWrongMnemonicWordsCount
	}
}

//nolint:funlen // its ok
func (c *MnemonicConfig) PrepareWith(dependentCfgList ...interface{}) error {
	for _, cfgSrv := range dependentCfgList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSrv = castedCfg
		default:
			continue
		}
	}

	return nil
}
