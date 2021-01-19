load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "node_package_cpu_match.go",
        "node_package_match.go",
        "node_package_mem_match.go",
        "node_package_network_match.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/scheduler/framework/plugins/nodepackage",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api/v1/resource:go_default_library",
        "//pkg/features:go_default_library",
        "//pkg/scheduler/framework/v1alpha1:go_default_library",
        "//pkg/scheduler/nodeinfo:go_default_library",
        "//pkg/scheduler/util:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/util/feature:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/listers/non.native.resource/v1alpha1:go_default_library",
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
