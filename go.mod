module github.com/rh01/terraform-provider-kubeflow-training

go 1.16

require (
	github.com/aws/aws-sdk-go v1.31.9 // indirect
	github.com/hashicorp/go-getter v1.4.2-0.20200106182914-9813cbd4eb02 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20191212124732-c6ae6269b9d7 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/kubevirt/terraform-provider-kubevirt v0.0.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/zclconf/go-cty-yaml v1.0.2 // indirect
	k8s.io/api v0.26.1
	k8s.io/apimachinery v0.26.1
	k8s.io/client-go v12.0.0+incompatible
)

replace (
	github.com/hashicorp/terraform => github.com/openshift/terraform v0.12.20-openshift-4
	github.com/hashicorp/terraform-plugin-sdk => github.com/openshift/hashicorp-terraform-plugin-sdk v1.14.0-openshift
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)
