# GoEvent - Event Management Backend

GoEvent is a backend event management system built with Go, designed to support modern event operations with automated scheduling, notifications, and professional email integration.

## Features

- **Event Scheduling & Automation**: Schedule events, reminders, and automatic end events using gocron.
- **Notification System**: Automatic notifications for participants and admins (reminders, confirmations, cancellations, event updates).
- **Email Integration**: Mailjet integration for sending welcome, reminder, confirmation, cancellation, and event update emails.
- **Modular Architecture**: Clean structure with dependency injection, adapters, and interfaces to avoid circular imports.
- **Comprehensive Documentation & Testing**: Technical documentation and testing guides included in this README.

## Project Structure

- `cmd/` : Application entry point (main.go)
- `internal/` : Main modules (event, participant, schedule, notification, email)
- `pkg/` : Configuration, database, middleware
- `config/` : Configuration files

## Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <repo-url>
   cd Goevent
   ```
2. **Install dependencies**
   ```bash
   go mod tidy
   ```
3. **Configure Mailjet and database**
   - Edit configuration files in `config/` and set your Mailjet credentials and database connection.
4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

## Main API Endpoints

- `/api/auth/` : Auth
- `/api/events` : Event management
- `/api/participants` : Participant registration & management
- `/api/schedule` : Event scheduling
- `/api/notification` : Notifications & email

## API Documentation

### Event Endpoints (Organizer)

| Endpoint          | Method | Auth Required | Role      | Description                      |
| ----------------- | ------ | ------------- | --------- | -------------------------------- |
| `/api/events/`    | POST   | Yes           | Organizer | Create a new event               |
| `/api/events/`    | GET    | Yes           | Organizer | Get all events by organizer/user |
| `/api/events/:id` | PUT    | Yes           | Organizer | Update event by ID               |
| `/api/events/:id` | DELETE | Yes           | Organizer | Delete event by ID               |

All endpoints require authentication (JWT) and organizer role.

### Participant Endpoints

| Endpoint                | Method | Auth Required | Role      | Description                    |
| ----------------------- | ------ | ------------- | --------- | ------------------------------ |
| `/api/participants/`    | POST   | Yes           | User      | Register as participant        |
| `/api/participants/`    | GET    | Yes           | Organizer | Get all participants for event |
| `/api/participants/:id` | PUT    | Yes           | User      | Update participant info        |
| `/api/participants/:id` | DELETE | Yes           | User      | Cancel participation           |

### Schedule Endpoints

| Endpoint            | Method | Auth Required | Role      | Description           |
| ------------------- | ------ | ------------- | --------- | --------------------- |
| `/api/schedule/`    | POST   | Yes           | Organizer | Create event schedule |
| `/api/schedule/`    | GET    | Yes           | Organizer | Get all schedules     |
| `/api/schedule/:id` | PUT    | Yes           | Organizer | Update schedule by ID |
| `/api/schedule/:id` | DELETE | Yes           | Organizer | Delete schedule by ID |

### Notification & Email Endpoints

| Endpoint             | Method | Auth Required | Role      | Description             |
| -------------------- | ------ | ------------- | --------- | ----------------------- |
| `/api/notification/` | POST   | Yes           | Organizer | Send notification/email |
| `/api/notification/` | GET    | Yes           | Organizer | Get all notifications   |

## Email Integration (Mailjet)

- Automatic email delivery for welcome, reminders, confirmations, cancellations, and event updates.
- Configuration and implementation in `internal/notification/email/service.go`.
- To test email integration, ensure Mailjet credentials are set in config and trigger relevant endpoints.

## Scheduler (gocron)

- Automated job scheduling for reminders and event endings.
- Implementation in `internal/schedule/scheduler.go`.
- To test scheduler, create events and schedules, and verify automatic reminders and event status updates.

## Troubleshooting & Testing

- Circular imports are resolved using adapters and interfaces.
- For email and scheduler testing, use the relevant endpoints and check logs/output.
- Regular build and test recommended.

## Portfolio Highlights

- Clean, modular, and scalable architecture.
- Professional email integration (Mailjet).
- Automated scheduler (gocron).
- Comprehensive documentation for backend engineering portfolio.

---
