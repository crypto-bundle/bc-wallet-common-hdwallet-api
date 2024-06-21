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

import "time"

func (c *HdWalletConfig) GetHostName() string {
	return c.baseAppCfgSvc.GetHostName()
}
func (c *HdWalletConfig) GetEnvironmentName() string {
	return c.baseAppCfgSvc.GetEnvironmentName()
}
func (c *HdWalletConfig) IsProd() bool {
	return c.baseAppCfgSvc.IsProd()
}
func (c *HdWalletConfig) IsStage() bool {
	return c.baseAppCfgSvc.IsStage()
}
func (c *HdWalletConfig) IsTest() bool {
	return c.baseAppCfgSvc.IsTest()
}
func (c *HdWalletConfig) IsDev() bool {
	return c.baseAppCfgSvc.IsDev()
}
func (c *HdWalletConfig) IsDebug() bool {
	return c.baseAppCfgSvc.IsDebug()
}
func (c *HdWalletConfig) IsLocal() bool {
	return c.baseAppCfgSvc.IsLocal()
}
func (c *HdWalletConfig) GetStageName() string {
	return c.baseAppCfgSvc.GetStageName()
}
func (c *HdWalletConfig) GetApplicationPID() int {
	return c.baseAppCfgSvc.GetApplicationPID()
}
func (c *HdWalletConfig) GetApplicationName() string {
	return c.baseAppCfgSvc.GetApplicationName()
}
func (c *HdWalletConfig) SetApplicationName(appName string) {
	c.baseAppCfgSvc.SetApplicationName(appName)
}

func (c *HdWalletConfig) GetReleaseTag() string {
	return c.baseAppCfgSvc.GetReleaseTag()
}
func (c *HdWalletConfig) GetCommitID() string {
	return c.baseAppCfgSvc.GetCommitID()
}
func (c *HdWalletConfig) GetShortCommitID() string {
	return c.baseAppCfgSvc.GetShortCommitID()
}
func (c *HdWalletConfig) GetBuildNumber() uint64 {
	return c.baseAppCfgSvc.GetBuildNumber()
}
func (c *HdWalletConfig) GetBuildDateTS() int64 {
	return c.baseAppCfgSvc.GetBuildDateTS()
}
func (c *HdWalletConfig) GetBuildDate() time.Time {
	return c.baseAppCfgSvc.GetBuildDate()
}

func (c *HdWalletConfig) GetMinimalLogLevel() string {
	return c.loggerCfgSvc.GetMinimalLogLevel()
}
func (c *HdWalletConfig) IsStacktraceEnabled() bool {
	return c.loggerCfgSvc.IsStacktraceEnabled()
}
func (c *HdWalletConfig) GetSkipBuildInfo() bool {
	return c.loggerCfgSvc.GetSkipBuildInfo()
}

// GetProviderName is for getting event filter by processing provider
func (c *HdWalletConfig) GetProviderName() string {
	return c.processingEnvCfgSvc.GetProviderName()
}

// GetNetworkName is for getting event filter by processing network
func (c *HdWalletConfig) GetNetworkName() string {
	return c.processingEnvCfgSvc.GetNetworkName()
}
