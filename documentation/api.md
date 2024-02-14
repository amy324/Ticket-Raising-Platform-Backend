
---

# Ticket Platform API Documentation

Welcome to the Ticket Platform API documentation. This API provides endpoints for managing user authentication, ticket creation, ticket conversations, and administrative tasks.

## Base URL

The base URL for all endpoints is:

```
https://ticketplatform.onrender.com/
```

## Authentication

### Register

- **URL**: `/register`
- **Method**: `POST`
- **Description**: Register a new user account.
- **Request Body**:
  - `email` (string): User's email address.
  - `password` (string): User's password.
  - `firstName` (string): User's first name.
  - `lastName` (string): User's last name.
- **Response**: 
  - `200 OK`: User successfully registered.
  - `400 Bad Request`: Invalid request body.

### Login

- **URL**: `/login`
- **Method**: `POST`
- **Description**: Log in to an existing user account.
- **Request Body**:
  - `email` (string): User's email address.
  - `password` (string): User's password.
- **Response**: 
  - `200 OK`: User successfully logged in. Returns access and refresh tokens.
  - `401 Unauthorized`: Invalid credentials.

### Logout

- **URL**: `/logout`
- **Method**: `POST`
- **Description**: Log out a user from their account.
- **Request Header**:
  - `Authorization` (string): Bearer token.
- **Response**: 
  - `200 OK`: User successfully logged out.
  - `401 Unauthorized`: Invalid or expired token.

### Refresh Token

- **URL**: `/tokens/refresh`
- **Method**: `POST`
- **Description**: Refresh the access token using a valid refresh token.
- **Request Header**:
  - `Authorization` (string): Refresh token.
- **Response**: 
  - `200 OK`: Access token successfully refreshed.
  - `401 Unauthorized`: Invalid or expired refresh token.

## Tickets

### Create Ticket

- **URL**: `/tickets`
- **Method**: `POST`
- **Description**: Create a new support ticket.
- **Request Body**:
  - `subject` (string): Subject of the ticket.
  - `issue` (string): Description of the issue.
- **Response**: 
  - `200 OK`: Ticket successfully created.
  - `400 Bad Request`: Invalid request body.

### Get All Tickets

- **URL**: `/tickets`
- **Method**: `GET`
- **Description**: Retrieve all support tickets.
- **Response**: 
  - `200 OK`: List of tickets retrieved successfully.

### Get Ticket by ID

- **URL**: `/tickets/{ticketID}`
- **Method**: `GET`
- **Description**: Retrieve a support ticket by its ID.
- **Response**: 
  - `200 OK`: Ticket retrieved successfully.
  - `404 Not Found`: Ticket not found.

### Close Ticket

- **URL**: `/tickets/{ticketID}`
- **Method**: `DELETE`
- **Description**: Close a support ticket by its ID.
- **Response**: 
  - `200 OK`: Ticket successfully closed.
  - `404 Not Found`: Ticket not found.

### Add Conversation to Ticket

- **URL**: `/tickets/{ticketID}/conversation`
- **Method**: `POST`
- **Description**: Add a new conversation message to a support ticket.
- **Request Body**:
  - `message` (string): Message to add to the conversation.
- **Response**: 
  - `200 OK`: Conversation message added successfully.
  - `400 Bad Request`: Invalid request body.

## Administration

### View All Tickets (Admin)

- **URL**: `/admin/tickets`
- **Method**: `GET`
- **Description**: Retrieve all support tickets (admin access required).
- **Response**: 
  - `200 OK`: List of tickets retrieved successfully.
  - `403 Forbidden`: Access denied.

### Get Ticket by ID (Admin)

- **URL**: `/admin/tickets/{ticketID}`
- **Method**: `GET`
- **Description**: Retrieve a support ticket by its ID (admin access required).
- **Response**: 
  - `200 OK`: Ticket retrieved successfully.
  - `403 Forbidden`: Access denied.
  - `404 Not Found`: Ticket not found.

### Add Conversation to Ticket (Admin)

- **URL**: `/admin/tickets/{ticketID}/conversation`
- **Method**: `POST`
- **Description**: Add a new conversation message to a support ticket (admin access required).
- **Request Body**:
  - `message` (string): Message to add to the conversation.
- **Response**: 
  - `200 OK`: Conversation message added successfully.
  - `400 Bad Request`: Invalid request body.
  - `403 Forbidden`: Access denied.
  - `404 Not Found`: Ticket not found.

---
