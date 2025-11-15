# Event Service API Documentation (Postman)

Berikut adalah dokumentasi endpoint Event untuk integrasi dengan Postman:

## 1. Create Event (Organizer Only)

- **Endpoint:** `/api/event/`
- **Method:** POST
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

```json
{
  "title": "string",
  "description": "string",
  "location": "string",
  "start_time": "2025-11-15T09:00:00Z",
  "end_time": "2025-11-15T12:00:00Z"
}
```

- **Response:**

```json
{
  "message": "event created successfully",
  "event": { ... }
}
```

## 2. Get All Events

- **Endpoint:** `/api/event/`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "events retrieved successfully",
  "events": [ ... ]
}
```

## 3. Get Event by ID

- **Endpoint:** `/api/event/{id}`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "event retrieved successfully",
  "event": { ... }
}
```

## 4. Update Event

- **Endpoint:** `/api/event/{id}`
- **Method:** PUT
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

```json
{
  "title": "string",
  "description": "string",
  "location": "string",
  "start_time": "2025-11-15T09:00:00Z",
  "end_time": "2025-11-15T12:00:00Z"
}
```

- **Response:**

```json
{
  "message": "event updated successfully",
  "event": { ... }
}
```

## 5. Delete Event

- **Endpoint:** `/api/event/{id}`
- **Method:** DELETE
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "event deleted successfully"
}
```

## 6. Get Events by Organizer

- **Endpoint:** `/api/event/organizer/{organizer_id}`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "events retrieved successfully",
  "events": [ ... ]
}
```

---

**Catatan:**

- Semua endpoint yang membutuhkan autentikasi harus mengirimkan header `Authorization: Bearer {jwt-token}`.
- Response `{ ... }` menyesuaikan dengan struktur event pada database.
- Untuk testing di Postman, pastikan JWT token valid dan role sesuai dengan endpoint yang diakses.
