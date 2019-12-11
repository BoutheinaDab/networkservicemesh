---
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      networkservicemesh.io/app: "vpn-gateway"
      networkservicemesh.io/impl: "secure-intranet-connectivity"
      app: iperf-vpn-server
  replicas: 2
  template:
    metadata:
      labels:
        networkservicemesh.io/app: "vpn-gateway"
        networkservicemesh.io/impl: "secure-intranet-connectivity"
        app: iperf-vpn-server
    spec:
      serviceAccount: nse-acc
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: networkservicemesh.io/app
                    operator: In
                    values:
                      - vpn-gateway
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: vpn-gateway
          image: {{ .Values.registry }}/{{ .Values.org }}/test-common:{{ .Values.tag }}
          command: ["/bin/icmp-responder-nse"]
          imagePullPolicy: {{ .Values.pullPolicy }}
          env:
            - name: ADVERTISE_NSE_NAME
              value: "secure-intranet-connectivity"
            - name: ADVERTISE_NSE_LABELS
              value: "app=vpn-gateway"
            - name: IP_ADDRESS
              value: "172.16.1.0/24"
{{- if .Values.global.JaegerTracing }}
            - name: TRACER_ENABLED
              value: "true"
            - name: JAEGER_AGENT_HOST
              value: jaeger.nsm-system
            - name: JAEGER_AGENT_PORT
              value: "6831"
{{- end }}
          resources:
            limits:
              networkservicemesh.io/socket: 1
        - name: nginx
          image: {{ .Values.registry }}/networkservicemesh/nginx:latest
        - name: iperf-vpn-server
          image: jmarhee/iperf:latest
          #command: ['iperf', '-s', '-u', '-i', '5']
          command: ['iperf', '-s', '-u']
          #command: ['iperf', '-s']
          ports:
          - containerPort: 5001

metadata:
  name: vpn-gateway-nse
  namespace: {{ .Release.Namespace }}
  labels:
    app: iperf-vpn-server

---
kind: Service
apiVersion: v1
metadata:
  name: iperf-vpn-server
  namespace: nsm-system
spec:
  selector:
    app: iperf-vpn-server
  ports:
  - protocol: UDP
    port: 5001
    targetPort: 5001

