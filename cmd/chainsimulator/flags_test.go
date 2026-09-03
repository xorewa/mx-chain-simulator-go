package main

import (
	"flag"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestDetermineOverrideConfigFiles(t *testing.T) {
	tests := []struct {
		name          string
		overrideFiles string
		skipDefault   bool
		expected      []string
	}{
		{
			name:          "default remains unchanged",
			overrideFiles: nodeOverrideDefaultPath,
			expected:      []string{nodeOverrideDefaultPath},
		},
		{
			name:          "default is prepended to custom overrides",
			overrideFiles: "first.toml, second.toml",
			expected:      []string{nodeOverrideDefaultPath, "first.toml", "second.toml"},
		},
		{
			name:          "production config can exclude bundled override",
			overrideFiles: nodeOverrideDefaultPath,
			skipDefault:   true,
			expected:      []string{},
		},
		{
			name:          "explicit custom overrides remain when bundled override is skipped",
			overrideFiles: nodeOverrideDefaultPath + ", production.toml",
			skipDefault:   true,
			expected:      []string{"production.toml"},
		},
		{
			name:          "empty entries are ignored",
			overrideFiles: " ,production.toml, ",
			skipDefault:   true,
			expected:      []string{"production.toml"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet(test.name, flag.ContinueOnError)
			flagSet.String(nodeOverrideConfigurationFile.Name, "", "")
			flagSet.Bool(skipDefaultNodeOverride.Name, false, "")
			require.NoError(t, flagSet.Set(nodeOverrideConfigurationFile.Name, test.overrideFiles))
			if test.skipDefault {
				require.NoError(t, flagSet.Set(skipDefaultNodeOverride.Name, "true"))
			}

			ctx := cli.NewContext(cli.NewApp(), flagSet, nil)
			require.Equal(t, test.expected, determineOverrideConfigFiles(ctx))
		})
	}
}

func TestHasModuleReplacement(t *testing.T) {
	buildInfo := &debug.BuildInfo{
		Deps: []*debug.Module{
			{
				Path:    mxChainNodeModulePath,
				Version: "v1.11.1",
				Replace: &debug.Module{
					Path:    "github.com/xorewa/mx-chain-go",
					Version: "v0.0.0-20260903052749-1b527afd34ff",
				},
			},
		},
	}

	require.True(t, hasModuleReplacement(buildInfo, mxChainNodeModulePath))
	require.ErrorIs(t, validateAutomaticNodeConfigDownload(buildInfo), errAutomaticNodeConfigDownloadWithReplacement)
	require.False(t, hasModuleReplacement(buildInfo, "github.com/multiversx/mx-chain-proxy-go"))
	withoutReplacement := &debug.BuildInfo{
		Deps: []*debug.Module{{Path: mxChainNodeModulePath, Version: "v1.11.8"}},
	}
	require.False(t, hasModuleReplacement(withoutReplacement, mxChainNodeModulePath))
	require.NoError(t, validateAutomaticNodeConfigDownload(withoutReplacement))
}
