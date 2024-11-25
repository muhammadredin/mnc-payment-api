# Payment API

Welcome to the **Payment API**! This API allows you to handle user authentication, customer information, and transaction management. To get started, ensure you have the `.env` file configured with the necessary environment variables.

## Configuration

### .env File

Ensure your `.env` file contains the following parameters:

```txt
APPLICATION_NAME=My Payment API
SERVER_PORT=8081
LOGIN_EXPIRATION_DURATION=10
JWT_SIGNATURE_KEY=my-super-secret-key
```

## Features

### Authentication

#### 1. **Register** - `/api/public/auth/register`

Create a new customer account.

- **Request Body Example**:

    ```json
    {
        "username": "johndoe",
        "password": "password"
    }
    ```

- **Response Body Example**:

    ```json
    {
        "status_code": 201,
        "message": "Successfully created a new customer",
        "data": []
    }
    ```

#### 2. **Login** - `/api/public/auth/login`

Authenticate a user and obtain access and refresh tokens.

- **Request Body Example**:

    ```json
    {
        "username": "johndoe",
        "password": "password"
    }
    ```

- **Response Body Example**:

    ```json
    {
        "status_code": 200,
        "message": "Successfully logged in",
        "data": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "d4fc3603-43fb-4b1c-b477-4759f196e070",
            "customer_id": "ef749bd6-4f11-404b-b00f-caa5eb5e81d8"
        }
    }
    ```

#### 3. **Logout** - `/api/public/auth/logout`

Logout the authenticated user. A valid refresh token is required in the cookie.

- **Response Body Example**:

    ```json
    {
        "status_code": 200,
        "message": "Successfully logged out",
        "data": []
    }
    ```

#### 4. **Refresh Token** - `/api/public/auth/refresh-token`

Obtain a new access token using a valid refresh token in the cookie.

- **Response Body Example**:

    ```json
    {
        "status_code": 200,
        "message": "Successfully refreshed token",
        "data": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "d4fc3603-43fb-4b1c-b477-4759f196e070",
            "customer_id": "ef749bd6-4f11-404b-b00f-caa5eb5e81d8"
        }
    }
    ```

---

### Customer

#### 5. **Get Customer by Id** - `/api/customers/{id}`

Retrieve the details of a customer by their unique ID. The user must be authenticated as the current customer.

- **Response Body Example**:

    ```json
    {
        "status_code": 200,
        "message": "Successfully retrieved customer data",
        "data": {
            "id": "ef749bd6-4f11-404b-b00f-caa5eb5e81d8",
            "username": "johndoe",
            "wallet_id": "1070f292-5d68-4b30-b37f-32042675ef2a",
            "balance": 945000
        }
    }
    ```

---

### Transaction

#### 6. **Create Transaction** - `/api/transactions`

Create a transaction between two wallets. The user must be authenticated.

- **Request Body Example**:

    ```json
    {
        "from_wallet_id": "1070f292-5d68-4b30-b37f-32042675ef2a",
        "to_wallet_id": "198a1bff-50a7-4a3f-a18c-a724dea104de",
        "amount": 500000,
        "message": "Salary"
    }
    ```

- **Response Body Example**:

    ```json
    {
        "status_code": 201,
        "message": "Successfully created a transaction",
        "data": {
            "id": "6453869a-01b6-49a6-bbff-8102023c2622",
            "from_wallet_id": "1070f292-5d68-4b30-b37f-32042675ef2a",
            "to_wallet_id": "198a1bff-50a7-4a3f-a18c-a724dea104de",
            "created_at": "2024-11-25T23:09:29+07:00",
            "amount": 5000,
            "message": "Salary"
        }
    }
    ```

---

## Additional Notes

- All API requests that involve customer data require authentication via access tokens.
- JWT tokens (access & refresh) are used for user authentication and session management.
- Make sure your `.env` file is correctly configured before running the API.
- Ensure that you have a valid refresh token in the cookie for certain authentication routes like `logout` and `refresh-token`.
- To create new transaction user must have sufficient balance otherwise will return an error.
- For full testing purpose we created JSON file for API with request body with predefined user that has enough amount of balance to do transaction in test folder.
---

## Getting Started

1. Clone this repository.
2. Create a `.env` file in the root of your project with the required configuration (see above).
3. Run the application using the following command:

    ```bash
    go run main.go
    ```

4. Access the API at `http://localhost:8081`.