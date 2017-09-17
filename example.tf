
resource kubectl_generic_object test_deployment {
  // language=yaml
  yaml = <<EOF
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: tf-test
  namespace: integration-testing
  annotations:
    prometheus.io/scrape: "true"
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: tf-test
    spec:
      containers:
      - name: server
        image: busybox:latest
        args: [ "while true; do sleep 30; done" ]
        resources:
          requests: {cpu: 50m, memory: 50m}
EOF
}
