load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "constants.go",
        "ttypes.go",
        "zipkincollector.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-client-go/thrift-gen/zipkincore",
    importpath = "github.com/uber/jaeger-client-go/thrift-gen/zipkincore",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/uber/jaeger-client-go/thrift:go_default_library"],
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
