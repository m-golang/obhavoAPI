# Ob-havo API Service

A robust and efficient weather data service that allows users to fetch weather information for different locations. This service also supports authentication via API keys, caching weather data in Redis, and managing user accounts.

## Table of Contents

- [Overview](#overview)
- [Technologies Used](#technologies-used)
- [Setup](#setup)
- [API Endpoints](#api-endpoints)
  - [User Registration](#user-registration)
  - [User Authentication](#user-authentication)
  - [User Dashboard](#user-dashboard)
  - [User Logout](#User-Logut)
  - [Fetch Weather Data](#fetch-weather-data)
  - [Fetch Bulk Weather Data](#fetch-bulk-weather-data)
- [Error Handling](#error-handling)
- [Redis Cache](#redis-cache)
- [Cron Job](#cron-job-for-periodic-cache-updates)

## Overview

The Weather API Service provides weather data retrieval capabilities, including support for multiple locations and bulk queries. It also supports user management with authentication via username and password, and generates an API key for authorized access. All weather data is cached in Redis to improve response times and minimize external API calls.

## Technologies Used

- Go (Golang)
- Redis
- Weather API (via [WeatherAPI.com](https://www.weatherapi.com/))
- MySQL (or other database for user management)
- JWT-based authentication with secure cookies.
- bcrypt (for password hashing)
- HTTP(S) API communication

## Setup

### Prerequisites

1. Go - Install Go version 1.18 or above. You can download it from the [official Go website](https://go.dev/).
2. Redis - Ensure you have Redis running locally or via Docker.
3. Weather API Key - Sign up at [WeatherAPI.com](https://www.weatherapi.com/) and get an API key.
4. MySQL (or other database) - Set up a MySQL database for user management.

### Steps to Run Locally

1. Clone this repository:

   ```bash
   git clone https://github.com/m-golang/obhavoAPI
   cd obhavoAPI

   ```

2. Create a `.env` file with the following environment variables:

   ```bash
   DB_USER_NAME=your-db-username
   DB_USER_PASSWORD=your-db-password
   DB_NAME=your-db-name
   JWT_SECRET_KEY=your-secret_key-for-JWT
   API_KEY_FOR_WEATHERAPI=your-weatherapi-com-api-key
   REDIS_ADDR=localhost:6379
   REDIS_PASS=your-redis-password

   ```

3. Start the application:
   ```bash
   go mod tidy
   go run ./cmd/...
   ```

## API Endpoints

1. ### User Registration

- **Endpoint:** `POST /api/v1/signup`
- **Description:** Registers a new user with the given details.
- **Request Body:**
  ``bash
  {
    "name": "John",
    "surname": "Doe",
    "username": "johndoe",
    "password": "password123"
  }

- **Response:**

  ```bash
  {
    "message": "User registered successfully."
  }

  ```

- **Errors:**
  - `400 Bad Request` - Missing or invalid data.
  - `409 Conflict` - Username already exists.

2. ### User Authentication

   - **Endpoint:** `POST /api/v1/login`
   - **Description:** Authenticates a user and returns message.
   - **Request Body:**

   ```bash
    {
      "username": "johndoe",
      "password": "password123"
    }
   ```

   - **Response:**

   ```bash
   {
     "message": "Login complete! Explore what's new!"
   }
   ```

   - **Errors:**
   - `401 Unauthorized` - Invalid credentials.
   - `404 Not Found` - User not found.

3. ### User Dashboard
   - **Endpoint:** `GET /api/v1/user/dashboard`
   - **Description:** Authenticated user gets API key.
   - **Response:**

   ```bash
   {
     "your API key": {your-API-key},
   }
   ```
4. ### User Logut
   - **Endpoint:** `GET /api/v1/logout`
   - **Description:** User logout.
   - **Response:**

   ```bash
   {
    "message": "You are now logged out. Have a great day!"
   }
   ```
5. ### Fetch Weather Data

   - **Call:** `GET localhost:8080/api/v1/weather.current?key={your-api-key}&q={location}`
   - **Description:** Fetches weather data for a specific location.
   - **Query Parameters:**
     - q (required): Location name (e.g., "Tashkent").
   - **Response:**

   ```bash
   {
       "location": {
           "name": "Tashkent",
           "country": "Uzbekistan",
           "lat": 34.517,
           "lon": 69.183,
           "temp_c": -2.1,
           "temp_color": "#B3DFFD",
           "wind_kph": 7.6,
           "wind_color": "#E0F7FA",
           "cloud": 5,
           "cloud_color": "#FFF9C4"
       }
   }
   ```

   - **Errors:**
   - `404 Not Found` - Location not found.
   - `500 Internal` Server Error - Error fetching data.

6. ### Fetch Bulk Weather Data

   - **Call:** `POST localhost:8080/api/v1/weather.current?key={your-api-key}&q=bulk`
   - **Description:** Fetches weather data for multiple locations.
   - **Bulk Request Example:**

   ```bash
   curl --location --request POST 'localhost:8080/api/v1/weather.current?key={your-api-key}&q=bulk' \
   --header 'Content-Type: application/json' \
   --data '{
             "locations": [
                           {
                             "q": "new york"
                           },
                           {
                             "q": "london"
                           },
                           {
                             "q": "tashkent"
                           }
                           ]
           }'
   ```

   - **Response:**

   ```bash
   {
     "bulk": [
               {
                 "name": "New York",
                 "country": "United States of America",
                 "lat": 40.7142,
                 "lon": -74.0064,
                 "temp_c": 1.7,
                 "temp_color": "#E6F7FF",
                 "wind_kph": 15.1,
                 "wind_color": "#B2EBF2",
                 "cloud": 75,
                 "cloud_color": "#9E9E9E"
               },
               {
                 "name": "London",
                 "country": "United Kingdom",
                 "lat": 51.5171,
                 "lon": -0.1062,
                 "temp_c": 3.1,
                 "temp_color": "#E6F7FF",
                 "wind_kph": 4.7,
                 "wind_color": "#E0F7FA",
                 "cloud": 100,
                 "cloud_color": "#616161"
               },
               {
                 "name": "Tashkent",
                 "country": "Uzbekistan",
                 "lat": 41.3167,
                 "lon": 69.25,
                 "temp_c": 1.3,
                 "temp_color": "#E6F7FF",
                 "wind_kph": 3.6,
                 "wind_color": "#E0F7FA",
                 "cloud": 100,
                 "cloud_color": "#616161"
               }
             ]
   }
   ```

   - **Errors:**
   - `404 Not Found` - One or more locations not found.

   - **Request Body:**

   ```bash
   curl --location --request POST 'localhost:8080/api/v1/weather.current?key={your-api-key}&q=bulk' \
   --header 'Content-Type: application/json' \
   --data '{
             "locations": [
                             {
                               "q": "new york"
                             },
                             {
                               "q": "london"
                             },
                             {
                               "q": "tashkent"
                             },
                             {
                               "q": "locationNotFound"
                             }
                           ]
           }
   ```

   - **Response:**

   ```bash
   {
     "bulk": [
               ...
               {
               "name": "Tashkent",
               "country": "Uzbekistan",
               "lat": 41.3167,
               "lon": 69.25,
               "temp_c": 2.3,
               "temp_color": "#E6F7FF",
               "wind_kph": 3.6,
               "wind_color": "#E0F7FA",
               "cloud": 100,
               "cloud_color": "#616161"
               }
               ],
               "not_found": [
               "'locationNotFound' not found"
             ]
   }
   ```

## Error Handling

The API follows RESTful conventions for error handling. Some common error responses include: - **400 Bad Request** - Invalid or missing input data. - **401 Unauthorized** - Invalid authentication or API key. - **404 Not Found - Requested** resource (e.g., location) not found. - **500 Internal Server Error** - Unexpected server errors.

## Redis Cache

### Weather Data Caching

Weather data for locations is cached in Redis to improve performance and reduce unnecessary API calls. The cache stores the latest weather data for a location for up to 30 minutes. After 30 minutes, the cached data expires, and a new request is made to the weather API to refresh the data.

### Cron Job for Periodic Cache Updates

A cron job is set up to automatically refresh the weather data cache at regular intervals. This helps ensure that cached data is up-to-date, even if no new requests are made.

### Cron Job Details:
- **Job Frequency:** Every 30 minutes.
- **Job Function:** The cron job fetches weather data for a predefined list of locations (e.g., major cities or countries) and updates the Redis cache.
- **Purpose:** To keep the cache updated periodically and minimize delays for users accessing weather data, ensuring that they always get the latest information.
