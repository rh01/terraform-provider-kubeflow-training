 
provider "terraform-provider-kubeflow-training" {
  # Configuration options...
}

provider "kubeflow" {
  # Configuration options...
}
provider "kubernetes" {
  host = var.host

  client_certificate     = base64decode(var.client_certificate)
  client_key             = base64decode(var.client_key)
  cluster_ca_certificate = base64decode(var.cluster_ca_certificate)
}
terraform {
  required_providers {
    terraform-provider-kubeflow-training = {
      source  = "example.com/exampleprovider/kubeflow-training"
      version = "1.0.0"
      # Other parameters...
    }
    kubeflow = {
      source  = "example.com/exampleprovider/kubeflow-training"
      version = "1.0.0"
      # Other parameters...
    }
  }
}

resource "kubeflow_pytorch_job" "pytorch_job" {

  metadata {
    name = "pytorch-job"
    namespace = "kubeflow"
  }

  spec {
    
    pytorch_replica_specs {
      master {
        replicas = 1
        template {
          spec {
            container {
              args = [ "python", "train.py" ]
              command = [ "python", "train.py" ]
              image = "gcr.io/kubeflow-examples/mnist"
              name = "pytorch"
            }
          }
        }
        restart_policy = "Never"
      }

      worker {
        
      }

    }

    
    
  }
}


