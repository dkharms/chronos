| Name | Previous | Current | Ratio | Verdict |
|------|----------|---------|-------|---------|
{{- range $calc := . }}
{{ printf "| `%s` | | | | |" $calc.Name }}
{{- range $metric := $calc.MetricDiff }}
{{- $ratio := $metric.Ratio }}
{{ printf "| | `%s` `%s` | `%s` `%s` | `%s` | `%s` |"
		(slice $metric.PreviousCommit 0 6) (printf "%.2f %s" $metric.PreviousValue $metric.Unit)
		(slice $metric.CurrentCommit 0 6) (printf "%.2f %s" $metric.CurrentValue $metric.Unit)
		(printf "%.2f" $ratio)
		$metric.Emoji
}}
{{- end }}
{{- end }}
