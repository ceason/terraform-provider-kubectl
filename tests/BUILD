load("@rules_terraform//terraform:def.bzl", "terraform_integration_test", "terraform_workspace")
load("@rules_terraform//terraform:def.bzl", "terraform_provider")

terraform_provider(
    name = "provider",
    file = "//:terraform-provider-kubectl",
    visibility = ["//visibility:public"],
)

terraform_workspace(
    name = "workspace",
    srcs = ["create-delete.tf"],
    plugins = [":provider"],
)

sh_test(
    name = "create_delete_test",
    size = "small",
    srcs = ["create-delete.sh"],
    data = [":workspace"],
    tags = ["manual"],
)

terraform_integration_test(
    name = "create_delete_integration_test",
    size = "small",
    srctest = ":create_delete_test",
    tags = [
        # "manual",
    ],
    terraform_workspace = ":workspace",
)
