load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "csi_client.go",
        "csidriver.go",
        "csinodeinfo.go",
        "csistoragepool.go",
        "doc.go",
        "generated_expansion.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes/typed/csi/v1alpha1",
    importpath = "code.byted.org/kubernetes/clientsets/k8s/kubernetes/typed/csi/v1alpha1",
    visibility = ["//visibility:public"],
    deps = [
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/types:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//staging/src/k8s.io/client-go/rest:go_default_library",
        "//vendor/code.byted.org/kubernetes/apis/k8s/csi/v1alpha1:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes/scheme:go_default_library",
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
