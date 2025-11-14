package notification

// Helper functions untuk membuat notification request tanpa circular dependency

// NewCancellationRequest membuat request untuk notifikasi pembatalan event
func NewCancellationRequest(userID uint, eventID uint, message string) *CreateNotificationRequest {
return &CreateNotificationRequest{
UserID:  userID,
EventID: &eventID,
Type:    NotifCancellation,
Message: message,
}
}

// NewUpdateRequest membuat request untuk notifikasi update event
func NewUpdateRequest(userID uint, eventID uint, message string) *CreateNotificationRequest {
return &CreateNotificationRequest{
UserID:  userID,
EventID: &eventID,
Type:    NotifUpdate,
Message: message,
}
}

// NewReminderRequest membuat request untuk notifikasi reminder event
func NewReminderRequest(userID uint, eventID uint, message string) *CreateNotificationRequest {
return &CreateNotificationRequest{
UserID:  userID,
EventID: &eventID,
Type:    NotifReminder,
Message: message,
}
}
