# Banking API Documentation

**Author**: Joaquin Gonzalez Alvarez  
**Date**: November 20, 2024

## Introduction

This document provides detailed documentation for a simple banking application, developed as part of a technical test for the company Valesa. The purpose of the application is to demonstrate backend development skills, particularly in creating a functional banking system using the Go programming language.

## Application Overview

The application is a basic banking system built in Go, utilizing the **Gin** framework to create a RESTful API. The system allows users to:

- Create and manage bank accounts.
- Perform financial transactions, including deposits and withdrawals.
- Transfer money between accounts.

Data is stored in memory using Go data structures, with no persistent database.

### Core Data Structures

The application uses two primary data structures:

- **Account**: Represents a bank account, including the account ID, owner's name, and balance.
- **Transaction**: Represents a financial transaction (either a deposit or withdrawal), linked to a specific account and amount.

### Controllers and Features

The API exposes several controllers to handle client requests:

- **Account Controller**: For creating, retrieving, and updating accounts.
- **Transaction Controller**: For performing transactions (deposit/withdrawal) on accounts.
- **Transfer Controller**: For transferring money between accounts.

The data access logic (in-memory storage) is modularized for better organization and code reuse.

### Error Handling

Common errors that are handled include:

- Insufficient balance when attempting a withdrawal.
- Invalid transaction amounts (e.g., negative or zero).
- Transfer attempts between the same account.
- Access attempts to non-existent accounts.

## Installation Guide

To install and run this application locally, follow the steps below:

### Prerequisites

Ensure the following are installed:

- **Go** version 1.18 or higher. [Download Go here](https://go.dev/dl/).
- **Git** for cloning the repository.

### Installation Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/iamximo/ValesaChallange.git
   ```

2. Navigate to the project directory:

   ```bash
   cd ValesaChallange
   ```

3. Install the required dependencies:

   ```bash
   go mod tidy
   ```

4. Build and run the server:

   ```bash
   go run main.go
   ```

The server will be available at `http://localhost:8080`.

### Running the Application

Once the application is running, it will be accessible at `http://localhost:8080`. You can interact with the API using tools like **Postman**, **curl**, or integrate it into your frontend.

## API Request Examples

Below are examples of `curl` commands for interacting with the Banking API.

### Create Account (POST /accounts)

```bash
curl -X POST http://localhost:8080/accounts \
     -H "Content-Type: application/json" \
     -d '{
           "owner": "JOQUIN",
           "initial_balance": 0
         }'
```

### Get Account by ID (GET /accounts/{id})

```bash
curl -X GET http://localhost:8080/accounts/990c804f-cb50-4994-8df3-2b1787bf87ae
```

### Get All Accounts (GET /accounts)

```bash
curl -X GET http://localhost:8080/accounts
```

### Create Transaction (POST /accounts/{id}/transactions)

```bash
curl -X POST http://localhost:8080/accounts/990c804f-cb50-4994-8df3-2b1787bf87ae/transactions \
     -H "Content-Type: application/json" \
     -d '{
           "type": "deposit",
           "amount": 6
         }'
```

### Get Transactions (GET /accounts/{id}/transactions)

```bash
curl -X GET http://localhost:8080/accounts/8e26b8b1-59e1-4ef7-8dc1-764f8eea55da/transactions
```

### Transfer Money (POST /transfer)

```bash
curl -X POST http://localhost:8080/transfer \
     -H "Content-Type: application/json" \
     -d '{
           "from_account_id": "8e26b8b1-59e1-4ef7-8dc1-764f8eea55da",
           "to_account_id": "67f67506-55d7-4d58-9e64-8280d8904a7d",
           "amount": 0.1
         }'
```

## API Endpoints

The application exposes several RESTful endpoints for interacting with the banking system. Below are the details of these endpoints along with example requests and responses.

### Request and Response Format

All requests and responses are in **JSON** format.

### Usage Examples

#### Create an Account

**Request:**

```http
POST /accounts
Content-Type: application/json

{
  "owner": "Joaquin",
  "initial_balance": 1000.0
}
```

**Response:**

```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": "f06b22bb-d9f0-42a3-9b9b-b234f9f329d9",
  "owner": "Joaquin",
  "balance": 1000.0
}
```

#### Get Account by ID

**Request:**

```http
GET /accounts/{id}
```

**Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": "f06b22bb-d9f0-42a3-9b9b-b234f9f329d9",
  "owner": "Joaquin",
  "balance": 1000.0
}
```

#### Perform a Transaction

**Request:**

```http
POST /accounts/{id}/transactions
Content-Type: application/json

{
  "type": "deposit",
  "amount": 500.0
}
```

**Response:**

```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": "acfe2d8f-fb83-455b-bce3-f93e2f4c02db",
  "accountId": "f06b22bb-d9f0-42a3-9b9b-b234f9f329d9",
  "type": "deposit",
  "amount": 500.0,
  "timestamp": "2024-11-20T10:32:00Z"
}
```

#### Transfer Between Accounts

**Request:**

```http
POST /transfer
Content-Type: application/json

{
  "from_account_id": "f06b22bb-d9f0-42a3-9b9b-b234f9f329d9",
  "to_account_id": "8ab21f7f-476e-4873-b032-9d183d0fb2c1",
  "amount": 200.0
}
```

**Response:**

```json
HTTP/1.1 201 Created
Content-Type: application/json

[
  {
    "id": "df82307c-d1bc-4bdb-99f3-e28f3b7f0638",
    "accountId": "f06b22bb-d9f0-42a3-9b9b-b234f9f329d9",
    "type": "withdrawal",
    "amount": 200.0,
    "timestamp": "2024-11-20T10:35:00Z"
  },
  {
    "id": "e11e9f7d-e5ab-4f8a-b60e-6744e334c121",
    "accountId": "8ab21f7f-476e-4873-b032-9d183d0fb2c1",
    "type": "deposit",
    "amount": 200.0,
    "timestamp": "2024-11-20T10:35:00Z"
  }
]
```

## Running Tests

To run the tests for the Banking API:

1. Navigate to the project directory:

   ```bash
   cd ValesaChallange
   ```
   
2. Run the tests with:

   ```bash
   go test ./...
   ```

### Resetting In-Memory Storage Between Tests

To reset the in-memory storage between tests, use the `Reset()` function from the `storage` package. This ensures that the tests do not interfere with each other.

## Future Improvements

To enhance the application, several improvements can be made:

- **Database Integration**: Introduce a database layer for persistent storage. This will separate business logic from data storage, improving maintainability and scalability.
- **Concurrency Handling**: Optimize concurrency management by leveraging database systems or ORM packages that handle concurrent transactions more efficiently.
- **Pagination**: Implement pagination for retrieving large datasets. This will improve performance and reduce memory consumption when dealing with a significant number of accounts or transactions.
```
