{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a list of values for ADVERTISE_IP
*/}}
{{- define "advertiseIp" -}}
{{- $hosts := list -}}
{{- range .Values.ingress.hosts -}}
{{- if kindIs "string" . }}
    {{- $hosts = printf "https://%s" . | append $hosts -}}
    {{- $hosts = printf "https://%s:443" . | append $hosts -}}
{{- else }}
    {{- $hosts = printf "https://%s" .host | append $hosts -}}
    {{- $hosts = printf "https://%s:443" .host | append $hosts -}}
{{- end }}
{{- end -}}
{{ join "," $hosts }}
{{- end -}}
