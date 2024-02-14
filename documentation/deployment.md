
---

# Deployment Guide

This guide offers comprehensive instructions for deploying the backend application across different environments.

## Prerequisites

Before proceeding with the deployment, ensure the following prerequisites are met:

- Access to the backend application's codebase via cloning (see [README](../README.md) for details).
- Access to a MySQL database server (e.g., Alwaysdata)
- Access to a mail testing service (e.g., Mailtrap)
- Proper configuration of environment variables for connecting to the database and email service

## Deployment Process

### 1. Setup MySQL Database

1. Access your MySQL database server (e.g., Alwaysdata) and create a new database.
2. Utilize the provided SQL schemas to create essential tables (`users`, `tickets`, `conversations`, `access_tokens`).

Example SQL commands:

```sql
-- Create the 'users' table
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

-- Create the 'access_tokens' table
CREATE TABLE `access_tokens` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `accessJWT` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `expires_at` timestamp NULL DEFAULT NULL
);

-- Create the 'tickets' table
CREATE TABLE `tickets` (
  `id` int(11) NOT NULL,
  `userId` bigint(20) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `subject` varchar(255) NOT NULL,
  `issue` varchar(255) DEFAULT NULL,
  `status` varchar(50) NOT NULL,
  `dateOpened` timestamp NULL DEFAULT current_timestamp()
);

-- Create the 'conversations' table
CREATE TABLE `conversations` (
  `id` int(11) NOT NULL,
  `ticketId` int(11) NOT NULL,
  `sender` varchar(255) NOT NULL,
  `message` text NOT NULL,
  `messageSentAt` timestamp NULL DEFAULT current_timestamp()
);
```

### 2. Configure Environment Variables

Ensure the following environment variables are properly configured:

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

1. Clone the repository containing the backend application.
2. Configure the required environment variables in a `.env` file.
3. If applicable, build the application using the appropriate build command, I used `go build` to build `web.exe` in my case.
4. Start the application server.

Alternatively, simply run on a localhost using the `go run` command.

Alternatively, feel free 

### 4. Testing

Verify the proper functioning of the deployed backend application by testing the endpoints using tools like Postman.

---
