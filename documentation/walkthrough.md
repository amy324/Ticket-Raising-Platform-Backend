

---

# Walkthrough: Using the Ticket Platform API

Welcome to the walkthrough of the Ticket Platform API. In this guide, we'll explore the functionality of the API endpoints, allowing you to register as a new user, verify your account, log in, raise tickets, interact with conversations, and more. T
## Prerequisites

Before diving into the walkthrough, ensure you have an API testing tool like Postman installed. Alternatively, you can directly interact with the deployed API instance at https://ticketplatform.onrender.com/. Confirm the server's operation by visiting the link; if you see `{"message": "Hello, World!"}`, the server is up and running.

For security reasons, please refrain from using real information during testing. Instead, employ dummy data for all interactions.

## Registering as a New User

To register as a new user, send a POST request to the registration endpoint (`/register`) with the following JSON body. Use similar dummy data to below:

```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "firstName": "John",
  "lastName": "Doe"
}
```

You will receive a response confirming successful registration along with a PIN code. Note that in real-world scenarios, the PIN code would typically be sent via email for verification purposes. While this PIN is sent to the Mailtrap account I configured, I included it in the JSON response a syou will not have access to this unless you configure your own mailer

Response:
```
{
    "message": "User registered successfully",
    "pin": "517945",
    "userID": 29
}
```

Email recieved via Mailtrap:
```From: ticketplatform@email.com
To: user@example.com
Subject: Your PIN Code

Your PIN code is: Verification code for user user@example.com: 517945
```

The user has now been added to the MySQL database's `user` table:

| id | email            | first_name | last_name | password                                                     | pin_number | user_active | is_admin | refreshJWT                                                                                                                                    |
|----|------------------|------------|-----------|--------------------------------------------------------------|------------|-------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| 29 | user@example.com | John       | Doe       | hashed password | 517945     | 0           | 0        | refreshJWT |

## Verifying Your Account

To verify your account, send a POST request to the verification endpoint (`/verify-pin`) with the PIN code received in the registration response and your email:

```json
{
    "pin": "517945",
    "email": "user@example.com"
}
```

Upon successful verification, you'll receive a confirmation message.

Response:
```
{
    "message": "PIN verified successfully"
}
```

Attempting to login before verifying will result in the following response:
```
User does not exist or has not been activated. Please try re-registering your account
```
Successful verification updates the user's entry in the `users` table, meaning we can successfully login with this user's details:
| id | email            | first_name | last_name | password                                                     | pin_number | user_active | is_admin | refreshJWT                                                                                                                                    |
|----|------------------|------------|-----------|--------------------------------------------------------------|------------|-------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| 29 | user@example.com | John       | Doe       | hashed password | N/A - verified     | 1        | 0        | refreshJWT |


## Logging In

Now that your account is verified, you can log in by sending a POST request to the login endpoint (`/login`) with your email and password:

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

You will receive an access token and a refresh token, allowing you to access protected routes.

Response:
```
{
    "accessToken": <accessJWT>,
    "message": "Login successful",
    "refreshToken": <refreshJWT>,
    "user": {
        "ID": 29,
        "Email": "user@example.com",
        "FirstName": "John",
        "LastName": "Doe",
        "Password": "<hashed password>",
        "PinNumber": "",
        "UserActive": 1,
        "IsAdmin": 0,
        "RefreshJWT": ""
    }
}
```
Note in a real-world scenario access tokens should never be shared, but in the above response the actual tokens will be printed instead of the accessJWT and refreshJWT placeholders. Once again, please do not enter any real details when testing any of this program's features!

An entry for the `access_tokens` table has now been created:
| id | user_id | email            | accessJWT   | created_at          | updated_at          | expires_at          |   |
|----|---------|------------------|-------------|---------------------|---------------------|---------------------|---|
| 4  | 29      | user@example.com | accessJWT   | 2024-02-14 16:40:49 | 2024-02-14 16:40:49 | 2024-02-14 16:10:48 |   |

## Accessing User Profile

To access your user profile, send a GET request to the profile endpoint (`/profile`) with the access token provided during login. This simulates navigating to your profile page or dashboard etc. Include the access token in the Authorization header:

```
Authorization:  <access token>
```

