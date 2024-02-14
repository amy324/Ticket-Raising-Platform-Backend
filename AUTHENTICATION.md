

---

# User Authentication Documentation

This document outlines the user authentication mechanism implemented in the ticket-raising platform project. The authentication system employs JWT (JSON Web Tokens) for user authentication and authorization. It ensures secure access to resources and prevents unauthorized access to sensitive data.

## Authentication Flow

The authentication flow involves the following steps:

1. **User Registration**: Users register with the platform by providing necessary details such as email, password, and optionally first name and last name. The provided information is validated, and if successful, a new user account is created in the database.

```go
// RegisterHandler handles user registration requests.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // Extract user registration details from the request body
    // Validate and sanitize the input
    // Create a new user account in the database
}
```

2. **User Login**: Upon successful registration, users can log in to the platform using their registered email and password. The server verifies the credentials and generates JWT tokens for authenticated users.

```go
// LoginHandler handles user login requests.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // Extract login credentials from the request body
    // Verify the credentials against stored data in the database
    // Generate JWT tokens upon successful authentication
}
```

3. **JWT Generation**: Upon successful login, the server generates JWT access and refresh tokens. The access token contains user identification information and is used to access protected resources. The refresh token allows users to obtain a new access token without requiring reauthentication.

```go
// GenerateJWT generates a new JWT token.
func GenerateJWT(user User) (string, string, error) {
    // Generate access token with user claims
    // Generate refresh token and associate it with the user
    // Return both tokens
}
```

4. **Access Control**: Endpoints requiring authentication are protected using middleware functions. These middleware functions validate JWT tokens sent by clients and grant access only to authenticated users.

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

5. **Token Refresh**: When an access token expires, users can use the refresh token to obtain a new access token without having to log in again. This mechanism improves user experience by minimizing the need for frequent logins.

```go
// refreshAccessToken generates a new access token using a refresh token.
func refreshAccessToken(w http.ResponseWriter, r *http.Request, refreshToken string, db *sql.DB) {
    // Validate the refresh token and retrieve associated user information
    // Generate a new access token if the refresh token is valid
    // Return the new access token to the client
}
```

6. **Admin Privileges**: Certain endpoints and functionalities are restricted to admin users. Admin privileges are assigned during user registration or through manual configuration in the database.

```go
// ViewAllTicketsHandler retrieves all tickets from the database.
func ViewAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
    // Authenticate the request and verify admin privileges
    // Retrieve all tickets from the database
    // Return ticket data as JSON response
}
```

## User Privileges

1. **User Roles**: Each user account has a role assigned to it, determining the level of access and privileges. The roles include:
   - **User**: Standard users who can raise tickets, view their own tickets, and perform other basic operations.
   - **Admin**: Administrators who have additional privileges such as viewing all tickets, adding conversations to any ticket, and closing tickets.

2. **Authorization**: Access to certain endpoints and resources is restricted based on the user's role. For example:
   - Regular users can only view and manage their own tickets.
   - Administrators have access to all tickets and can perform administrative actions.

3. **JWT Claims**: JWT tokens contain claims that specify the user's role and permissions. These claims are verified on the server side to enforce access control policies.

## Implementation Details

- **JWT Token Generation**: JWT tokens are generated using a secure algorithm and signed with a secret key to prevent tampering.
- **Token Expiry**: Access tokens have a short expiry time (e.g., 15 minutes), while refresh tokens have a longer validity period (e.g., 7 days). This helps mitigate the risk of token misuse.
- **Token Refresh**: When the access token expires, the client can send the refresh token to the server to obtain a new access token without requiring the user to log in again.
- **Authorization Middleware**: Middleware functions are used to enforce authentication and authorization checks on protected endpoints. These middleware functions validate the JWT tokens and verify the user's role before allowing access to the requested resource.

## Security Considerations

- **Secure Storage**: User credentials and sensitive information are securely stored in the database using hashed passwords and other encryption techniques.
- **HTTPS**: API endpoints are served over HTTPS to ensure data confidentiality and integrity during transmission.
- **Input Validation**: Request data is validated and sanitized to prevent common security vulnerabilities such as SQL injection and cross-site scripting (XSS) attacks.

---



## Conclusion

The user authentication system implemented in the ticket-raising platform ensures secure access to resources, protects sensitive data, and provides a seamless user experience. By leveraging JWT tokens and access control mechanisms, the platform maintains the integrity and confidentiality of user accounts and information.

---

