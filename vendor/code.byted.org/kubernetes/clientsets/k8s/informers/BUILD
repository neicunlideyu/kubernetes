load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "factory.go",
        "generic.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/code.byted.org/kubernetes/clientsets/k8s/informers",
    importpath = "code.byted.org/kubernetes/clientsets/k8s/informers",
    visibility = ["//visibility:public"],
    deps = [
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//staging/src/k8s.io/client-go/tools/cache:go_default_library",
        "//vendor/code.byted.org/kubernetes/apis/k8s/csi/v1alpha1:go_default_library",
        "//vendor/code.byted.org/kubernetes/apis/k8s/non.native.resource/v1alpha1:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/csi:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/internalinterfaces:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/non.native.resource:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes:go_default_library",
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
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/csi:all-srcs",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/internalinterfaces:all-srcs",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/informers/non.native.resource:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
