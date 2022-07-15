{{/*
Expand the name of the chart.
*/}}
{{- define "simple-bank.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "simple-bank.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "simple-bank.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "simple-bank.labels" -}}
helm.sh/chart: {{ include "simple-bank.chart" . }}
{{ include "simple-bank.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "simple-bank.selectorLabels" -}}
app.kubernetes.io/name: {{ include "simple-bank.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "simple-bank.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "simple-bank.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "simple-bank.db.env" -}}
- name: DATASOURCE_URL
  value: api-postgresql:5432
- name: DATASOURCE_DRIVER
  value: postgres
- name: DATASOURCE_DB
  value: {{ .Values.postgresql.postgresqlDatabase }}
- name: DATASOURCE_USER
  value: {{ .Values.postgresql.postgresqlUsername }}
{{- if .Values.postgresql.postgresqlPassword }}
- name: DATASOURCE_PASSWORD
  value: {{ .Values.postgresql.postgresqlPassword }}
{{- else }}
- name: DATASOURCE_PASSWORD
  valueFrom:
    secretKeyRef:
      name: api-postgresql
      key: postgressql-password
{{- end }}
- name: TOKEN_SYMMETRIC_KEY
  value: asd1231fsasf1231asds123asd123
- name: ACCESS_TOKEN_DURATION
  value: 15m
{{- end }}
