# Deployment de Mongo-User
resource "kubernetes_deployment" "mongo_user" {
  metadata {
    name      = "mongo-user-deploy"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
    labels = {
      app = "mongo-user"
    }
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "mongo-user"
      }
    }

    template {
      metadata {
        labels = {
          app = "mongo-user"
        }
      }

      spec {
        container {
          name  = "mongo-user"
          image = "mongo:latest"
          port {
            container_port = 27017
          }
          command = ["mongod", "--quiet", "--logpath", "/dev/null"]
          volume_mount {
            mount_path = "/data/db"
            name       = "mongo-user-storage"
          }
        }

        volume {
          name = "mongo-user-storage"

          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.mongo_user.metadata[0].name
          }
        }
      }
    }
  }
}

# PVC de Mongo-User
resource "kubernetes_persistent_volume_claim" "mongo_user" {
  metadata {
    name      = "mongo-user-pvc"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
  }

  spec {
    access_modes = ["ReadWriteOnce"]

    resources {
      requests = {
        storage = "256Mi"
      }
    }
  }
}

# Service de Mongo-User
resource "kubernetes_service" "mongo_user" {
  metadata {
    name      = "mongo-user"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
  }

  spec {
    selector = {
      app = "mongo-user"
    }

    port {
      name        = "mongo"
      port        = 27017
      target_port = 27017
    }
  }
}

# ConfigMap de User-Service
resource "kubernetes_config_map" "user_env" {
  metadata {
    name      = "user-env"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
  }

  data = {
    MONGO_URI     = "mongodb://mongo-user:27017"
    MONGO_DB_NAME = "trello_clone"
    JWT_SECRET    = "supersecretkey"
    REDIS_ADDR    = "redis:6379"
    REDIS_PASS    = ""
  }
}

# Deployment de User-Service
resource "kubernetes_deployment" "user_service" {
  metadata {
    name      = "user-service"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
    labels = {
      app = "user-service"
    }
  }

  depends_on = [kubernetes_deployment.mongo_user,
  kubernetes_deployment.kafka]

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "user-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "user-service"
        }
      }

      spec {
        container {
          name              = "user-service"
          image             = "user-service:local"
          image_pull_policy = "Never" # imagePullPolicy: Never es construccion local
          port {
            container_port = 8080
          }
          env {
            name = "MONGO_URI"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_env.metadata[0].name
                key  = "MONGO_URI"
              }
            }
          }

          env {
            name = "MONGO_DB_NAME"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_env.metadata[0].name
                key  = "MONGO_DB_NAME"
              }
            }
          }

          env {
            name = "JWT_SECRET"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_env.metadata[0].name
                key  = "JWT_SECRET"
              }
            }
          }

          env {
            name = "REDIS_ADDR"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_env.metadata[0].name
                key  = "REDIS_ADDR"
              }
            }
          }

          env {
            name = "REDIS_PASS"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_env.metadata[0].name
                key  = "REDIS_PASS"
              }
            }
          }


        }
      }
    }
  }
}

# Service de User-Service
resource "kubernetes_service" "user_service" {
  metadata {
    name      = "user-service"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
  }

  spec {
    selector = {
      app = "user-service"
    }

    port {
      name        = "http"
      port        = 8080
      target_port = 8080
    }

    type = "NodePort"
  }
}
