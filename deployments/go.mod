module github.com/kubearmor/KubeArmor/deployments

go 1.18

replace (
	github.com/kubearmor/KubeArmor => ../
	github.com/kubearmor/KubeArmor/pkg/KubeArmorController => ../pkg/KubeArmorController
	github.com/kubearmor/KubeArmor/pkg/KubeArmorHostPolicy => ../pkg/KubeArmorHostPolicy
	github.com/kubearmor/KubeArmor/pkg/KubeArmorPolicy => ../pkg/KubeArmorPolicy
	k8s.io/api => k8s.io/api v0.22.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.3
)

require (
	github.com/clarketm/json v1.17.1
	github.com/kubearmor/KubeArmor/pkg/KubeArmorController v0.0.0-00010101000000-000000000000
	k8s.io/api v0.22.3
	k8s.io/apimachinery v0.22.3
	sigs.k8s.io/yaml v1.2.0
)

require (
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.22.3 // indirect
	k8s.io/klog/v2 v2.10.0 // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
)
