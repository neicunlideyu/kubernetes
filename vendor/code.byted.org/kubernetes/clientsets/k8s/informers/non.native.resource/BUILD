load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["interface.go"],
    importmap = "k8s.io/kubernetes/vendor/code.byted.org/kubernetes/clientsets/k8s/informers/non.native.resource",
    importpath = "code.byted.org/kubernetes/clientsets/k8s/informers/non.native.resource",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/internalinterfaces:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/non.native.resource/v1alpha1:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/non.native.resource/v1alpha1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
