# FluxGuard: Distributed Rate Limiter

FluxGuard is a high-performance, distributed rate-limiting service built in **Go**, designed to protect microservices from traffic spikes and abuse. It implements the **Token Bucket algorithm** using **Redis Lua scripts** for atomic state management and exposes a low-latency **gRPC** interface.

## Tech Stack
- **Language:** Go (Golang) 1.22+
- **Communication:** gRPC (Protobuf)
- **Database:** Redis (with Lua scripting for atomicity)
- **Architecture:** Distributed System, Fail-open design

## Key Features
- **Atomic Token Bucket:** Uses custom Lua scripts to ensure race-condition-free counting across distributed clusters.
- **gRPC Interface:** Replaces slow JSON/HTTP with strictly typed, binary Protobuf messages.
- **Fail-Open Strategy:** Designed to allow traffic default if the rate-limiting service experiences downtime (availability over consistency).

## How to Run
**1. Start Redis**
```bash
sudo service redis-server start