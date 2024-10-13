package templates_test

import (
	"github.com/YReshetko/go-annotation/annotations/cobra/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestExecute_InitCommands(t *testing.T) {
	testcases := []struct {
		name         string
		model        templates.InitCommands
		expectedData string
	}{
		{
			name: "with tag",
			model: templates.InitCommands{
				BuildTag: "test",
				Imports: []templates.Import{
					{Alias: "a1", Package: "git.com/some/path"},
					{Alias: "a2", Package: "git.com/some/another/path"},
				},
				Commands: []templates.Command{
					{
						IsRoot:        true,
						VarName:       "root",
						Use:           "use",
						Example:       "example",
						Short:         "short",
						Long:          "long",
						SilenceUsage:  true,
						SilenceErrors: true,
						Flags: []templates.Flag{
							{
								Type:         templates.StringFlagType,
								Name:         "flag1",
								Shorthand:    "f",
								DefaultValue: "10",
								Description:  "my flag 1",
								IsRequired:   false,
								IsPersistent: false,
							},
							{
								Type:         templates.IntFlagType,
								Name:         "flag2",
								DefaultValue: "10",
								Description:  "my flag 2",
								IsRequired:   false,
								IsPersistent: false,
							},
							{
								Type:         templates.BoolFlagType,
								Name:         "flag3",
								Shorthand:    "f",
								DefaultValue: "true",
								Description:  "my flag 3",
								IsRequired:   true,
								IsPersistent: false,
							},
							{
								Type:         templates.Float32FlagType,
								Name:         "flag4",
								Shorthand:    "f",
								DefaultValue: "3.14",
								Description:  "my flag 4",
								IsRequired:   false,
								IsPersistent: true,
							},
							{
								Type:         templates.Uint32FlagType,
								Name:         "flag5",
								DefaultValue: "3",
								Description:  "my flag 5",
								IsRequired:   true,
								IsPersistent: true,
							},
						},
					},
					{
						VarName:       "child",
						ParentVarName: "root",
						Handlers: []templates.Handler{
							{
								MethodName:           "Run1",
								ExecutorPackageAlias: "a1",
								ExecutorTypeName:     "RootCommandExecutor",
								IsPreRun:             false,
								IsPostRun:            false,
								IsPersistentRun:      false,
								HasReturn:            false,
							},
							{
								MethodName:           "Run2",
								ExecutorPackageAlias: "a1",
								ExecutorTypeName:     "RootCommandExecutor",
								IsPreRun:             true,
								IsPostRun:            false,
								IsPersistentRun:      false,
								HasReturn:            false,
							},
							{
								MethodName:           "Run3",
								ExecutorPackageAlias: "a1",
								ExecutorTypeName:     "RootCommandExecutor",
								IsPreRun:             false,
								IsPostRun:            true,
								IsPersistentRun:      false,
								HasReturn:            false,
							},
							{
								MethodName:           "Run4",
								ExecutorPackageAlias: "a1",
								ExecutorTypeName:     "RootCommandExecutor",
								IsPreRun:             false,
								IsPostRun:            true,
								IsPersistentRun:      true,
								HasReturn:            false,
							},
							{
								MethodName:           "Run5",
								ExecutorPackageAlias: "a1",
								ExecutorTypeName:     "RootCommandExecutor",
								IsPreRun:             true,
								IsPostRun:            false,
								IsPersistentRun:      true,
								HasReturn:            true,
							},
						},
					},
				},
			},
			expectedData: "./fixtures/init-command-with-tag.go.txt",
		},
		{
			name: "without tag",
			model: templates.InitCommands{
				Imports: []templates.Import{
					{Alias: "", Package: "git.com/some/path"},
					{Alias: "", Package: "git.com/some/another/path"},
				},
			},
			expectedData: "./fixtures/init-command-without-tag.go.txt",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			expectedData := loadFile(t, testcase.expectedData)
			actualData, err := templates.Execute(templates.InitCommandsTpl, testcase.model)
			require.NoError(t, err)
			assert.Equal(t, string(expectedData), string(actualData))
		})
	}
}

func loadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}
