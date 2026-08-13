

---

# Architecture Documentation

## Overview
The platform is built with scalability, reliability, and maintainability in mind, using a modular structure that separates concerns across distinct components, making development, testing, and deployment easier to manage independently.

# Components

#### 1. Backend Server

The core of the platform, handling HTTP requests, business logic, and data access.

**Stack:** Go, Gorilla Mux for routing, standard library HTTP handling.
**Responsibilities:**
* User authentication and authorization
* Ticket lifecycle management (creation, retrieval, closure)
* MySQL integration for persistent storage
* RESTful API endpoints for client-server communication
* Error handling and logging

#### 2. Database

Stores users, tickets, conversations, and access tokens.
**Stack:** MySQL.
**Schema:**

* Tables: `Users`, `Tickets`, `Conversations`, `AccessTokens`
* One-to-many relationships between users → tickets, and tickets → conversations
* Indexes on frequently queried columns to keep lookups fast

#### 3. Authentication

Custom middleware handling authentication and token-based authorization.

**Stack:** Go, JWT.
**Functionality:**
* Token generation and validation using JSON Web Tokens
* Role-based access control (RBAC) to distinguish regular users from admins
* Refresh tokens to extend session validity without repeated logins

#### 4. External Services

**Stack:** Mailtrap (email simulation).

Used for sending notification emails on ticket updates and account actions — Mailtrap stands in for a real SMTP provider so the flow can be tested without sending to real inboxes.

### Deployment

Currently deployed on Render, connecting to a MySQL database for storage. Render was chosen for simplicity during development; the app itself isn't tied to it and could move to any host running Go with access to a MySQL instance.

### Future Considerations

This project demonstrates the core architecture, but a production deployment would need further work:

* **Scaling:** horizontal scaling to handle increased traffic
* **Security:** HTTPS, stricter input validation, more robust password storage
* **Monitoring:** real-time performance tracking and error detection
* **Containerization:** Docker, for portability and consistency across environments
* **CI/CD:** automated testing and deployment pipeline
* **Data protection:** compliance with relevant data protection regulations
* **API testing:** a proper test suite covering functionality and performance
* **Auth improvements:** more granular, fine-grained authorization controls

--- 

