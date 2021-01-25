load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["json.go"],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-client-go/internal/spanlog",
    importpath = "github.com/uber/jaeger-client-go/internal/spanlog",
    visibility = ["//vendor/github.com/uber/jaeger-client-go:__subpackages__"],
    deps = ["//vendor/github.com/opentracing/opentracing-go/log:go_default_library"],
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
