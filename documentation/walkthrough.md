## Walkthrough: Using the Ticket Platform API

A step-by-step guide to the API — registering, verifying your account, logging in, raising tickets, and interacting with conversations.

### Contents

* [Prerequisites](#prerequisites)
* [Registering a New User](#registering-a-new-user)
* [Verifying Your Account](#verifying-your-account)
* [Logging In](#logging-in)
* [Accessing Your Profile](#accessing-your-profile)
* [Raising Tickets](#raising-tickets)
* [Viewing Tickets](#viewing-tickets)
* [Interacting with Conversations](#interacting-with-conversations)
* [Closing Tickets](#closing-tickets)
* [Refreshing Tokens](#refreshing-tokens)
* [Logging Out](#logging-out)
* [Admin Privileges](#admin-privileges)

### Prerequisites

As this project only contains backend code and no GUI, you'll want an API testing tool like Postman or you could make the requests via your command line terminal.  

**Note:** protected endpoints below require the access token in the `Authorization` header. This is called out explicitly for the first few endpoints, then assumed from there on. Please use dummy data throughout; avoid entering real personal information.

### Registering a New User

`POST /register`

```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "firstName": "John",
  "lastName": "Doe"
}
```

You'll get a confirmation along with a PIN code. In a real deployment this PIN would only go out by email — here it's included in the response too, since you likely don't have access to my Mailtrap inbox.

Response:

```json
{
    "message": "User registered successfully",
    "pin": "517945",
    "userID": 29
}
```

Email received via Mailtrap:

```text
To: user@example.com
Subject: Your PIN Code

Your PIN code is: Verification code for user user@example.com: 517945
```

The user now exists in the `users` table:

| id | email                                       | first_name | last_name | password        | pin_number | user_active | is_admin | refreshJWT |
| -- | ------------------------------------------- | ---------- | --------- | --------------- | ---------- | ----------- | -------- | ---------- |
| 29 | [user@example.com](mailto:user@example.com) | John       | Doe       | hashed password | 517945     | 0           | 0        | refreshJWT |

### Verifying Your Account

`POST /verify-pin`

```json
{
    "pin": "517945",
    "email": "user@example.com"
}
```

Response:

```json
{
    "message": "PIN verified successfully"
}
```

Trying to log in before verifying returns:

```text
User does not exist or has not been activated. Please try re-registering your account
```

After verification, the user's row updates and login becomes possible:

| id | email                                       | first_name | last_name | password        | pin_number     | user_active | is_admin | refreshJWT |
| -- | ------------------------------------------- | ---------- | --------- | --------------- | -------------- | ----------- | -------- | ---------- |
| 29 | [user@example.com](mailto:user@example.com) | John       | Doe       | hashed password | N/A - verified | 1           | 0        | refreshJWT |

### Logging In

`POST /login`

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Response:

```json
{
    "accessToken": "<accessJWT>",
    "message": "Login successful",
    "refreshToken": "<refreshJWT>",
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

In a real deployment these tokens should never be shared. Here the response shows the actual token strings rather than placeholders. Again, stick to dummy data when testing.

An entry is created in `access_tokens`:

| id | user_id | email                                       | accessJWT | created_at          | updated_at          | expires_at          |
| -- | ------- | ------------------------------------------- | --------- | ------------------- | ------------------- | ------------------- |
| 4  | 29      | [user@example.com](mailto:user@example.com) | accessJWT | 2024-02-14 16:40:49 | 2024-02-14 16:40:49 | 2024-02-14 16:10:48 |

### Accessing Your Profile

`GET /profile`, with the access token in the `Authorization` header. This simulates loading a profile page or dashboard.

Response:

```json
{
    "ID": 29,
    "Email": "user@example.com",
    "FirstName": "John",
    "LastName": "Doe",
    "Password": "<hashed password>",
    "PinNumber": "",
    "UserActive": 1,
    "IsAdmin": 0,
    "RefreshJWT": "<refreshJWT>"
}
```

### Raising Tickets

`POST /tickets`

```json
{
    "subject": "Subject 1",
    "issue": "Issue 1"
}
```

Response:

```json
{
    "TicketID": 3
}
```

The ticket is now stored in the `tickets` table:

| id | userId | email                                       | subject   | issue   | status | dateOpened          |
| -- | ------ | ------------------------------------------- | --------- | ------- | ------ | ------------------- |
| 3  | 29     | [user@example.com](mailto:user@example.com) | Subject 1 | Issue 1 | open   | 2024-02-14 16:50:25 |

### Viewing Tickets

`GET /tickets` returns every ticket tied to your account:

```json
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

For a single ticket, `GET /tickets/{ticketID}` (e.g. `/tickets/3`):

```json
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

Every new ticket gets an automated operator message.

You can't view another account's tickets — requesting a ticket ID that isn't yours (whether or not it exists) returns:

```text
No ticket associated with this ID
```

### Interacting with Conversations

`POST /tickets/{ticketID}/conversation`

```json
{
    "message": "Reply 1"
}
```

Response:

```text
Message successfully sent
```

The message is added to `conversations`:

| id | ticketId | sender   | message                                                                                                           | messageSentAt       |
| -- | -------- | -------- | ----------------------------------------------------------------------------------------------------------------- | ------------------- |
| 5  | 3        | operator | We will be in touch with you shortly. In the meantime please feel free to reply to this message with more details | 2024-02-14 16:50:26 |
| 8  | 3        | John     | Reply 1                                                                                                           | 2024-02-14 15:59:29 |

A follow-up `GET /tickets/3` now shows both messages:

```json
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

### Closing Tickets

`DELETE /tickets/{ticketID}` (e.g. `/tickets/3`), with the access token in the header and appropriate permissions.

Response:

```text
Ticket 3 closed successfully
```

All related rows in `tickets` and `conversations` are deleted.

### Refreshing Tokens

`POST /tokens/refresh`, with the refresh token in the `Authorization` header:

Response:

```json
{
    "accessToken": "<access token>",
    "message": "Token refreshed successfully"
}
```

The `accessJWT`, `updated_at`, and `expires_at` fields update in `access_tokens`.

### Logging Out

`POST /logout`, with the current access token in the header.

Response:

```json
{
    "message": "Logout successful"
}
```

The corresponding `access_tokens` entry is deleted. Logging in again updates `refreshJWT` in `users` and creates a new `access_tokens` row.

### Admin Privileges

Admin users log in through the same `/login` endpoint. An admin row looks like:

| id | email                                                       | first_name | last_name | password        | pin_number     | user_active | is_admin | refreshJWT |
| -- | ----------------------------------------------------------- | ---------- | --------- | --------------- | -------------- | ----------- | -------- | ---------- |
| 5  | [admin@ticketplatform.com](mailto:admin@ticketplatform.com) | Admin      | User      | hashed password | N/A - verified | 1           | 1        | refreshJWT |

To test this, you'll need to manually flag a user as admin in the database.

**Viewing all tickets:** `GET /admin/tickets`, with the admin's access token.

```json
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

**Viewing a specific ticket:** `GET /admin/tickets/{ticketID}` (e.g. `/admin/tickets/4`):

```json
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

**Replying to a ticket:** `POST /admin/tickets/{ticketID}/conversation` (e.g. `/admin/tickets/4/conversation`):

```json
{"message": "Please provide more details of your situation"}
```

Response:

```json
{
    "conversationID": 9
}
```

The message appears the next time the user fetches the ticket:

```json
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
            "message": "Please provide more details of your situation",
            "messageSentAt": "2024-02-14T16:33:40Z"
        }
    ]
}
```

Messages sent by admins are always attributed to `operator`.

Non-admin users hitting an admin route get:

```text
Access denied. Admin privilege required.
```
