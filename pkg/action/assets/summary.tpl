| Name | Previous | Current | Ratio | Verdict |
|------|----------|---------|-------|---------|
{{- range $calc := . }}
{{ printf "| `%s` | `%s` | `%s` | ---- | ---- |"
		$calc.Name
		$calc.PreviousCommit
		$calc.CurrentCommit
}}
{{- range $metric := $calc.MetricDiff }}
{{ printf "| ---- | `%s` | `%s` | `%s` | `%s` |"
		(printf "%.2f %s" $metric.PreviousValue $metric.Unit)
		(printf "%.2f %s" $metric.CurrentValue $metric.Unit)
		(printf "%.2f" $metric.Ratio)
		$metric.Emoji
}}
{{- end }}
{{- end }}
