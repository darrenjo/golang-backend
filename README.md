# Backend API Project

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)
- [License](#license)

# Introduction

This is a backend API project using Golang for managing user profiles and photos. The API allows users to register, login, upload, and manage their profile photos. It uses JWT for authentication and MySQL for data storage.

# Features

- User registration and login
- JWT authentication
- Upload and manage profile photos
- Replace existing profile photos
- Delete profile photos
- Retrieve profile photos

# Technologies Used

- Go
- Gorilla Mux
- JWT
- MySQL

# Installation

1. **Create a new database**

2. **Update environment variables:**

   - Create a `.env` file in the root directory of the project.
   - Add your database credentials and other necessary configurations to the .env file.

3. **Install dependencies:**
   ```
   go mod tidy
   ```

# Running the Application

1. **Run the application:**
   ```
   go run main.go
   ```

# API Endpoints

### User Endpoints

- Register:
  - URL: `/users/register`
  - Method: `POST`
  - Body:
    ```json
    {
      "username": "testuser",
      "email": "testuser@example.com",
      "password": "password123"
    }
    ```
- Login:
  - URL: `/users/login`
  - Method: `POST`
  - Body:
    ```json
    {
      "email": "testuser@example.com",
      "password": "password123"
    }
    ```

### Photo Endpoints

- Set Profile Photo:
  - URL: `/api/photos/profile`
  - Method: `POST`
  - Headers: `Authorization: Bearer <token>`
  - Body:
    ```json
    {
      "photo_url": "http://example.com/photo.jpg"
    }
    ```
- Get Profile Photo:

  - URL: `/api/photos/profile`
  - Method: `GET`
  - Headers: `Authorization: Bearer <token>`

- Delete Profile Photo:
  - URL: `/api/photos/{photoId}`
  - Method: `DELETE`
  - Headers: `Authorization: Bearer <token>`

# Environment Variables

Create a `.env` file in the root directory and add the following variables:

```
DB_USER=YOUR_DB_USERNAME
DB_PASSWORD=YOUR_DB_PASSWORD
DB_HOST=YOUR_DB_HOST
DB_PORT=YOUR_DB_PORT
DB_NAME=YOUR_DB_NAME
JWT_SECRET=YOUR_JWT_SECRET_KEY
PORT=YOUR_DESIRED_PORT
```

In my case:

```
DB_USER=root
DB_PASSWORD=
DB_HOST=localhost
DB_PORT=3306
DB_NAME=usergolang
JWT_SECRET=5gpuoMSGH7r7X-kpGy9SbRcnyygLpicbX0_J5goqyrI
PORT=3000
```

# Test Case

### User Endpoints

- Register

  **Endpoint:** `POST` `/users/register`

  **Description:** Test user registration functionality.

  #### Test Case 1: Successful Registration

  - **Request:**

    ```json
    {
      "username": "testuser",
      "email": "testuser@example.com",
      "password": "password123"
    }
    ```

  - **Expected Response:**
    - Status Code: 201 Created
    - Response:
    ```json
    {
      "id": 1,
      "username": "testuser",
      "email": "testuser@example.com",
      "password": "$2a$14$Rl7C1SsyWrF3yrOMuKmJN.PXJDd2nG0Cln6lrfUPZiy58a.bh3OMC",
      "photos": null,
      "created_at": "2024-05-31T09:52:18.6555727+07:00",
      "updated_at": "2024-05-31T09:52:18.6555727+07:00"
    }
    ```

  #### Test Case 2: Registration with Existing Email

  - **Request:**
    ```json
    {
      "username": "testuser2",
      "email": "testuser@example.com",
      "password": "password123"
    }
    ```
  - Expected Response:
    - Status Code: 400 Bad Request
    - Response:
    ```json
    {
      "error": "Error saving user"
    }
    ```

  #### Test Case 3: Registration with Existing Username

  - **Request:**
    ```json
    {
      "username": "testuser",
      "email": "testuser2@example.com",
      "password": "password123"
    }
    ```
  - Expected Response:
    - Status Code: 400 Bad Request
    - Response:
    ```json
    {
      "error": "Error saving user"
    }
    ```

- Login

  **Endpoint:** `POST` `/users/login`

  **Description:** Test user login functionality.

  #### Test Case 1: Successful Login

  - **Request:**

    ```json
    {
      "email": "testuser@example.com",
      "password": "password123"
    }
    ```

  - **Expected Response:**
    - Status Code: 200 OK
    - Response:
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTczODM1NDQsInVzZXJfaWQiOjF9.BHCTYIA3y5G-36jK7W-7oAcnfCHmyr9vbBbCvCySvHM"
    }
    ```

  #### Test Case 2: Invalid Credentials

  - **Request:**

    ```json
    {
      "email": "testuser@example.com",
      "password": "asdsdpassword123"
    }
    ```

  - **Expected Response:**
    - Status Code: 401 Unauthorized
    - Response:
    ```json
    {
      "error": "Invalid email or password"
    }
    ```

---

### Photo Endpoints

- Upload Profile Photo

  **Endpoint:** `POST` `/api/photos/profile`

  **Description:** Test profile photo upload functionality.

  #### Test Case 1: Successful Upload

  - **Request:**

    ```json
    {
      "photo_url": "http://example.com/testuser.jpg"
    }
    ```

    - **Headers:**
      ```
      Authorization: Bearer user_jwt_token
      ```

  - **Expected Response:**
    - Status Code: 200 OK
    - Response:
    ```json
    {
      "id": 1,
      "photo_url": "http://example.com/testuser.jpg",
      "user_id": 1,
      "is_profile": true,
      "created_at": "2024-05-30T13:53:32Z",
      "updated_at": "2024-05-31T10:05:34.6880182+07:00"
    }
    ```

  #### Test Case 2: Replace Profile Photo

  - **Request:**

    ```json
    {
      "photo_url": "http://example.com/testuser99999999.jpg"
    }
    ```

    - **Headers:**
      ```
      Authorization: Bearer user_jwt_token
      ```

  - **Expected Response:**
    - Status Code: 200 OK
    - Response:
    ```json
    {
      "id": 2,
      "photo_url": "http://example.com/testuser99999999.jpg",
      "user_id": 1,
      "is_profile": true,
      "created_at": "2024-05-30T13:53:32Z",
      "updated_at": "2024-05-31T10:05:34.6880182+07:00"
    }
    ```

- Retrieve Profile Photo

  **Endpoint:** `GET` `/api/photos/profile`

  **Description:** Test retrieving the profile photo URL.

  #### Test Case 1: Successful Retrieval

  - **Request:**

    - **Headers:**
      ```
      Authorization: Bearer user_jwt_token
      ```

  - **Expected Response:**
    - Status Code: 200 OK
    - Response:
      ```json
      {
        "photo_url": "http://example.com/testuser99999999.jpg"
      }
      ```

- Delete Profile Photo

  **Endpoint:** `DELETE` `/api/photos/{photoId}`

  **Description:** Test profile photo deletion functionality.

  #### Test Case 1: Successful Deletion

  - **Request:**
    - **Headers:**
      ```
      Authorization: Bearer user_jwt_token
      ```
  - **Expected Response:**
    - Status Code: 200 OK
    - Response:
      ```json
      {
        "result": "success"
      }
      ```

  #### Test Case 2: Delete Other User Photo

  - **Request:**
    - **Headers:**
      ```
      Authorization: Bearer user_jwt_token
      ```
  - **Expected Response:**
    - Status Code: 404 NOT FOUND
    - Response:
      ```json
      {
        "error": "Photo not found or user does not have access"
      }
      ```
