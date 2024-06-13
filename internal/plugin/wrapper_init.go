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

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/app"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
)

func (w *wrapper) Init(_ context.Context) error {
	p, err := plugin.Open(w.pluginPath)
	if err != nil {
		return err
	}

	err = w.initPluginName(p)
	if err != nil {
		return err
	}

	err = w.initLdFlagManager(p)
	if err != nil {
		return err
	}

	err = w.initChainIDFlow(p)
	if err != nil {
		return err
	}

	err = w.initCoinTypeFlow(p)
	if err != nil {
		return err
	}

	err = w.initMnemonicGenerator(p)
	if err != nil {
		return err
	}

	err = w.initMnemonicValidator(p)
	if err != nil {
		return err
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

	w.walletMakerClb = unitMakerFunc

	return nil
}

func (w *wrapper) initPluginName(p *plugin.Plugin) error {
	getPluginNameFunc, err := stringFuncSymbolLookUp(p, getPluginNameSymbol)
	if err != nil {
		return err
	}

	w.pluginName = getPluginNameFunc()

	return nil
}

func (w *wrapper) initChainIDFlow(p *plugin.Plugin) error {
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

	getPluginSupportedChainIDsFunc, ok := getPluginSupportedChainIDsFuncSymbol.(func() string)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGetSupportedChainIDsSymbol)
	}

	setPluginChainIDFuncSymbol, err := p.Lookup(pluginSetChainIDSymbol)
	if err != nil {
		return err
	}

	setPluginChainIDFunc, ok := setPluginChainIDFuncSymbol.(func(int) error)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginSetChainIDSymbol)
	}

	err = setPluginChainIDFunc(w.chainID)
	if err != nil {
		return err
	}

	if getPluginChainIDFunc() != w.chainID {
		return fmt.Errorf("%w: %s, %d - %d", ErrMismatchedPluginValue, app.PluginChainIDTag,
			getPluginChainIDFunc(), w.chainID)
	}

	w.getChainIDFunc = getPluginChainIDFunc
	w.getSupportedChainIDsFunc = getPluginSupportedChainIDsFunc

	return nil
}

func (w *wrapper) initCoinTypeFlow(p *plugin.Plugin) error {
	getPluginCoinTypeFuncSymbol, err := p.Lookup(pluginGetCoinTypeSymbol)
	if err != nil {
		return err
	}

	getPluginCoinTypeFunc, ok := getPluginCoinTypeFuncSymbol.(func() int)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGetCoinTypeSymbol)
	}

	getPluginSupportedCoinTypesFuncSymbol, err := p.Lookup(pluginGetSupportedCoinTypesSymbol)
	if err != nil {
		return err
	}

	getPluginSupportedCoinTypesFunc, ok := getPluginSupportedCoinTypesFuncSymbol.(func() string)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGetSupportedCoinTypesSymbol)
	}

	setPluginCoinTypeFuncSymbol, err := p.Lookup(pluginSetCoinTypeSymbol)
	if err != nil {
		return err
	}

	setPluginCoinTypeFunc, ok := setPluginCoinTypeFuncSymbol.(func(int) error)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginSetCoinTypeSymbol)
	}

	err = setPluginCoinTypeFunc(w.coinType)
	if err != nil {
		return err
	}

	if getPluginCoinTypeFunc() != w.coinType {
		return fmt.Errorf("%w: %s, %d - %d", ErrMismatchedPluginValue, app.PluginCoinTypeTag,
			getPluginCoinTypeFunc(), w.coinType)
	}

	w.getCoinTypeFunc = getPluginCoinTypeFunc
	w.getSupportedCoinTypeFunc = getPluginSupportedCoinTypesFunc

	return nil
}

func (w *wrapper) initMnemonicValidator(p *plugin.Plugin) error {
	validateMnemonicFuncSymbol, err := p.Lookup(pluginValidateMnemonicSymbol)
	if err != nil {
		return err
	}

	validateMnemoFunc, ok := validateMnemonicFuncSymbol.(func(mnemonic string) bool)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginValidateMnemonicSymbol)
	}

	w.validateFunc = validateMnemoFunc

	return nil
}

func (w *wrapper) initMnemonicGenerator(p *plugin.Plugin) error {
	generateMnemonicFuncSymbol, err := p.Lookup(pluginGenerateMnemonicSymbol)
	if err != nil {
		return err
	}

	generateMnemoFunc, ok := generateMnemonicFuncSymbol.(func() (string, error))
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGenerateMnemonicSymbol)
	}

	w.generateFunc = generateMnemoFunc

	return nil
}

func (w *wrapper) initLdFlagManager(p *plugin.Plugin) error {
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

	flagManagerSvc, err := commonConfig.NewLdFlagsManager(getPluginReleaseTagFunc(),
		getPluginCommitIDFunc(), getPluginShortCommitIDFunc(),
		getPluginBuildNumberFunc(), getPluginBuildDateTSFunc())
	if err != nil {
		return err
	}

	w.ldFlagManager = flagManagerSvc

	return nil
}
