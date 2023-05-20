 
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
    namespace = "ai-training"
  }

  

  spec {
    
    pytorch_replica_specs {
      master {
        replicas = 1
        restart_policy = "Never"

        template {
          spec {
            container {
              args = [ "python", "train.py" ]
              command = [ "python", "train.py" ]
              image = "gcr.io/kubeflow-examples/mnist"
              name = "pytorch"
              env {
                name = "NCCL_DEBUG"
                value = "INFO"
              }
              env {
                name = "NCCL_IB_DISABLE"
                value = "0"
              }

              env {
                name = "MASTER"
                value = "1"
              }
              image_pull_policy = "IfNotPresent"
              resources {
               limits = {
                "nvidia.com/gpu" = 8
                "rdma/hca"       = 1
               }

              }
              security_context {
                capabilities {
                  add = [ "IPC_LOCK" ]
                }
                 
              }
              volume_mount {
                 mount_path = "/dev/shm"
                      name      = "cache-volume"
              }

              volume_mount {
                 mount_path = "/mnt/pfs"
                      name      = "data"
              }
            }
            image_pull_secrets {
              name = "regcred"
            }
            volume {
              name = "cache-volume"
              empty_dir {
                medium = "Memory"
              }
            }
            volume {
              name = "data"
              persistent_volume_claim {
                claim_name = "pfs-pvc-model"
              }
            }
            scheduler_name = "volcano"
          }
        }
      }

      worker {
        replicas = 1
        restart_policy = "Never"

        template {
          spec {
            container {
              args = [ "python", "train.py" ]
              command = [ "python", "train.py" ]
              image = "gcr.io/kubeflow-examples/mnist"
              name = "pytorch"
              env {
                name = "NCCL_DEBUG"
                value = "INFO"
              }
              env {
                name = "NCCL_IB_DISABLE"
                value = "0"
              }

        
              image_pull_policy = "IfNotPresent"
              resources {
               limits = {
                "nvidia.com/gpu" = 8
                "rdma/hca"       = 1
               }

              }
              security_context {
                capabilities {
                  add = [ "IPC_LOCK" ]
                }
                 
              }
              volume_mount {
                 mount_path = "/dev/shm"
                      name      = "cache-volume"
              }

              volume_mount {
                 mount_path = "/mnt/pfs"
                      name      = "data"
              }
            }
            image_pull_secrets {
              name = "regcred"
            }
            volume {
              name = "cache-volume"
              empty_dir {
                medium = "Memory"
              }
            }
            volume {
              name = "data"
              persistent_volume_claim {
                claim_name = "pfs-pvc-model"
              }
            }
            scheduler_name = "volcano"
          }
        }
      }

    }
    
  }
}


