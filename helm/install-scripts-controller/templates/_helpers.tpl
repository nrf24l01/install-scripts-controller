{{- define "install-scripts-controller.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "install-scripts-controller.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "install-scripts-controller.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | quote }}
app.kubernetes.io/name: {{ include "install-scripts-controller.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end -}}

{{- define "install-scripts-controller.selectorLabels" -}}
app.kubernetes.io/name: {{ include "install-scripts-controller.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end -}}

{{- define "install-scripts-controller.configSecretName" -}}
{{- default (printf "%s-config" (include "install-scripts-controller.fullname" .)) .Values.secrets.name -}}
{{- end -}}

{{- define "install-scripts-controller.ingressTls" -}}
{{- $tls := list -}}
{{- range .Values.ingress.tls -}}
{{- if .secretName -}}
{{- $tls = append $tls . -}}
{{- end -}}
{{- end -}}
{{- if $tls -}}
{{- toYaml $tls -}}
{{- end -}}
{{- end -}}
