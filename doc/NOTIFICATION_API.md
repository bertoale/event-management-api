# Notification Service API Documentation (Postman)

Berikut adalah dokumentasi endpoint Notification untuk integrasi dengan Postman:

## 1. Get My Notifications

- **Endpoint:** `/api/notification/`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "notifications retrieved successfully",
  "notifications": [ ... ]
}
```

## 2. Mark Notification as Read

- **Endpoint:** `/api/notification/{id}/read`
- **Method:** PUT
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "notification marked as read successfully"
}
```

## 3. Delete Notification

- **Endpoint:** `/api/notification/{id}`
- **Method:** DELETE
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "notification deleted successfully"
}
```

## 4. Create Notification (Admin Only)

- **Endpoint:** `/api/notification/`
- **Method:** POST
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

```json
{
  "user_id": 2,
  "type": "reminder|update|cancellation",
  "message": "string"
}
```

- **Response:**

```json
{
  "message": "notification created successfully",
  "notification": { ... }
}
```

---

**Catatan:**

- Semua endpoint yang membutuhkan autentikasi harus mengirimkan header `Authorization: Bearer {jwt-token}`.
- Endpoint POST hanya bisa diakses oleh admin.
- Response `{ ... }` menyesuaikan dengan struktur notification pada database.
- Untuk testing di Postman, pastikan JWT token valid dan user memiliki hak akses yang sesuai.
