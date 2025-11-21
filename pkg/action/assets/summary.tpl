| Name | Previous | Current | Ratio | Verdict |
|------|----------|---------|-------|---------|
{{- range $calc := . }}
{{- printf "| %s | | | | |" $calc.Name }}
{{- range $metric := $calc.MetricDiff }}
{{- $ratio := $metric.Ratio }}
{{- printf "| | (%s) %s | (%s) %s | %s | %s |"
		`({{ $metric.PreviousCommit }})` `{{ printf "%.2f %s" $metric.PreviousValue $metric.Unit }}`
		`({{ $metric.CurrentCommit 0 6}})` `{{ printf "%.2f %s" $metric.CurrentValue $metric.Unit }}`
		`{{ if ne $ratio $ratio }}N/A{{ else }}{{ printf "%.2f" $ratio }}{{ end }}`
		`{{ $metric.Emoji }}`
}}
{{- end }}
{{- end }}
