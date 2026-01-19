### 🔴 Performance Degradation

The following benchmarks have **degraded** compared to the previous run:

| Name | Previous | Current | Ratio | Verdict |
|------|----------|---------|-------|---------|
{{- range $calc := . }}
{{ printf "| `%s` | `%s` | `%s` | | |"
		$calc.Name
		(slice $calc.PreviousCommit 0 6)
		(slice $calc.CurrentCommit 0 6)
}}
{{- range $metric := $calc.MetricDiff }}
{{ printf "| | `%s` | `%s` | `%s` | `%s` |"
		(printf "%.2f %s" $metric.PreviousValue $metric.Unit)
		(printf "%.2f %s" $metric.CurrentValue $metric.Unit)
		(printf "%.2f" $metric.Ratio)
		$metric.Emoji
}}
{{- end }}
{{- end }}
