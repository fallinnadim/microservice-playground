# Inventory Worker Service

A Go-based microservice using Hexagonal Architecture to process Kafka events and manage inventory.

## 📋 Prerequisites

Ensure you have the following installed on your machine:
* **Docker** (Desktop or Engine)
* **Docker Compose**
* **Make** (to use the shortcut commands)

---

## 🚀 How to Run

For ease of use, all necessary configuration files, including **`.env`** and **SSL certificates**, are already committed to this repository. You do not need to create them manually.

### 1. Start the Application
This command builds the Go services, starts the infrastructure (Postgres, Redis, Kafka, Zookeeper), and automatically initializes the `order.created` topic with 10 partitions.
```bash
make run-up