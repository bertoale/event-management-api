# Schedule Service API Documentation (Postman)

Berikut adalah dokumentasi endpoint Schedule untuk integrasi dengan Postman:

## 1. Create Schedule

- **Endpoint:** `/api/event/{id}/schedule`
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

- **Endpoint:** `/api/event/{id}/schedule`
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

## 3. Get Schedule by ID

- **Endpoint:** `/api/event/{id}/schedule`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "schedule retrieved successfully",
  "schedule": { ... }
}
```

## 4. Update Schedule

- **Endpoint:** `/api/schedule/{id}`
- **Method:** PUT
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

```json
{
  "reminder_time": "2025-11-15T08:00:00Z"
}
```

- **Response:**

```json
{
  "message": "schedule updated successfully",
  "schedule": { ... }
}
```

## 5. Delete Schedule

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
