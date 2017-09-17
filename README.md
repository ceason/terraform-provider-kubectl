## Motivation
> todo: write this

## Usage
This provider has a few prerequsites right now:
- `kubectl` must exist on your $PATH (may change in future)
- Download the plugin from the [Releases](https://github.com/ceason/terraform-provider-kubectl/releases) page
- [Install](https://terraform.io/docs/plugins/basics.html) it

Example "generic object" resource
```hcl
provider kubectl {
  context = "kubeconfig-context-name"
}

resource kubectl_generic_object example_deployment {
  yaml = <<EOF
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: tf-test
  namespace: sdlc-demo-test
  annotations:
    blarghtest: "something"
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
          requests: {cpu: 100m, memory: 50m}
EOF
}
```



## Design notes
> todo: make this section make more sense

#### Logging
This provider currently uses `glog`. To configure it, you must pass flags via environment variable like this:
```shell
export TF_PROVIDER_KUBECTL_FLAGS="-v=4 -log_dir=/tmp/glog"
```

#### Flow of execution
- `create`
    - generate and set `id` globally unique id
    - call `apply`
- `apply`
    - merge `metadata.labels.$PRUNE_LABEL=$id` with `yaml`
    - execute command `kubectl --context=$CONTEXT apply --prune -l $PRUNE_LABEL=$id -f -` (reads from stdin)
- `delete`
    - execute command `kubectl --context=$CONTEXT delete -f -` (reads from stdin)
- `read`
    - noop for now
    - pretty sure this is unnecessary for min functionality because `apply` handles all of the state transition logic for us

## Todo
- More tests
    - test kubectl wrapper via object lifecycle (create, update, delete)
- Input validation
    - Ensure `namespace` specified for all candidate objects (probably by whitelisting server-level stuff)
    - Ensure single object passed in to yaml
    - Ensure specified kubectl context exists
- More flexible provider config
    - Perhaps via environment variables, etc. The official `kubernetes` provider seems nice
    - Maybe download `kubectl` if not present. Could even specify version
