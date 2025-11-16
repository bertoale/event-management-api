# Schedule Service API Documentation (Postman)

Berikut adalah dokumentasi endpoint Schedule untuk integrasi dengan Postman:

## 1. Create Schedule

- **Endpoint:** `/api/schedule/event/{eventId}
- **Method:** POST
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

```json
{
  "job_type": "string",
  "run_at": "2025-11-15T10:00:00Z"
}
```

- **Response:**

```json
{
  "message": "schedule created successfully",
  "schedule": { ... }
}
```

## 2. Get All Schedules

- **Endpoint:** `/api/schedule/event/{eventId}`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "schedules retrieved successfully",
  "schedules": [ ... ]
}
```

## 3. Delete Schedule

- **Endpoint:** `/api/schedule/{id}`
- **Method:** DELETE
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "schedule deleted successfully"
}
```

---

**Catatan:**

- Semua endpoint yang membutuhkan autentikasi harus mengirimkan header `Authorization: Bearer {jwt-token}`.
- Response `{ ... }` menyesuaikan dengan struktur schedule pada database.
- Untuk testing di Postman, pastikan JWT token valid dan user memiliki hak akses yang sesuai.
