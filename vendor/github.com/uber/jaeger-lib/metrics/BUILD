load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "counter.go",
        "factory.go",
        "gauge.go",
        "histogram.go",
        "keys.go",
        "metrics.go",
        "stopwatch.go",
        "timer.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-lib/metrics",
    importpath = "github.com/uber/jaeger-lib/metrics",
    visibility = ["//visibility:public"],
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
