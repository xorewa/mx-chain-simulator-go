package config

import (
	"os"
	"path/filepath"
	"testing"

	chainConfig "github.com/multiversx/mx-chain-go/config"
	proxyConfig "github.com/multiversx/mx-chain-proxy-go/config"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadNodeOverrideConfigs(t *testing.T) {
	t.Parallel()

	testString := `
# Test comment
OverridableConfigTomlValues = [
    { File = "config.toml", Path = "A", Value = "B" },
    { File = "external.toml", Path = "C", Value = "D" }
]
`

	expectedConfig := OverrideConfigs{
		OverridableConfigTomlValues: []chainConfig.OverridableConfig{
			{
				File:  "config.toml",
				Path:  "A",
				Value: "B",
			},
			{
				File:  "external.toml",
				Path:  "C",
				Value: "D",
			},
		},
	}

	cfg := OverrideConfigs{}

	err := toml.Unmarshal([]byte(testString), &cfg)
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, cfg)
}

func TestBundledChainSimulatorConfigsContainSupernovaRequiredFields(t *testing.T) {
	t.Parallel()

	nodeConfig := chainConfig.Config{}
	mustLoadToml(t, filepath.Join("..", "cmd", "chainsimulator", "config", "node", "config", "config.toml"), &nodeConfig)

	require.GreaterOrEqual(t, len(nodeConfig.GeneralSettings.ChainParametersByEpoch), 3)
	require.GreaterOrEqual(t, len(nodeConfig.Versions.VersionsByEpochs), 3)
	require.NotEmpty(t, nodeConfig.GeneralSettings.ProcessConfigsByEpoch)
	require.NotEmpty(t, nodeConfig.GeneralSettings.ProcessConfigsByRound)
	require.NotEmpty(t, nodeConfig.GeneralSettings.EpochStartConfigsByEpoch)
	require.NotEmpty(t, nodeConfig.GeneralSettings.EpochStartConfigsByRound)
	require.NotEmpty(t, nodeConfig.GeneralSettings.ConsensusConfigsByEpoch)
	require.NotEmpty(t, nodeConfig.Antiflood.ConfigsByRound)
	require.NotZero(t, nodeConfig.TxCacheBounds.MaxNumBytesPerSenderUpperBound)
	require.NotZero(t, nodeConfig.TxCacheBounds.MaxTrackedBlocks)
	require.NotZero(t, nodeConfig.TxCacheSelection.SelectionMaxNumTxs)
	require.NotZero(t, nodeConfig.TxCacheSelection.SelectionGasRequested)
	require.NotZero(t, nodeConfig.PostProcessTransactionsCache.Type)
	require.NotZero(t, nodeConfig.ExecutedMiniBlocksCache.Type)
	require.NotZero(t, nodeConfig.DirectSentTransactions.CacheSpanInSec)
	require.NotZero(t, nodeConfig.DirectSentTransactions.CacheExpiryInSec)
	require.NotZero(t, nodeConfig.ExecutionResultInclusionEstimator.SafetyMargin)
	require.NotZero(t, nodeConfig.ExecutionResultInclusionEstimator.MaxResultsPerBlock)

	epochConfig := chainConfig.EpochConfig{}
	mustLoadToml(t, filepath.Join("..", "cmd", "chainsimulator", "config", "node", "config", "enableEpochs.toml"), &epochConfig)
	require.Equal(t, uint32(2), epochConfig.EnableEpochs.SupernovaEnableEpoch)

	roundConfig := chainConfig.RoundConfig{}
	mustLoadToml(t, filepath.Join("..", "cmd", "chainsimulator", "config", "node", "config", "enableRounds.toml"), &roundConfig)
	require.Equal(t, "440", roundConfig.RoundActivations["SupernovaEnableRound"].Round)

	economicsConfig := chainConfig.EconomicsConfig{}
	mustLoadToml(t, filepath.Join("..", "cmd", "chainsimulator", "config", "node", "config", "economics.toml"), &economicsConfig)
	require.Equal(t, uint64(200), economicsConfig.FeeSettings.BlockCapacityOverestimationFactor)
	require.Equal(t, uint64(10), economicsConfig.FeeSettings.PercentDecreaseLimitsStep)

	proxyCfg := proxyConfig.Config{}
	mustLoadToml(t, filepath.Join("..", "cmd", "chainsimulator", "config", "proxy", "config", "config.toml"), &proxyCfg)
	require.NotZero(t, proxyCfg.GeneralSettings.BlockCacheDurationSec)
}

func mustLoadToml(t *testing.T, filename string, target interface{}) {
	t.Helper()

	bytes, err := os.ReadFile(filename)
	require.NoError(t, err)
	require.NoError(t, toml.Unmarshal(bytes, target))
}
