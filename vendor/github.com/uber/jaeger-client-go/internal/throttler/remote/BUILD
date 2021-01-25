load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "options.go",
        "throttler.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-client-go/internal/throttler/remote",
    importpath = "github.com/uber/jaeger-client-go/internal/throttler/remote",
    visibility = ["//vendor/github.com/uber/jaeger-client-go:__subpackages__"],
    deps = [
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go/utils:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
