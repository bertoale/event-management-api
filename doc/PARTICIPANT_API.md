# Participant Service API Documentation (Postman)

Berikut adalah dokumentasi endpoint Participant untuk integrasi dengan Postman:

## 1. Register as Participant

- **Endpoint:** `/api/participants/{id}`
- **Method:** POST
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Request Body:**

- **Response:**

```json
{
  "message": "participant registered successfully",
  "participant": { ... }
}
```

## 2. Cancel Participation

- **Endpoint:** `/api/participants/{id}`
- **Method:** DELETE
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "participant cancelled successfully"
}
```

## 3. Get Participants by Event

- **Endpoint:** `//api/participants/{id}`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "participants retrieved successfully",
  "participants": [ ... ]
}
```

## 4. Get My Participation

- **Endpoint:** `/api/participant/my`
- **Method:** GET
- **Headers:**
  - Authorization: Bearer {jwt-token}
- **Response:**

```json
{
  "message": "participation retrieved successfully",
  "participant": { ... }
}
```

---

**Catatan:**

- Semua endpoint yang membutuhkan autentikasi harus mengirimkan header `Authorization: Bearer {jwt-token}`.
- Response `{ ... }` menyesuaikan dengan struktur participant pada database.
- Untuk testing di Postman, pastikan JWT token valid dan user memiliki hak akses yang sesuai.
