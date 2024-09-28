{{- define "istio-upgrade-consumer.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "istio-upgrade-consumer.podlabels" -}}
{{- with .Values.podLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}
