load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

exports_files(glob(["*.patch"]))

# gazelle:prefix github.com/ceason/terraform-provider-kubectl

# 'dep ensure' updates the vendor/ directory, which means we need to
# re-run gazelle. This script takes care of that.
genrule(
    name = "dep-ensure",
    outs = ["dep-ensure.sh"],
    cmd = """cat <<'EOF' > $@
#!/usr/bin/env bash
set -euo pipefail
cd "$$BUILD_WORKSPACE_DIRECTORY"
$(location @com_github_golang_dep//cmd/dep) ensure
$(location @bazel_gazelle//cmd/gazelle) update
git add Gopkg.{lock,toml} vendor/. $$(find . -name BUILD -o -name BUILD.bazel)
EOF""",
    executable = True,
    tools = [
        "@bazel_gazelle//cmd/gazelle",
        "@com_github_golang_dep//cmd/dep",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "data_source_kubectl_namespace.go",
        "kubectl_cli.go",
        "main.go",
        "provider.go",
        "resource_kubectl_generic_object.go",
        "types.go",
    ],
    importpath = "github.com/ceason/terraform-provider-kubectl",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/hashicorp/terraform/helper/schema:go_default_library",
        "//vendor/github.com/hashicorp/terraform/plugin:go_default_library",
        "//vendor/github.com/hashicorp/terraform/terraform:go_default_library",
        "//vendor/gopkg.in/yaml.v2:go_default_library",
        "@io_k8s_kubernetes//pkg/kubectl/cmd:go_default_library",
    ],
)

go_binary(
    name = "terraform-provider-kubectl",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    size = "small",
    srcs = ["resource_kubectl_generic_object_test.go"],
    data = glob(["testdata/**"]),
    embed = [":go_default_library"],
    deps = ["//vendor/gopkg.in/yaml.v2:go_default_library"],
)
