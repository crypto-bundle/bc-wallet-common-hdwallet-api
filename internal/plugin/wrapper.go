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

package plugin

import (
	"context"
	"fmt"
	"plugin"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
)

const (
	getPluginNameSymbol          = "GetPluginName"
	getPluginReleaseTagSymbol    = "GetPluginReleaseTag"
	getPluginCommitIDSymbol      = "GetPluginCommitID"
	getPluginShortCommitIDSymbol = "GetPluginShortCommitID"
	getPluginBuildNumberSymbol   = "GetPluginBuildNumber"
	getPluginBuildDateTSSymbol   = "GetPluginBuildDateTS"

	pluginGetChainIDSymbol           = "GetChainID"
	pluginGetSupportedChainIDsSymbol = "GetSupportedChainIDs"

	pluginGenerateMnemonicSymbol = "GenerateMnemonic"
	pluginValidateMnemonicSymbol = "ValidateMnemonic"
	pluginNewPoolUnitSymbol      = "NewPoolUnit"
)

type wrapper struct {
	pluginPath string
	pluginName string

	generateFunc   generateMnemonicFunction
	validateFunc   validateMnemonicFunction
	walletMakerClb walletMakerFunc

	getChainIDFunc           getChainIDFunction
	getSupportedChainIDsFunc getSupportedChainIDsFunction

	ldFlagManager
}

func (w *wrapper) GetPluginName() string {
	return w.pluginName
}

func (w *wrapper) GetChainID() int {
	return w.getChainIDFunc()
}

func (w *wrapper) GetSupportedChainIDs() []int {
	return w.getSupportedChainIDsFunc()
}

func (w *wrapper) GetMnemonicGeneratorFunc() func() (string, error) {
	return w.generateFunc
}

func (w *wrapper) GetMnemonicValidatorFunc() func(mnemonic string) bool {
	return w.validateFunc
}

func (w *wrapper) GetMakeWalletCallback() func(walletUUID string,
	mnemonicDecryptedData string,
) (interface{}, error) {
	return w.walletMakerClb
}

func (w *wrapper) Init(_ context.Context) error {
	p, err := plugin.Open(w.pluginPath)
	if err != nil {
		return err
	}

	getPluginNameFunc, err := stringFuncSymbolLookUp(p, getPluginNameSymbol)
	if err != nil {
		return err
	}

	getPluginReleaseTagFunc, err := stringFuncSymbolLookUp(p, getPluginReleaseTagSymbol)
	if err != nil {
		return err
	}

	getPluginCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginCommitIDSymbol)
	if err != nil {
		return err
	}

	getPluginShortCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginShortCommitIDSymbol)
	if err != nil {
		return err
	}

	getPluginBuildNumberFunc, err := stringFuncSymbolLookUp(p, getPluginBuildNumberSymbol)
	if err != nil {
		return err
	}

	getPluginBuildDateTSFunc, err := stringFuncSymbolLookUp(p, getPluginBuildDateTSSymbol)
	if err != nil {
		return err
	}

	getPluginChainIDFuncSymbol, err := p.Lookup(pluginGetChainIDSymbol)
	if err != nil {
		return err
	}

	getPluginChainIDFunc, ok := getPluginChainIDFuncSymbol.(func() int)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGetChainIDSymbol)
	}

	getPluginSupportedChainIDsFuncSymbol, err := p.Lookup(pluginGetSupportedChainIDsSymbol)
	if err != nil {
		return err
	}

	getPluginSupportedChainIDsFunc, ok := getPluginSupportedChainIDsFuncSymbol.(func() []int)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGetSupportedChainIDsSymbol)
	}

	generateMnemonicFuncSymbol, err := p.Lookup(pluginGenerateMnemonicSymbol)
	if err != nil {
		return err
	}

	generateMnemoFunc, ok := generateMnemonicFuncSymbol.(func() (string, error))
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGenerateMnemonicSymbol)
	}

	validateMnemonicFuncSymbol, err := p.Lookup(pluginValidateMnemonicSymbol)
	if err != nil {
		return err
	}

	validateMnemoFunc, ok := validateMnemonicFuncSymbol.(func(mnemonic string) bool)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginValidateMnemonicSymbol)
	}

	unitMakerFuncSymbol, err := p.Lookup(pluginNewPoolUnitSymbol)
	if err != nil {
		return err
	}

	unitMakerFunc, ok := unitMakerFuncSymbol.(func(walletUUID string,
		mnemonicDecryptedData string,
	) (interface{}, error))
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginNewPoolUnitSymbol)
	}

	flagManagerSvc, err := commonConfig.NewLdFlagsManager(getPluginReleaseTagFunc(),
		getPluginCommitIDFunc(), getPluginShortCommitIDFunc(),
		getPluginBuildNumberFunc(), getPluginBuildDateTSFunc())
	if err != nil {
		return err
	}

	w.getChainIDFunc = getPluginChainIDFunc
	w.getSupportedChainIDsFunc = getPluginSupportedChainIDsFunc
	w.generateFunc = generateMnemoFunc
	w.validateFunc = validateMnemoFunc
	w.ldFlagManager = flagManagerSvc
	w.pluginName = getPluginNameFunc()
	w.walletMakerClb = unitMakerFunc

	return nil
}

func NewPlugin(pluginPath string) *wrapper {
	return &wrapper{
		pluginPath:     pluginPath,
		pluginName:     "",
		walletMakerClb: nil,
		ldFlagManager:  nil,
	}
}
