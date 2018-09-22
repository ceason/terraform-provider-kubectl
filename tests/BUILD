load("@rules_terraform//terraform:def.bzl", "terraform_distribution_dir", "terraform_integration_test", "terraform_module", "terraform_workspace")

terraform_workspace(
    name = "workspace",
    srcs = ["create-delete.tf"],
    plugins = ["//:plugin"],
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
    terraform_workspace = ":workspace",
    tags = ["manual"],
)
