resource random_string uniqifier {
  length  = 5
  special = false
  upper   = false
}

data kubectl_namespace current {}

provider kubectl {
  namespace = "kube-public"
}

resource kubernetes_config_map test_configmap {
  metadata {
    name      = "createdelete-test-configmap-${random_string.uniqifier.result}"
    namespace = "${data.kubectl_namespace.current.id}"
  }
  data {
    TEST_ASDF.X = "1234"
    KEY_NUM2    = "asdfX"
  }
}

resource local_file vars {
  filename = "test_vars.sh"
  content  = <<EOF
TEST_SVCACCT=${local.test_svcacct}
TEST_CONFIGMAP=createdelete-test-configmap-${random_string.uniqifier.result}
TEST_NAMESPACE=${data.kubectl_namespace.current.id}
EOF
}


variable "asdf" {
  default = "something"
}

locals {
  test_svcacct="test-create-delete-${var.asdf}"
}


resource kubectl_generic_object test_svcacct {
  #// language=yaml
  yaml = <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ${local.test_svcacct}
imagePullSecrets:
- name: some-test-name
EOF
}