load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "endpoints.go",
        "metrics.go",
        "normalizer.go",
        "observer.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-client-go/rpcmetrics",
    importpath = "github.com/uber/jaeger-client-go/rpcmetrics",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/opentracing/opentracing-go:go_default_library",
        "//vendor/github.com/opentracing/opentracing-go/ext:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go:go_default_library",
        "//vendor/github.com/uber/jaeger-lib/metrics:go_default_library",
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
