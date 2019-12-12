---
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      networkservicemesh.io/app: "vpn-gateway-nsc"
      networkservicemesh.io/impl: "secure-intranet-connectivity"
  replicas: 2
  template:
    metadata:
      labels:
        networkservicemesh.io/app: "vpn-gateway-nsc"
        networkservicemesh.io/impl: "secure-intranet-connectivity"
        app: iperf-vpn-client
    spec:
      serviceAccount: nsc-acc
      containers:
        - name: alpine-img
          image: alpine:latest
          imagePullPolicy: {{ .Values.pullPolicy }}
          command: ['tail', '-f', '/dev/null']
        - name: iperf-vpn-client
          image: jmarhee/iperf:latest
          #command: ["iperf", "-c", "iperf-vpn-server", "-P", "10"]
          #command: ["iperf", "-c", "iperf-vpn-server", "-u", "-b", "100m", "-t", "40"]
          command: ["iperf", "-c", "iperf-vpn-server", "-u", "-b", "100m", "-t", "20", "-i", "2"]
          #command: ["iperf", "-c", "iperf-vpn-server", "-u"]

        - name: iperf-vpn-firewall-client
          image: jmarhee/iperf:latest
          command: ["iperf", "-c", "iperf-firewall-server", "-P", "100"]
          #command: ["iperf", "-c", "iperf-firewall-server", "-u"]


metadata:
  name: vpn-gateway-nsc
  namespace: {{ .Release.Namespace }}
  annotations:
    ns.networkservicemesh.io: secure-intranet-connectivity
