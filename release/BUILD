load("@rules_terraform//experimental:publishing.bzl", "ghrelease_assets", "ghrelease_publisher", "ghrelease_test_suite")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

VERSION = "0.3"

ARCH_LIST = ["amd64"]

OS_LIST = [
    "darwin",
    "linux",
    "windows",
]

ghrelease_publisher(
    name = "release",
    args = ["--draft"],
    version = VERSION,
    deps = [
        ":prerelease-tests",
    ] + [
        ":%s-%s-plugin" % (os, arch)
        for os in OS_LIST
        for arch in ARCH_LIST
    ],
)

ghrelease_test_suite(
    name = "prerelease-tests",
    tests = [
        "//...",
        "//tests:create_delete_integration_test",
    ],
)

[go_binary(
    name = "terraform-provider-kubectl-%s-%s" % (os, arch),
    embed = ["//:go_default_library"],
    goarch = arch,
    goos = os,
    pure = "on",
    tags = ["manual"],
) for os in OS_LIST for arch in ARCH_LIST]

[ghrelease_assets(
    name = "%s-%s-plugin" % (os, arch),
    tags = ["manual"],
    bazel_flags = [
        "--platforms=@io_bazel_rules_go//go/toolchain:%s_%s" % (os, arch),
    ],
    data = [
        ":terraform-provider-kubectl-%s-%s" % (os, arch),
    ],
) for os in OS_LIST for arch in ARCH_LIST]
