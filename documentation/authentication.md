
---

# User Authentication Documentation

This document provides an in-depth overview of the user authentication mechanism integrated into the ticket-raising platform project. Employing JWT (JSON Web Tokens), this authentication system ensures secure access to resources while preventing unauthorized access to sensitive data.

## Authentication Flow

The authentication flow comprises the following steps:

1. **User Registration**: Users register by providing essential details such as email and password, optionally including first and last names. Upon successful validation, a new user account is created in the database.

```go
// RegisterHandler manages user registration requests.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // Extract, validate, and sanitize user registration details
    // Create a new user account in the database
}
```

2. **Pin Verification**: A 6-digit PIN is sent to the user's provided email address for account activation. Upon successful PIN verification, the user account is marked as active in the database and ready for login.

```go
// VerifyPinHandler handles PIN verification.
func VerifyPinHandler(w http.ResponseWriter, r *http.Request) {
    // Retrieve user by email
    // Retrieve PIN for the user from the database
    // Compare the provided PIN with the stored one
}
```

3. **User Login**: Registered users can log in using their email and password. The server verifies the credentials and issues JWT tokens upon successful authentication.

```go
// LoginHandler manages user login requests.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // Extract login credentials from the request
    // Verify credentials against database records
    // Generate JWT tokens upon successful authentication
}
```

4. **JWT Generation**: After successful login, the server generates JWT access and refresh tokens. The access token contains user identification details and provides access to protected resources. The refresh token facilitates obtaining a new access token without requiring reauthentication.

```go
// GenerateJWT generates a new JWT token.
func GenerateJWT(user User) (string, string, error) {
    // Generate access token with user claims
    // Generate refresh token associated with the user
    // Return both tokens
}
```

5. **Access Control**: Protected endpoints require authentication, enforced by middleware functions. These functions validate JWT tokens and grant access only to authenticated users.

```go
// validateAccessToken validates JWT access tokens.
func validateAccessToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract and validate JWT access token from request header
        // Grant access to protected resource if token is valid
        // Otherwise, deny access with appropriate error response
    })
}
```

6. **Token Refresh**: When an access token expires, users can use the refresh token to obtain a new access token without logging in again. This mechanism enhances user experience by minimizing the need for frequent logins.

```go
// refreshAccessToken generates a new access token using a refresh token.
func refreshAccessToken(w http.ResponseWriter, r *http.Request, refreshToken string, db *sql.DB) {
    // Validate the refresh token and retrieve associated user information
    // Generate a new access token if the refresh token is valid
    // Return the new access token to the client
}
```

7. **Admin Privileges**: Certain functionalities are restricted to admin users, requiring additional privileges. Admin status is assigned during user registration or manually configured in the database.

```go
// ViewAllTicketsHandler retrieves all tickets from the database.
func ViewAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
    // Authenticate the request and verify admin privileges
    // Retrieve all tickets from the database
    // Return ticket data as JSON response
}
```

## User Privileges

1. **User Roles**: Each user account is assigned a role determining access levels and privileges. Roles include:
   - **User**: Standard users who can raise tickets, view their own tickets, and perform basic operations.
   - **Admin**: Administrators with additional privileges such as viewing all tickets, adding conversations to any ticket, and closing tickets.

2. **Authorization**: Endpoint access is restricted based on the user's role:
   - Regular users can manage their own tickets.
   - Administrators have access to all tickets and perform administrative tasks.

3. **JWT Claims**: JWT tokens include claims specifying the user's role and permissions, verified server-side to enforce access control policies.

## Implementation Details

- **JWT Token Generation**: Secure JWT tokens are generated using robust algorithms and signed with a secret key to prevent tampering.
- **Token Expiry**: Access tokens have a short lifespan (e.g., 15 minutes), while refresh tokens remain valid for a longer duration (e.g., 7 days) to mitigate misuse risks.
- **Token Refresh**: Upon access token expiry, clients utilize refresh tokens to obtain new access tokens without reauthentication.
- **Authorization Middleware**: Middleware functions enforce authentication and authorization checks on protected endpoints, validating JWT tokens and verifying user roles.

## Security Considerations

- **Secure Storage**: Sensitive data, including user credentials, are securely stored in the database using hashed passwords and encryption techniques.
- **HTTPS**: API endpoints are served over HTTPS to maintain data confidentiality and integrity during transmission.
- **Input Validation**: Request data undergoes validation and sanitization to prevent common security vulnerabilities like SQL injection and XSS attacks.

---

## Conclusion

The user authentication system integrated into the ticket-raising platform ensures secure resource access and protects sensitive data. Leveraging JWT tokens and access control mechanisms, the platform maintains user account integrity and confidentiality, providing a seamless user experience.

--- 

