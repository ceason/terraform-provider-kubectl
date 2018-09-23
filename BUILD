load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")
load("@rules_terraform//terraform:def.bzl", "terraform_plugin")

gazelle(
    name = "gazelle",
    prefix = "github.com/ceason/terraform-provider-kubectl",
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
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/hashicorp/terraform/helper/schema:go_default_library",
        "//vendor/github.com/hashicorp/terraform/plugin:go_default_library",
        "//vendor/github.com/hashicorp/terraform/terraform:go_default_library",
        "//vendor/gopkg.in/yaml.v2:go_default_library",
    ],
)

go_binary(
    name = "terraform-provider-kubectl",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["kubectl_cli_test.go"],
    data = glob(["testdata/**"]),
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)

go_binary(
    name = "linux_amd64",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
)

go_binary(
    name = "darwin_amd64",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "darwin",
    pure = "on",
)

terraform_plugin(
    name = "plugin",
    darwin_amd64 = ":darwin_amd64",
    linux_amd64 = ":linux_amd64",
    provider_name = "kubectl",
    version = "v0.3.0",
    visibility = ["//visibility:public"],
)
