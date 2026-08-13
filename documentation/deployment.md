# Deployment Instructions

Instructions for deploying the backend to a new environment.

### Prerequisites

* The codebase, cloned locally (see README)
* A MySQL database server (e.g. Alwaysdata)
* A mail testing service (e.g. Mailtrap)
* Environment variables configured for both

### 1. Set Up the MySQL Database

Create a new database on your MySQL server, then run the schema below to set up the required tables.

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

| Variable          | Purpose                       |
| ----------------- | ----------------------------- |
| `DB_USER`         | MySQL username                |
| `DB_PASSWORD`     | MySQL password                |
| `DB_HOST`         | MySQL host                    |
| `DB_PORT`         | MySQL port                    |
| `DB_DATABASE`     | MySQL database name           |
| `SMTP_HOST`       | Mailtrap SMTP host            |
| `SMTP_PORT`       | Mailtrap SMTP port            |
| `SMTP_USERNAME`   | Mailtrap username             |
| `SMTP_PASSWORD`   | Mailtrap password             |
| `JWT_ACCESS_KEY`  | JWT access token signing key  |
| `JWT_REFRESH_KEY` | JWT refresh token signing key |

### 3. Deploy the Application

1. Clone the repository.
2. Add the environment variables above to a `.env` file.
3. Build the app. I used `go build` to produce `web.exe` locally, but any standard Go build process will work.
4. Start the server, or run it directly with `go run` for local testing.

### 4. Test

Verify the deployment by hitting the endpoints with a tool like Postman.
