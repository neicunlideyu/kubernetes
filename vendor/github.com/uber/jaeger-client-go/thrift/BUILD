load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "application_exception.go",
        "binary_protocol.go",
        "compact_protocol.go",
        "exception.go",
        "memory_buffer.go",
        "messagetype.go",
        "numeric.go",
        "processor.go",
        "protocol.go",
        "protocol_exception.go",
        "protocol_factory.go",
        "rich_transport.go",
        "serializer.go",
        "simple_json_protocol.go",
        "transport.go",
        "transport_exception.go",
        "transport_factory.go",
        "type.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/uber/jaeger-client-go/thrift",
    importpath = "github.com/uber/jaeger-client-go/thrift",
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
