
---

# Deployment Guide

This guide provides step-by-step instructions on deploying the backend application in various environments.

## Prerequisites

Before deploying the application, ensure you have the following prerequisites:

- Access to the backend application codebase
- Access to a MySQL database server (e.g., Alwaysdata)
- Access to a mail testing service (e.g., Mailtrap)
- Environment variables configured for database connection and email service credentials

## Deployment Process

### 1. Setup MySQL Database

1. Access your MySQL database server (e.g., Alwaysdata) and create a new database.
2. Use the provided SQL schemas to create the necessary tables (`users`, `tickets`, `conversations`, `access_tokens`).

Example SQL commands:

```sql
CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `password` varchar(255) NOT NULL,
  `pin_number` varchar(255) DEFAULT NULL,
  `user_active` int(11) DEFAULT NULL,
  `is_admin` int(11) DEFAULT NULL,
  `refreshJWT` varchar(255) DEFAULT NULL
);

```

```sql
CREATE TABLE `access_tokens` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `accessJWT` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `expires_at` timestamp NULL DEFAULT NULL
);
```

```sql
CREATE TABLE `tickets` (
  `id` int(11) NOT NULL,
  `userId` bigint(20) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `subject` varchar(255) NOT NULL,
  `issue` varchar(255) DEFAULT NULL,
  `status` varchar(50) NOT NULL,
  `dateOpened` timestamp NULL DEFAULT current_timestamp()
) 
```

```sql
CREATE TABLE `conversations` (
  `id` int(11) NOT NULL,
  `ticketId` int(11) NOT NULL,
  `sender` varchar(255) NOT NULL,
  `message` text NOT NULL,
  `messageSentAt` timestamp NULL DEFAULT current_timestamp()
);
```

### 2. Configure Environment Variables

Ensure the following environment variables are set:

- `DB_USER`: MySQL database username
- `DB_PASSWORD`: MySQL database password
- `DB_HOST`: MySQL database host
- `DB_PORT`: MySQL database port
- `DB_DATABASE`: MySQL database name
- `SMTP_HOST`: Mailtrap SMTP host
- `SMTP_PORT`: Mailtrap SMTP port
- `SMTP_USERNAME`: Mailtrap username
- `SMTP_PASSWORD`: Mailtrap password
- `JWT_ACCESS_KEY`: JWT access key
- `JWT_REFRESH_KEY`: JWT refresh key

### 3. Deploy Backend Application

1. Clone the backend application repository.
2. Configure the necessary environment variables in a `.env` file.
3. Build the application using the appropriate build command (if applicable).
4. Start the application server.

### 4. Testing

Verify that the deployed backend application is running correctly by testing the endpoints using a tool like Postman.

---
