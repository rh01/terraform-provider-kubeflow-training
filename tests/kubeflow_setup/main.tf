provider "kubernetes" {}

provider "k8s" {}

provider "helm" {}

module "kubeflow" {
  providers = {
    kubernetes = kubernetes
    k8s        = k8s
    helm       = helm
  }

  source  = "datarootsio/kubeflow/module"
  version = "~>0.0"

  ingress_gateway_ip  = "10.20.30.40"
  use_cert_manager    = true
  domain_name         = "foo.local"
  letsencrypt_email   = "foo@bar.local"
  kubeflow_components = ["pipelines"]
}