package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestDetermineOverrideConfigFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		overrideFiles       string
		skipDefaultOverride bool
		expected            []string
	}{
		{
			name:          "prepends default for existing callers",
			overrideFiles: "supernova.toml",
			expected:      []string{nodeOverrideDefaultPath, "supernova.toml"},
		},
		{
			name:                "explicit profile can exclude default bootstrap",
			overrideFiles:       "supernova.toml",
			skipDefaultOverride: true,
			expected:            []string{"supernova.toml"},
		},
		{
			name:          "explicit default remains authoritative",
			overrideFiles: "custom.toml," + nodeOverrideDefaultPath,
			expected:      []string{"custom.toml", nodeOverrideDefaultPath},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			flagSet := flag.NewFlagSet(test.name, flag.ContinueOnError)
			flagSet.String(nodeOverrideConfigurationFile.Name, "", "")
			flagSet.Bool(skipDefaultNodeOverride.Name, false, "")
			require.NoError(t, flagSet.Set(nodeOverrideConfigurationFile.Name, test.overrideFiles))
			if test.skipDefaultOverride {
				require.NoError(t, flagSet.Set(skipDefaultNodeOverride.Name, "true"))
			}

			ctx := cli.NewContext(cli.NewApp(), flagSet, nil)
			require.Equal(t, test.expected, determineOverrideConfigFiles(ctx))
		})
	}
}

func TestDefaultNodeOverridePathExists(t *testing.T) {
	t.Parallel()

	// Locate the bundled file from the source tree because `go test` runs in a
	// temporary working directory. The CLI resolves nodeOverrideDefaultPath from
	// its documented cmd/chainsimulator working directory.
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	_, err := os.Stat(filepath.Join(filepath.Dir(thisFile), "config", nodeOverrideDefaultFilename))
	require.NoError(t, err)
}
