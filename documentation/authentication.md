# User Authentication

An overview of the authentication system behind the platform, built around JWTs to keep resource access secure without requiring users to log in repeatedly.

### Authentication Flow

**User registration:** users sign up with an email and password, optionally including their name. Once validated, a new account is created.

```go
// RegisterHandler manages user registration requests.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // Extract, validate, and sanitize user registration details
    // Create a new user account in the database
}
```

**PIN verification:** a 6-digit PIN is emailed to the user to activate their account. Once verified, the account is marked active and ready for login.

```go
// VerifyPinHandler handles PIN verification.
func VerifyPinHandler(w http.ResponseWriter, r *http.Request) {
    // Retrieve user by email
    // Retrieve PIN for the user from the database
    // Compare the provided PIN with the stored one
}
```

**Login:** registered users log in with email and password. On success, the server issues JWT tokens.

```go
// LoginHandler manages user login requests.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // Extract login credentials from the request
    // Verify credentials against database records
    // Generate JWT tokens upon successful authentication
}
```

**JWT generation:** after login, the server issues an access token (used to reach protected resources) and a refresh token (used to get a new access token without logging in again).

```go
// GenerateJWT generates a new JWT token.
func GenerateJWT(user User) (string, string, error) {
    // Generate access token with user claims
    // Generate refresh token associated with the user
    // Return both tokens
}
```

**Access control:** protected endpoints are gated by middleware that validates the JWT before allowing access.

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

**Token refresh** — when an access token expires, the refresh token can be used to get a new one without re-authenticating.

```go
// refreshAccessToken generates a new access token using a refresh token.
func refreshAccessToken(w http.ResponseWriter, r *http.Request, refreshToken string, db *sql.DB) {
    // Validate the refresh token and retrieve associated user information
    // Generate a new access token if the refresh token is valid
    // Return the new access token to the client
}
```

**Admin privileges:** some functionality is restricted to admins, assigned either at registration or manually in the database.

```go
// ViewAllTicketsHandler retrieves all tickets from the database.
func ViewAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
    // Authenticate the request and verify admin privileges
    // Retrieve all tickets from the database
    // Return ticket data as JSON response
}
```

### User Roles

* **User:** can raise tickets, view their own tickets, and perform basic operations.
* **Admin:** can view all tickets, add conversations to any ticket, and close tickets.

Role is embedded in the JWT claims and checked server-side on every protected request.

### Implementation Notes

* Tokens are signed with a secret key to prevent tampering.
* Access tokens are short-lived (e.g. 15 minutes); refresh tokens last longer (e.g. 7 days) to limit exposure if one is compromised.
* Middleware enforces both authentication and role checks on protected endpoints.

### Security

* Passwords are hashed before storage.
* API endpoints are served over HTTPS.
* Request data is validated and sanitized to guard against SQL injection and XSS.
