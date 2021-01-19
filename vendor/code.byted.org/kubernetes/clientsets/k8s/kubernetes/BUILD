load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "clientset.go",
        "doc.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes",
    importpath = "code.byted.org/kubernetes/clientsets/k8s/kubernetes",
    visibility = ["//visibility:public"],
    deps = [
        "//staging/src/k8s.io/client-go/discovery:go_default_library",
        "//staging/src/k8s.io/client-go/rest:go_default_library",
        "//staging/src/k8s.io/client-go/util/flowcontrol:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes/typed/csi/v1alpha1:go_default_library",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes/typed/non.native.resource/v1alpha1:go_default_library",
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
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes/scheme:all-srcs",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes/typed/csi/v1alpha1:all-srcs",
        "//vendor/code.byted.org/kubernetes/clientsets/k8s/kubernetes/typed/non.native.resource/v1alpha1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
