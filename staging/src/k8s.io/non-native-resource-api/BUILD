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
        "//staging/src/k8s.io/non-native-resource-api/pkg/apis/non.native.resource/v1alpha1:all-srcs",
        "//staging/src/k8s.io/non-native-resource-api/pkg/client/clientset/versioned:all-srcs",
        "//staging/src/k8s.io/non-native-resource-api/pkg/client/informers/externalversions:all-srcs",
        "//staging/src/k8s.io/non-native-resource-api/pkg/client/listers/non.native.resource/v1alpha1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