You'll receive a response containing your user details:
```{
    "ID": 29,
    "Email": "user@example.com",
    "FirstName": "John",
    "LastName": "Doe",
    "Password": <hashed password>,
    "PinNumber": "",
    "UserActive": 1,
    "IsAdmin": 0,
    "RefreshJWT": <refreshJWT>
}
```

## Raising Tickets

To raise a new ticket, send a POST request to the tickets endpoint (`/tickets`) with the subject and issue details in the request body. Use the access token obtained during login for authorization.

Example body:
```
{
    "subject": "Subject 1",
    "issue": "Issue 1"
}

```
This returns the following response:

```
{
    "TicketID": 3
}
```
An entry for this ticket has now been created in the `tickets` table in the MySQL database:

| id | user_id | email            | subject   | issue   | status | dateOpened          |
|----|---------|------------------|-----------|---------|--------|---------------------|
| 3  | 29      | user@example.com | Subject 1 | Issue 1 | open   | 2024-02-14 16:50:25 |



## Viewing Tickets

You can view all tickets associated with your account by sending a GET request to the tickets endpoint (`/tickets`). Again, include the access token in the Authorization header.

The response:
```
[
    {
        "id": 3,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 1",
        "issue": "Subject 1",
        "status": "open",
        "dateOpened": "2024-02-14T16:50:25Z"
    },
    {
        "id": 4,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 2",
        "issue": "Subject 2",
        "status": "open",
        "dateOpened": "2024-02-14T16:52:24Z"
    },
    {
        "id": 5,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 3",
        "issue": "Subject 3",
        "status": "open",
        "dateOpened": "2024-02-14T16:52:32Z"
    }
]

```
You can view an individual ticket associated with your account, send a GET request to the tickets endpoint with the individual ticket ID (in this case 3) (`/tickets/{ticketid}`). Again, include the access token in the Authorization header.
```
{
    "ticket": {
        "id": 3,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 1",
        "issue": "Issue 1",
        "status": "open",
        "dateOpened": "2024-02-14T16:50:25Z"
    },
    "conversations": [
        {
            "id": 5,
            "ticketId": 3,
            "sender": "operator",
            "message": "We will be in touch with you shortly. In the meantime please feel free to reply to this message with more details",
            "messageSentAt": "2024-02-14T16:50:26Z"
        }
    ]
}
```

An automated operater message gets sent to each ticket.

Note, you are unable to view tickets associated with other accounts. If you try using a different ticket ID, you will get the following response, regardless of if another user has a ticket with that ID or not:

```
No ticket associated with this ID
```


## Interacting with Conversations

To add a new message to a ticket conversation, send a POST request to the conversations endpoint (`/tickets/{ticketID}/conversation`) with the message content in the request body. Use the access token for authorization.

Body:

```
{
    "message": "Reply 1"
}
```

Giving the response:

```
Message successfully sent
```

An entry has been created for this message in the `conversations` table in the database:
| id | ticketId | sender   | message                                                                                                           | messageSentAt       |
|----|----------|----------|-------------------------------------------------------------------------------------------------------------------|---------------------|
| 5  | 3        | operator | We will be in touch with you shortly. In the meantime please feel free to reply to this message with more details | 2024-02-14 16:50:26 |
| 8  | 3        | John     | Reply 1                                                                                                           | 2024-02-14 15:59:29 |

Therefore, if you try the GET (`/tickets/{ticketid}`) endpoint, this message will now be displayed:
```
{
    "ticket": {
        "id": 3,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 1",
        "issue": "Issue 1",
        "status": "open",
        "dateOpened": "2024-02-14T16:50:25Z"
    },
    "conversations": [
        {
            "id": 5,
            "ticketId": 3,
            "sender": "operator",
            "message": "We will be in touch with you shortly. In the meantime please feel free to reply to this message with more details",
            "messageSentAt": "2024-02-14T16:50:26Z"
        },
        {
            "id": 8,
            "ticketId": 3,
            "sender": "John",
            "message": "Reply 1",
            "messageSentAt": "2024-02-14T15:59:29Z"
        }
    ]
}
```

## Closing Tickets

To close a ticket, send a DELETE request to the ticket endpoint (`/tickets/{ticketID}`) - in this case I used ticketID 3. Ensure you have the necessary permissions and include the access token in the Authorization header.

This will generate the following response:
```
Ticket 3 closed successfully
```
All entries associated with this ticket in the `tickets` and `conversation` tables have now been deleted from the database
## Refreshing Tokens

