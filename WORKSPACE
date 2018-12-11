load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "rules_terraform",
    commit = "87b181f15378d5034aed38c43274ab01b808f27f",
    remote = "git@github.com:ceason/rules_terraform.git",
)

git_repository(
    name = "io_bazel_rules_docker",
    commit = "e5785ceaef4eb7e0cc28bdb909fd1b10d5b991c3",
    remote = "git@github.com:bazelbuild/rules_docker.git",
)

git_repository(
    name = "io_bazel_rules_k8s",
    commit = "d6e1b65317246fe044482f9e042556c77e6893b8",
    remote = "git@github.com:bazelbuild/rules_k8s.git",
)

git_repository(
    name = "io_bazel_rules_go",
    commit = "e56822c37c2f3d4e6aff7937b570e9db9ab753ff",
    remote = "https://github.com/bazelbuild/rules_go.git",
)

git_repository(
    name = "bazel_gazelle",
    commit = "44ce230b3399a5d4472198740358fcd825b0c3c9",
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

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

load("@io_bazel_rules_docker//container:container.bzl", "repositories")

repositories()

load("@rules_terraform//terraform:dependencies.bzl", "terraform_repositories")

terraform_repositories()

go_repository(
    name = "com_github_golang_dep",
    importpath = "github.com/golang/dep",
    tag = "v0.5.0",
)
