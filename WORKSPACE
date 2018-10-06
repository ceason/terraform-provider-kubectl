load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "bazel_project_helpers",
    commit = "768dc16dea23ceccd84c2dfa5ad89738433c7003",
    remote = "git@rvgithub.com:redventures/bazel_project_helpers.git",
)

git_repository(
    name = "rules_terraform",
    commit = "505e93a4ac1ca2b6f0757c8c2065278a559355ea",
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
    commit = "40e2b78a314ebb91d0e690579ed3273683a3a1a1",
    remote = "https://github.com/bazelbuild/rules_go.git",
)

git_repository(
    name = "bazel_gazelle",
    tag = "0.14.0",
    remote = "https://github.com/bazelbuild/bazel-gazelle.git",
)

git_repository(
    name = "io_k8s_kubernetes",
    patches = ["@//:io_k8s_kubernetes.patch"],
    remote = "https://github.com/kubernetes/kubernetes.git",
    tag = "v1.12.0",
)

git_repository(
    name = "io_kubernetes_build",
    commit = "84d52408a061e87d45aebf5a0867246bdf66d180",
    remote = "https://github.com/kubernetes/repo-infra.git",
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