If your access token expires, you can refresh it by sending a POST request to the token refresh endpoint (`/tokens/refresh`) with your refresh token as the Authorization header:
```
Authorization:  <refresh token>
```

This generates the following response:
```
{
    "accessToken": "<access token>,
    "message": "Token refreshed successfully"
}
```

The accessJWT, updatedAt and expires_at values for the user in the `access_tokens` table will now update


## Logging Out

To log out of your account, send a POST request to the logout endpoint (`/logout`) with the access token in the Authorization header. If you refreshed, ensure you use the new access token. 

Response:

```
{
    "message": "Logout successful"
}
```

The `access_tokens` entry for this user has now been deleted. Logging in again, will update the refershJWT entry in the `users table` as well as create an `access_tokens` table entry.

## Admin Privileges

If you're an admin user, you have additional privileges such as viewing all tickets raised by users and messaging users directly. You can log in as an admin user using the same login endpoint (`/login`).

The `users` entry for an admin user in the database:

| id | email            | first_name | last_name | password                                                     | pin_number | user_active | is_admin | refreshJWT                                                                                                                                    |
|----|------------------|------------|-----------|--------------------------------------------------------------|------------|-------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| 5 | admin@ticketplatform.com| Admin       | User     | hashed password | N/A - verified    | 1           | 1       | refreshJWT |

To test, you will need to manually add an admin user to your MySQL database.

### Admin View All Open Tickets
An admin user can send a GET request to `/admin/tickets`, using their access token as the Authorization Header, to view all tickets raised by users. It pulls all the tickets from the `tickets` table:


```
[
    {
        "id": 2,
        "userId": 4,
        "email": "testuser@example.com",
        "subject": "ssl issue",
        "issue": "ssl certificate not recognised",
        "status": "open",
        "dateOpened": "2024-02-13T15:57:13Z"
    },
    {
        "id": 4,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 2",
        "issue": "Subject 2",
        "status": "open",
        "dateOpened": "2024-02-14T16:52:24Z"
    },
    {
        "id": 5,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 3",
        "issue": "Subject 3",
        "status": "open",
        "dateOpened": "2024-02-14T16:52:32Z"
    }
]
```
## Admin View Individual Ticket & Conversations
The admin can view a specific ticket and its coversations with a GET request to `/admin/tickets/{ticketID}` - in this case I used ticketID 4. Again, including the access token in the Authorization header.

```
{
    "ticket": {
        "id": 4,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 2",
        "issue": "Subject 2",
        "status": "open",
        "dateOpened": "2024-02-14T16:52:24Z"
    },
    "conversations": [
        {
            "id": 6,
            "ticketId": 4,
            "sender": "operator",
            "message": "We will be in touch with you shortly. In the meantime please feel free to reply to this message with more details",
            "messageSentAt": "2024-02-14T16:52:24Z"
        }
    ]
}
```

## Admin Add a Conversation To a Ticker
The admin user can also send users a message for a specific ticket by sending a POST request to `/admin/tickets/{ticketID}/conversation` - in this case 4. As always, using the access token for Authorization.

Body:

```
{"message": "Please provide more details of your situation "}
```

Response:
```
{
    "conversationID": 9
}
```

This message will now be added to the `conversations` table in the database, available for the user to viuew when they next make a GET (`/tickets/{ticketID}`) request for the ticket in question:


```
{
    "ticket": {
        "id": 4,
        "userId": 29,
        "email": "user@example.com",
        "subject": "Subject 2",
        "issue": "Subject 2",
        "status": "open",
        "dateOpened": "2024-02-14T16:52:24Z"
    },
    "conversations": [
        {
            "id": 6,
            "ticketId": 4,
            "sender": "operator",
            "message": "We will be in touch with you shortly. In the meantime please feel free to reply to this message with more details",
            "messageSentAt": "2024-02-14T16:52:24Z"
        },
        {
            "id": 9,
            "ticketId": 4,
            "sender": "operator",
            "message": "Please provide more details of your situation ",
            "messageSentAt": "2024-02-14T16:33:40Z"
        }
    ]
}


```
All messages sent by an admin user are automatically given the sender name "operator". 

Note: if a non-admin user tries to make any admin router requests, they will get the following response:

```
Access denied. Admin privilege required.
```

---

This concludes the walkthrough of the Ticket Platform API. Feel free to explore the endpoints further and test different scenarios using your preferred API testing tool.