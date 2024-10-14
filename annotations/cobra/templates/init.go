package templates

var initCommands = `
{{- if .BuildTag -}}
//go:build {{ .BuildTag }}
{{ end -}}

package main

import (
	"github.com/spf13/cobra"
	{{ range .Imports -}}
	{{ if .Alias}}{{ .Alias }} {{ end }}"{{ .Package }}"
	{{ end -}}
)

func init() {
	{{ range .Commands -}}
	{{ .VarName }} {{ if .IsRoot }}={{ else }}:={{ end }} &cobra.Command{}
	{{ .VarName }}.Use = "{{ .Use }}"
	{{ if .Example -}}{{ .VarName }}.Example = "{{ .Example }}"{{ end }}
	{{ if .Short -}}{{ .VarName }}.Short = "{{ .Short }}"{{ end }}
	{{ if .Long -}}{{ .VarName }}.Long = "{{ .Long }}"{{ end }}
	{{ if .SilenceUsage -}}{{ .VarName }}.SilenceUsage = true{{ end }}
	{{ if .SilenceErrors -}}{{ .VarName }}.SilenceErrors = true{{ end }}

	{{- $varName := .VarName -}}
	{{ range .Flags }}
	{{ if .IsPersistent -}}
		{{ $varName }}.PersistentFlags().{{ .Type }}
	{{- else -}}
		{{ $varName }}.Flags().{{ .Type }}
	{{- end -}}
	{{- if	.Shorthand -}}
		P("{{ .Name }}", "{{ .Shorthand }}", {{ if eq .Type "String" }}"{{ .DefaultValue }}"{{ else }}{{ .DefaultValue }}{{ end }}, "{{ .Description }}")
	{{ else -}}
		("{{ .Name }}", {{ if eq .Type "String" }}"{{ .DefaultValue }}"{{ else }}{{ if eq .Type "Bool" }}{{ eq .DefaultValue "true" }}{{ else }}{{ .DefaultValue }}{{ end }}{{ end }}, "{{ .Description }}")
	{{- end }}
	{{ if .IsRequired -}}
	if err := {{ $varName }}.{{ if .IsPersistent }}MarkPersistentFlagRequired{{ else }}MarkFlagRequired{{ end }}("{{ .Name }}"); err != nil {
		fatal(err)
	}
	{{- end }}
	{{- end }}

	{{ range .Handlers }}
	{{ $varName }}.
	{{- if .IsPersistentRun }}Persistent{{ end }}
	{{- if .IsPreRun }}Pre{{ else }}{{ if .IsPostRun }}Post{{ end }}{{ end }}RunE = func(cmd *cobra.Command, args []string) error {
		executor := {{ .ExecutorPackageAlias }}.{{ .ExecutorTypeName }}{}
		if err := {{ if .IsPersistentRun }}parsePersistFlags{{ else }}parseFlags{{ end }}(cmd, &executor); err != nil {
			return err
		}
		{{ if .HasReturn -}}
		return executor.{{ .MethodName }}(cmd, args)
		{{- else -}}
		executor.{{ .MethodName }}(cmd, args)
		return nil
		{{- end }}
	}
	{{- end }}
	
	{{ if .ParentVarName }}{{ .ParentVarName }}.AddCommand({{ .VarName }}){{ end }}
	{{ end }}
}
`
