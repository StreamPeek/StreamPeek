# StreamPeek

![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
![Release](https://img.shields.io/github/v/release/StreamPeek/StreamPeek)

**High-performance, lightweight REST API for Apache Kafka.** Produce, consume, and inspect message streams over simple HTTP without the JVM overhead.

Interacting with Apache Kafka shouldn't require wrestling with complex native drivers or setting up a heavy Java ecosystem just to test a message. **StreamPeek** is a zero-dependency REST gateway built specifically for developers, performance engineers, and QA teams. It acts as a high-speed bridge, allowing you to interact with Kafka using standard HTTP requests. 

Whether you are injecting test data from Postman, running high-throughput performance benchmarks, or asserting message delivery in CI/CD pipelines, StreamPeek stays out of your way and lets you work fast.

## ✨ Why StreamPeek?

* **⚡ High Performance:** Engineered for low latency and high throughput. Ideal for performance testing and stress-testing Kafka clusters without "JVM warmup" lag.
* **🪶 No JVM Required:** Shipped as a single, lightweight binary. No heavy enterprise proxies or memory-hungry runtimes to configure.
* **🧪 Built for Testing:** Designed specifically to make local development, QA environments, and automated integration testing painless.
* **🌐 Universal Compatibility:** Talk to Kafka using `curl`, Python `requests`, Postman, or any language that speaks HTTP.
* **📦 Stateless & Cloud-Native:** Instant startup times make it perfect for ephemeral Docker containers and scaling horizontally in Kubernetes.

## 🏎️ Performance & Testing focus
StreamPeek is optimized for scenarios where speed and reliability matter:
* **Benchmark your Cluster:** Use standard HTTP load-testing tools (like `wrk` or `k6`) to stress-test your brokers via REST.
* **CI/CD Integration:** Spin up a lightweight sidecar in your pipeline to verify message flow with a minimal resource footprint.
* **Rapid Prototyping:** Skip the boilerplate of Kafka producers and start sending thousands of events per second via a simple REST interface.

## 🚀 Quick Start

### 1. Run StreamPeek
The easiest way to get started is using Docker. Point StreamPeek to your existing Kafka broker:

```bash
docker run -p 8080:8080 \
  -e KAFKA_BROKERS="localhost:9092" \
  streampeek/streampeek:latest
