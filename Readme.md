# Inventory Worker Service

A Go-based microservice using Hexagonal Architecture to process Kafka events and manage inventory.

## 📋 Prerequisites

Ensure you have the following installed on your machine:
* **Docker** (Desktop or Engine)
* **Docker Compose**
* **Make** (to use the shortcut commands)

---

## 🚀 How to Run

Use the following commands to manage the application lifecycle:

### 1. Start the Application
This will build the services, start the infrastructure (Postgres, Redis, Kafka, Zookeeper), and auto-create the `order.created` topic.
```bash
make run-up