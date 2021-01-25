load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "context.go",
        "grpc.go",
        "logger.go",
        "store.go",
        "tracer.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/code.byted.org/tce/kube-tracing",
    importpath = "code.byted.org/tce/kube-tracing",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/opentracing/opentracing-go:go_default_library",
        "//vendor/github.com/opentracing/opentracing-go/log:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go:go_default_library",
        "//vendor/github.com/uber/jaeger-client-go/config:go_default_library",
        "//vendor/google.golang.org/grpc:go_default_library",
        "//vendor/google.golang.org/grpc/metadata:go_default_library",
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
