load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "bazel_project_helpers",
    commit = "768dc16dea23ceccd84c2dfa5ad89738433c7003",
    remote = "git@rvgithub.com:redventures/bazel_project_helpers.git",
)

git_repository(
    name = "rules_terraform",
    commit = "295e923397dc0dc7f0916e5c7f38c98c06af6e8e",
    remote = "git@github.com:ceason/rules_terraform.git",
)

git_repository(
    name = "io_bazel_rules_docker",
    commit = "7401cb256222615c497c0dee5a4de5724a4f4cc7",
    remote = "git@github.com:bazelbuild/rules_docker.git",
)

git_repository(
    name = "io_bazel_rules_k8s",
    commit = "d6e1b65317246fe044482f9e042556c77e6893b8",
    remote = "git@github.com:bazelbuild/rules_k8s.git",
)

git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    commit = "f8c9f2c6336536147458aaccbd1becf5cc80232a",
)


git_repository(
    name = "bazel_gazelle",
    remote = "https://github.com/bazelbuild/bazel-gazelle.git",
    commit = "993d887662ad83bb60b9ba1570270d4afcda91a1",
)


load("@io_bazel_rules_docker//container:container.bzl", "repositories")

repositories()

load("@io_bazel_rules_k8s//k8s:k8s.bzl", "k8s_defaults", "k8s_repositories")

k8s_repositories()

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@rules_terraform//terraform:dependencies.bzl", "terraform_repositories")

terraform_repositories()

load("@bazel_project_helpers//:dependencies.bzl", "bazel_project_helpers_repositories")

bazel_project_helpers_repositories()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()