resource random_string uniqifier {
  length  = 5
  special = false
  upper   = false
}

provider kubectl {

  namespace = "asdfblah"
}

provider kubernetes {
  cluster_ca_certificate = ""
  host                   = ""
  token                  = ""
  client_certificate     = ""
  client_key             = ""
}

resource kubernetes_config_map test_configmap {
  metadata {
    name = "createdelete-test-configmap-${random_string.uniqifier.result}"
  }
  data {
    TEST_ASDF = "1234"
  }
}


resource kubectl_generic_object test_svcacct {
  // language=yaml
  yaml = <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-create-delete

EOF
}