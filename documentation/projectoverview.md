# Introduction

The Ticket-Raising Platform is a backend application built to demonstrate Go backend development practices, covering the design, implementation, and deployment of a scalable, secure system for managing support tickets.

### Objectives

* **Backend development:** design and build a robust, feature-rich backend in Go
* **Scalability:** architect the system to handle growing traffic and ticket volume
* **Security:** protect user data with secure authentication, authorization, and defenses against common vulnerabilities
* **API design:** a clean, well-documented RESTful interface

### Features

* **Authentication & authorisation:** token-based auth enforcing access control across the app
* **Ticket management:** create, retrieve, update, and close tickets, with categorization, status tracking, and queue management
* **Real-time communication:** in-thread messaging between users and support staff
* **Email notifications:** — keeps users informed of ticket updates and status changes
* **Data persistence:** MySQL for storing users, tickets, conversations, and access tokens
* **Logging & error handling:** structured logging to support debugging and troubleshooting
