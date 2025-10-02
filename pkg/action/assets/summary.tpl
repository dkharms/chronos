| Name | Previous | Current | Ratio |
|------|----------|---------|-------|
{{- range $calc := . }}
{{- range $metric := $calc.MetricDiff }}
{{- $ratio := $metric.Ratio }}
| `{{ $calc.Name }}` | ({{ slice $metric.PreviousCommit 0 6}}) `{{ printf "%.2f" $metric.PreviousValue }} {{ $metric.Unit }}` | ({{ slice $metric.CurrentCommit 0 6}}) `{{ printf "%.2f" $metric.CurrentValue }} {{ $metric.Unit }}` | {{ if ne $ratio $ratio }}N/A{{ else }}{{ printf "%.2f" $ratio }}{{ end }} |
{{- end }}
{{- end }}
