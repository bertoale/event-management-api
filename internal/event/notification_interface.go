package event

// NotificationService interface untuk menghindari circular dependency
// Method menerima string untuk type karena tidak bisa import NotifType dari notification package
type NotificationService interface {
	SendNotificationWithEmailByString(userID uint, eventID uint, notifTypeStr string, message, userEmail, userName string) error
}


