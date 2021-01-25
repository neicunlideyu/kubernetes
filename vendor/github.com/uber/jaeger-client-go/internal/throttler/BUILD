load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["throttler.go"],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-client-go/internal/throttler",
    importpath = "github.com/uber/jaeger-client-go/internal/throttler",
    visibility = ["//vendor/github.com/uber/jaeger-client-go:__subpackages__"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//vendor/github.com/uber/jaeger-client-go/internal/throttler/remote:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
