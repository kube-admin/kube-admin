{{/*
公共标签
*/}}
{{- define "kube-admin.labels" -}}
app.kubernetes.io/name: kube-admin
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
backend Service 名称固定为 backend：前端 nginx 反代 http://backend:8080 依赖此 DNS 名。
*/}}
{{- define "kube-admin.backendServiceName" -}}
backend
{{- end -}}

{{/*
ServiceAccount 名称
*/}}
{{- define "kube-admin.serviceAccountName" -}}
{{- if .Values.serviceAccount.name -}}
{{- .Values.serviceAccount.name -}}
{{- else -}}
kube-admin-sa
{{- end -}}
{{- end -}}
