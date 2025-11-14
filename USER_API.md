# User Service API Documentation (Postman)

Berikut adalah dokumentasi endpoint User untuk integrasi dengan Postman:

## 1. Register User

- **Endpoint:** `/api/auth/register`
- **Method:** POST
- **Request Body:**

```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

- **Response:**

```json
{
  "message": "User registered successfully",
  "user": { ... }
}
```

## 2. Login User

- **Endpoint:** `/api/auth/login`
- **Method:** POST
- **Request Body:**

```json
{
  "email": "string",
  "password": "string"
}
```

- **Response:**

```json
{
  "message": "Login successful",
  "token": "jwt-token",
  "user": { ... }
}
```

## 3. Get Profile

- **Endpoint:** `/api/user/profile`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "profile retrieved successfully",
  "user": { ... }
}
```

## 4. Update Profile

- **Endpoint:** `/api/user/profile`
- **Method:** PUT
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

```json
{
  "name": "string",
  "email": "string"
}
```

- **Response:**

```json
{
  "message": "profile updated successfully",
  "user": { ... }
}
```

## 5. Change Password

- **Endpoint:** `/api/user/change-password`
- **Method:** POST
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

```json
{
  "old_password": "string",
  "new_password": "string"
}
```

- **Response:**

```json
{
  "message": "password changed successfully"
}
```

## 6. Get All Users (Admin only)

- **Endpoint:** `/api/user/`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "users retrieved successfully",
  "users": [ ... ]
}
```

## 7. Get User by ID (Admin only)

- **Endpoint:** `/api/user/{id}`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "user retrieved successfully",
  "user": { ... }
}
```

## 8. Delete User (Admin only)

- **Endpoint:** `/api/user/{id}`
- **Method:** DELETE
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "user deleted successfully"
}
```

## 9. Get Users by Role (Admin only)

- **Endpoint:** `/api/user/role/{role}`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "users retrieved successfully",
  "users": [ ... ]
}
```

---

**Catatan:**

- Semua endpoint yang membutuhkan autentikasi harus mengirimkan header `Authorization: Bearer {jwt-token}`.
- Response `{ ... }` menyesuaikan dengan struktur user pada database.
- Untuk testing di Postman, pastikan JWT token valid dan role sesuai dengan endpoint yang diakses.
