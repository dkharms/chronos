| Name | Unit | Previous Commit | Current Commit | Previous Value | Current Value | Ratio |
|------|------|-----------------|----------------|----------------|---------------|-------|
{{- range $calc := . }}
{{- range $metric := $calc.MetricDiff }}
{{- $ratio := $metric.Ratio }}
| {{ $calc.Name }} | {{ $metric.Unit }} | {{ $metric.PreviousCommit }} | {{ $metric.CurrentCommit }} | {{ printf "%.2f" $metric.PreviousValue }} | {{ printf "%.2f" $metric.CurrentValue }} | {{ if ne $ratio $ratio }}N/A{{ else }}{{ printf "%.2f" $ratio }}{{ end }} |
{{- end }}
{{- end }}
