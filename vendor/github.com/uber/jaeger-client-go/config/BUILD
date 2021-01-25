load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "config_env.go",
        "options.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-client-go/config",
    importpath = "github.com/uber/jaeger-client-go/config",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/opentracing/opentracing-go:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go/internal/baggage/remote:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go/internal/throttler/remote:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go/rpcmetrics:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go/transport:go_default_library",
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
