// This is a generated file. Do not edit directly.

module k8s.io/api

go 1.13

require (
	github.com/gogo/protobuf v1.3.1
	github.com/stretchr/testify v1.7.0
	k8s.io/apimachinery v0.18.10
)

replace (
	github.com/google/cadvisor => code.byted.org/tce/cadvisor v0.0.0-20201220050623-48be6ea97b8c
	github.com/google/go-cmp => github.com/google/go-cmp v0.3.0
	github.com/kr/text => github.com/kr/text v0.1.0
	github.com/opentracing/opentracing-go => code.byted.org/tce/opentracing-go v1.1.0-tce-kube-tracing
	github.com/pkg/errors => github.com/pkg/errors v0.8.1
	github.com/stretchr/testify => github.com/stretchr/testify v1.4.0
	github.com/uber/jaeger-client-go => code.byted.org/tce/jaeger-client-go v2.16.0-tce-kube-tracing+incompatible
	go.uber.org/atomic => go.uber.org/atomic v1.3.2
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // pinned to release-branch.go1.13
	golang.org/x/tools => golang.org/x/tools v0.0.0-20190821162956-65e3620a7ae7 // pinned to release-branch.go1.13
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	k8s.io/api => ../api
	k8s.io/apimachinery => ../apimachinery
	k8s.io/client-go => ../client-go
	k8s.io/code-generator => ../code-generator
	k8s.io/utils => code.byted.org/tce/k8s-utils v0.0.0-20201125131702-a289a73e4c95
)
