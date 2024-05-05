package app

const (
	ApplicationNameTag = "application"
	BlockChainNameTag  = "blockchain_name"

	WalletUUIDTag         = "wallet_uuid"
	MnemonicWalletUUIDTag = "mnemonic_wallet_uuid"
	MnemonicWalletHashTag = "mnemonic_wallet_hash"
	WalletPurposeTag      = "purpose"
	WalletIsHotTag        = "is_hot"

	GRPCBindPathTag = "grpc_bind_unix_socket_path"

	TickerEventTriggerTimeTag = "ticker_time"

	HDWalletAccountIndexTag  = "hd_wallet_account_index"
	HDWalletInternalIndexTag = "hd_wallet_internal_index"
	HDWalletAddressIndexTag  = "hd_wallet_address_index"
	HDWalletAddressTag       = "hd_wallet_address"

	NatsCacheBucketNameTag = "nats_kv_bucket_name"

	PluginNameTag          = "plugin_name"
	PluginReleaseTag       = "plugin_release_tag"
	PluginCommitIDTag      = "plugin_commit_id_tag"
	PluginShortCommitIDTag = "plugin_short_commit_id_tag"
	PluginBuildNumberTag   = "plugin_build_number_tag"
	PluginBuildDateTag     = "plugin_build_date_tag"
)
