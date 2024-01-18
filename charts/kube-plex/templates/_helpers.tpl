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
{{- if (or (eq .Values.service.type "ClusterIP") (eq .Values.service.type "LoadBalancer") (empty .Values.service.type)) -}}
  {{- if .Values.service.loadBalancerIP -}}{{- $hosts = printf "https://%s:32400" .Values.service.loadBalancerIP | append $hosts -}}{{- end -}}
  {{- if index .Values.service.annotations "dns.pfsense.org/hostname" -}}
    {{- range splitList "," (index .Values.service.annotations "dns.pfsense.org/hostname") -}}
      {{- $hosts = printf "https://%s:32400" . | append $hosts -}}
    {{- end -}}
  {{- end -}}
  {{- range .Values.ingress.hosts -}}
  {{- if kindIs "string" . }}
    {{- if not (eq (lower .) "chart-example.local") -}}
      {{- $hosts = printf "https://%s" . | append $hosts -}}
      {{- $hosts = printf "https://%s:443" . | append $hosts -}}
    {{- end -}}
  {{- else }}
    {{- if not (eq (lower .host) "chart-example.local") -}}
      {{- $hosts = printf "https://%s" .host | append $hosts -}}
      {{- $hosts = printf "https://%s:443" .host | append $hosts -}}
    {{- end -}}
  {{- end }}
  {{- end -}}
{{- end -}}
{{ join "," $hosts }}
{{- end -}}
