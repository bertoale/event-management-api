package notification

import "strings"

// Helper functions untuk membuat notification request tanpa circular dependency

// NewCancellationRequest membuat request untuk notifikasi pembatalan event
func NewCancellationRequest(userID uint, eventID uint, message string) *CreateNotificationRequest {
	return &CreateNotificationRequest{
		UserID:  userID,
		EventID: &eventID,
		Type:    string(NotifCancellation),
		Message: message,
	}
}

// NewUpdateRequest membuat request untuk notifikasi update event
func NewUpdateRequest(userID uint, eventID uint, message string) *CreateNotificationRequest {
	return &CreateNotificationRequest{
		UserID:  userID,
		EventID: &eventID,
		Type:    string(NotifUpdate),
		Message: message,
	}
}

// NewReminderRequest membuat request untuk notifikasi reminder event
func NewReminderRequest(userID uint, eventID uint, message string) *CreateNotificationRequest {
	return &CreateNotificationRequest{
		UserID:  userID,
		EventID: &eventID,
		Type:    string(NotifReminder),
		Message: message,
	}
}

// ParseNotifType mengubah string menjadi NotifType yang valid
func ParseNotifType(s string) (NotifType, bool) {
	switch strings.ToLower(s) {
	case string(NotifReminder):
		return NotifReminder, true
	case string(NotifUpdate):
		return NotifUpdate, true
	case string(NotifCancellation):
		return NotifCancellation, true
	default:
		return "", false
	}
}
