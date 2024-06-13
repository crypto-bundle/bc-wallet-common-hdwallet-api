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

const (
	getPluginNameSymbol          = "GetPluginName"
	getPluginReleaseTagSymbol    = "GetPluginReleaseTag"
	getPluginCommitIDSymbol      = "GetPluginCommitID"
	getPluginShortCommitIDSymbol = "GetPluginShortCommitID"
	getPluginBuildNumberSymbol   = "GetPluginBuildNumber"
	getPluginBuildDateTSSymbol   = "GetPluginBuildDateTS"

	pluginGetChainIDSymbol           = "GetChainID"
	pluginGetSupportedChainIDsSymbol = "GetSupportedChainIDsInfo"
	pluginSetChainIDSymbol           = "SetChainID"

	pluginGetCoinTypeSymbol           = "GetHdWalletCoinType"
	pluginGetSupportedCoinTypesSymbol = "GetSupportedCoinTypesInfo"
	pluginSetCoinTypeSymbol           = "SetHdWalletCoinType"

	pluginGenerateMnemonicSymbol = "GenerateMnemonic"
	pluginValidateMnemonicSymbol = "ValidateMnemonic"
	pluginNewPoolUnitSymbol      = "NewPoolUnit"
)

type wrapper struct {
	pluginPath string
	pluginName string

	coinType int
	chainID  int

	generateFunc   generateMnemonicFunction
	validateFunc   validateMnemonicFunction
	walletMakerClb walletMakerFunc

	getChainIDFunc           getChainIDFunction
	getSupportedChainIDsFunc getSupportedChainIDsFunction
	//setChainIDFunc           setChainIDFunction

	getCoinTypeFunc          getCoinTypeFunction
	getSupportedCoinTypeFunc getSupportedCoinTypesFunction
	//setCoinTypeFunc          setCoinTypeFunction

	ldFlagManager
}

func (w *wrapper) GetPluginName() string {
	return w.pluginName
}

func (w *wrapper) GetChainID() int {
	return w.getChainIDFunc()
}

func (w *wrapper) GetSupportedChainIDs() string {
	return w.getSupportedChainIDsFunc()
}

//func (w *wrapper) SetChainID(chainID int) error {
//	return w.setChainIDFunc(chainID)
//}

func (w *wrapper) GetCoinType() int {
	return w.getCoinTypeFunc()
}

func (w *wrapper) GetSupportedCoinTypesInfo() string {
	return w.getSupportedCoinTypeFunc()
}

//func (w *wrapper) SetCoinType(coinType int) error {
//	return w.setCoinTypeFunc(coinType)
//}

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

func NewPlugin(pluginPath string,
	coinType int,
	chainID int,
) *wrapper {
	return &wrapper{
		pluginPath: pluginPath,
		coinType:   coinType,
		chainID:    chainID,

		pluginName:     "",
		walletMakerClb: nil,
		ldFlagManager:  nil,
	}
}
