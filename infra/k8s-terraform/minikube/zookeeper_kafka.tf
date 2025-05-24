# Namespace global
resource "kubernetes_namespace" "minikafka" {
  metadata {
    name = "minikafka"
  }
}

# Deployment de Zookeeper-Service
resource "kubernetes_deployment" "zookeeper" {
  metadata {
    name      = "zookeeper-deploy"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
    labels = {
      app = "zookeeper"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "zookeeper"
      }
    }

    template {
      metadata {
        labels = {
          app = "zookeeper"
        }
      }

      spec {
        container {
          name  = "zookeeper"
          image = "confluentinc/cp-zookeeper:latest"

          env {
            name  = "ZOOKEEPER_CLIENT_PORT"
            value = "2181"
          }

          port {
            container_port = 2181
          }
        }
      }
    }
  }
}

# Service de Zookeeper-Service
resource "kubernetes_service" "zookeeper" {
  metadata {
    name      = "zookeeper-service"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
  }

  spec {
    selector = {
      app = "zookeeper"
    }

    port {
      name        = "client"
      port        = 2181
      target_port = 2181
    }
  }
}

# Deployment de kafka-Service
resource "kubernetes_deployment" "kafka" {
  metadata {
    name      = "kafka-deploy"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
    labels = {
      app = "kafka"
    }
  }

  depends_on = [kubernetes_deployment.zookeeper, kubernetes_service.zookeeper]

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "kafka"
      }
    }

    template {
      metadata {
        labels = {
          app = "kafka"
        }
      }

      spec {
        init_container {
          name  = "wait-for-zookeeper"
          image = "busybox:1.28"
          command = [
            "sh", "-c",
            "until nc -z zookeeper-service.minikafka.svc.cluster.local 2181; do echo waiting for zookeeper; sleep 2; done;"
          ]
        }

        container {
          name  = "kafka"
          image = "confluentinc/cp-kafka:latest"

          env {
            name  = "KAFKA_BROKER_ID"
            value = "1"
          }

          env {
            name  = "KAFKA_ZOOKEEPER_CONNECT"
            value = "zookeeper-service:2181"
          }

          env {
            name  = "KAFKA_LISTENERS"
            value = "PLAINTEXT://:9092"
          }

          env {
            name  = "KAFKA_ADVERTISED_LISTENERS"
            value = "PLAINTEXT://kafka-service:9092"
          }

          env {
            name  = "KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR"
            value = "1"
          }

          env {
            name  = "KAFKA_TRANSACTION_STATE_LOG_MIN_ISR"
            value = "1"
          }

          env {
            name  = "KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR"
            value = "1"
          }


          port {
            container_port = 9092
          }
        }
      }
    }
  }
}

# Service de kafka-Service
resource "kubernetes_service" "kafka-service" {
  metadata {
    name      = "kafka-service"
    namespace = kubernetes_namespace.minikafka.metadata[0].name
  }

  spec {
    selector = {
      app = "kafka"
    }

    port {
      name        = "broker"
      port        = 9092
      target_port = 9092
    }
  }
}
