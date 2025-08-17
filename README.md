# Order Pipeline — Async Order Processing

## What?

Project built to explore Dockerization, message queues, caching, and microservice architecture.

FastAPI service receives orders via HTTP and stores them in PostgreSQL.
Orders are published to RabbitMQ for asynchronous processing.
Go worker consumes orders from RabbitMQ, processes them, and caches results in Redis.

Decouples API from background processing for better performance and scalability.
Docker Compose manages all services for easy local setup.

Demonstrates common patterns for building scalable, resilient backend systems.

---

## Why?

When you get an order through an API, you don’t want the client to hang around while you do a bunch of stuff like validating payment, updating inventory, or sending emails. That takes time and can slow everything down.

So, instead of doing everything in one shot, this setup **decouples** the “accept order” part from the “process order” part. This way:

* Your API stays snappy and responsive.
* You can scale your processing workers separately if things get busy.
* If the processor crashes or needs to restart, orders don’t get lost — they’re safely queued.
* You can add retries and error handling in the processing step without complicating the API.

---

## What again?

* **FastAPI service**: This is the front door. It takes in order requests via HTTP, saves the order details in PostgreSQL, and pushes a message about the new order into RabbitMQ.

* **RabbitMQ**: Think of it as a reliable post office that holds your order messages. It queues them so they don’t get lost and lets worker services grab orders to process when they’re ready.

* **Go Worker**: This guy listens on RabbitMQ for new orders. When it gets one, it does the actual “work” — like updating status, maybe charging the payment (not implemented here, but you get the idea), and storing the processing result in Redis for quick access.

* **PostgreSQL**: Keeps your order data safe and persistent.

* **Redis**: Used for caching the processed results to serve them fast if needed.

---