package email

import (
	"fmt"
	"go-event/pkg/config"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type Service interface {
	SendEmail(to, toName, subject, htmlBody, textBody string) error
	SendWelcomeEmail(to, toName string) error
	SendReminderEmail(to, toName, eventTitle, eventDate string) error
	SendRegistrationConfirmationEmail(to, toName, eventTitle, eventDate, eventLocation string) error
	SendCancellationEmail(to, toName, eventTitle string) error
	SendUpdateEmail(to, toName, eventTitle, updateMessage string) error
}

type service struct {
	client *mailjet.Client
	cfg    *config.Config
}

// SendEmail implements Service.
func (s *service) SendEmail(to, toName, subject, htmlBody, textBody string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: s.cfg.MailSenderEmail,
				Name:  s.cfg.MailSenderName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
					Name:  toName,
				},
			},
			Subject:  subject,
			TextPart: textBody,
			HTMLPart: htmlBody,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := s.client.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendWelcomeEmail implements Service.
func (s *service) SendWelcomeEmail(to, toName string) error {
	subject := "🎉 Selamat Datang di GoEvent!"
	
	htmlBody := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff;">
					<div style="text-align: center; padding: 20px 0; background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); border-radius: 8px 8px 0 0;">
						<h1 style="color: #ffffff; margin: 0; font-size: 28px;">🎉 Selamat Datang!</h1>
					</div>
					<div style="padding: 30px; background-color: #f0fdf4; border-radius: 0 0 8px 8px;">
						<p style="font-size: 16px;">Halo <strong>%s</strong>,</p>
						<p style="font-size: 16px;">Terima kasih telah bergabung dengan <strong>GoEvent</strong>! 🎊</p>
						<div style="background-color: #ffffff; padding: 20px; border-left: 4px solid #10b981; border-radius: 4px; margin: 20px 0;">
							<p style="margin: 0; font-size: 16px;">
								Akun Anda telah berhasil dibuat. Sekarang Anda dapat:
							</p>
							<ul style="margin: 15px 0; padding-left: 20px; font-size: 16px;">
								<li>✨ Membuat dan mengelola event</li>
								<li>🎫 Mendaftar ke berbagai event menarik</li>
								<li>🔔 Mendapatkan notifikasi event</li>
								<li>📊 Melihat riwayat partisipasi Anda</li>
							</ul>
						</div>
						<p style="font-size: 16px;">Mulai jelajahi event yang tersedia dan ciptakan pengalaman tak terlupakan bersama kami!</p>
						<div style="text-align: center; margin-top: 30px;">
							<a href="http://localhost:3000" style="display: inline-block; padding: 12px 30px; background-color: #10b981; color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: bold;">Mulai Sekarang</a>
						</div>
					</div>
					<div style="text-align: center; padding: 20px; background-color: #f3f4f6; border-radius: 0 0 8px 8px;">
						<p style="font-size: 12px; color: #6b7280; margin: 0;">
							Email ini dikirim secara otomatis oleh <strong>GoEvent App</strong><br>
							Mohon tidak membalas email ini.
						</p>
					</div>
				</div>
			</body>
		</html>
	`, toName)

	textBody := fmt.Sprintf("🎉 Selamat Datang!\n\nHalo %s,\n\nTerima kasih telah bergabung dengan GoEvent! 🎊\n\nAkun Anda telah berhasil dibuat. Sekarang Anda dapat:\n- ✨ Membuat dan mengelola event\n- 🎫 Mendaftar ke berbagai event menarik\n- 🔔 Mendapatkan notifikasi event\n- 📊 Melihat riwayat partisipasi Anda\n\nMulai jelajahi event yang tersedia dan ciptakan pengalaman tak terlupakan bersama kami!\n\n---\nGoEvent App\nEmail ini dikirim secara otomatis. Mohon tidak membalas email ini.", 
		toName)

	return s.SendEmail(to, toName, subject, htmlBody, textBody)
}

// SendReminderEmail implements Service.
func (s *service) SendReminderEmail(to, toName, eventTitle, eventDate string) error {
	subject := fmt.Sprintf("Reminder: Event '%s' akan segera dimulai", eventTitle)
	
	htmlBody := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff;">
					<div style="text-align: center; padding: 20px 0; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); border-radius: 8px 8px 0 0;">
						<h1 style="color: #ffffff; margin: 0; font-size: 28px;">🔔 Event Reminder</h1>
					</div>
					<div style="padding: 30px; background-color: #f9fafb; border-radius: 0 0 8px 8px;">
						<p style="font-size: 16px;">Halo <strong>%s</strong>,</p>
						<p style="font-size: 16px;">Ini adalah pengingat bahwa event berikut akan segera dimulai:</p>
						<div style="background-color: #ffffff; padding: 20px; border-left: 4px solid #667eea; border-radius: 4px; margin: 20px 0;">
							<h2 style="color: #667eea; margin-top: 0; font-size: 22px;">%s</h2>
							<p style="margin: 10px 0; font-size: 16px;">
								<strong>📅 Waktu:</strong> %s
							</p>
						</div>
						<p style="font-size: 16px;">Pastikan Anda sudah siap dan jangan sampai terlewat!</p>
						<div style="text-align: center; margin-top: 30px;">
							<p style="font-size: 14px; color: #666;">Sampai jumpa di event! 👋</p>
						</div>
					</div>
					<div style="text-align: center; padding: 20px; background-color: #f3f4f6; border-radius: 0 0 8px 8px;">
						<p style="font-size: 12px; color: #6b7280; margin: 0;">
							Email ini dikirim secara otomatis oleh <strong>GoEvent App</strong><br>
							Mohon tidak membalas email ini.
						</p>
					</div>
				</div>
			</body>
		</html>
	`, toName, eventTitle, eventDate)

	textBody := fmt.Sprintf("🔔 Event Reminder\n\nHalo %s,\n\nIni adalah pengingat bahwa event '%s' akan segera dimulai pada %s.\n\nPastikan Anda sudah siap dan jangan sampai terlewat!\n\nSampai jumpa di event! 👋\n\n---\nGoEvent App\nEmail ini dikirim secara otomatis. Mohon tidak membalas email ini.", 
		toName, eventTitle, eventDate)

	return s.SendEmail(to, toName, subject, htmlBody, textBody)
}

// SendRegistrationConfirmationEmail implements Service.
func (s *service) SendRegistrationConfirmationEmail(to, toName, eventTitle, eventDate, eventLocation string) error {
	subject := fmt.Sprintf("✅ Konfirmasi Pendaftaran: %s", eventTitle)
	
	htmlBody := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff;">
					<div style="text-align: center; padding: 20px 0; background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); border-radius: 8px 8px 0 0;">
						<h1 style="color: #ffffff; margin: 0; font-size: 28px;">✅ Pendaftaran Berhasil!</h1>
					</div>
					<div style="padding: 30px; background-color: #f0fdf4; border-radius: 0 0 8px 8px;">
						<p style="font-size: 16px;">Halo <strong>%s</strong>,</p>
						<p style="font-size: 16px;">Selamat! Pendaftaran Anda untuk event berikut telah berhasil dikonfirmasi:</p>
						<div style="background-color: #ffffff; padding: 25px; border-left: 4px solid #10b981; border-radius: 4px; margin: 20px 0; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
							<h2 style="color: #059669; margin-top: 0; font-size: 22px;">%s</h2>
							<div style="margin-top: 15px;">
								<p style="margin: 8px 0; font-size: 16px;">
									<strong>📅 Waktu:</strong> %s
								</p>
								<p style="margin: 8px 0; font-size: 16px;">
									<strong>📍 Lokasi:</strong> %s
								</p>
							</div>
						</div>
						<div style="background-color: #d1fae5; padding: 15px; border-radius: 6px; margin: 20px 0;">
							<p style="margin: 0; font-size: 14px; color: #065f46;">
								<strong>💡 Tips:</strong> Simpan email ini sebagai referensi dan jangan lupa untuk hadir tepat waktu!
							</p>
						</div>
						<p style="font-size: 16px;">Kami akan mengirimkan pengingat menjelang event dimulai.</p>
						<div style="text-align: center; margin-top: 30px;">
							<p style="font-size: 14px; color: #666;">Sampai jumpa di event! 🎉</p>
						</div>
					</div>
					<div style="text-align: center; padding: 20px; background-color: #f3f4f6; border-radius: 0 0 8px 8px;">
						<p style="font-size: 12px; color: #6b7280; margin: 0;">
							Email ini dikirim secara otomatis oleh <strong>GoEvent App</strong><br>
							Mohon tidak membalas email ini.
						</p>
					</div>
				</div>
			</body>
		</html>
	`, toName, eventTitle, eventDate, eventLocation)

	textBody := fmt.Sprintf("✅ Pendaftaran Berhasil!\n\nHalo %s,\n\nSelamat! Pendaftaran Anda untuk event berikut telah berhasil dikonfirmasi:\n\n%s\n\n📅 Waktu: %s\n📍 Lokasi: %s\n\n💡 Tips: Simpan email ini sebagai referensi dan jangan lupa untuk hadir tepat waktu!\n\nKami akan mengirimkan pengingat menjelang event dimulai.\n\nSampai jumpa di event! 🎉\n\n---\nGoEvent App\nEmail ini dikirim secara otomatis. Mohon tidak membalas email ini.", 
		toName, eventTitle, eventDate, eventLocation)

	return s.SendEmail(to, toName, subject, htmlBody, textBody)
}

// SendCancellationEmail implements Service.
func (s *service) SendCancellationEmail(to, toName, eventTitle string) error {
	subject := fmt.Sprintf("❌ Pembatalan Event: %s", eventTitle)
	
	htmlBody := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff;">
					<div style="text-align: center; padding: 20px 0; background: linear-gradient(135deg, #ef4444 0%%, #dc2626 100%%); border-radius: 8px 8px 0 0;">
						<h1 style="color: #ffffff; margin: 0; font-size: 28px;">❌ Event Dibatalkan</h1>
					</div>
					<div style="padding: 30px; background-color: #fef2f2; border-radius: 0 0 8px 8px;">
						<p style="font-size: 16px;">Halo <strong>%s</strong>,</p>
						<p style="font-size: 16px;">Kami informasikan bahwa event berikut telah dibatalkan:</p>
						<div style="background-color: #ffffff; padding: 20px; border-left: 4px solid #ef4444; border-radius: 4px; margin: 20px 0;">
							<h2 style="color: #dc2626; margin-top: 0; font-size: 22px;">%s</h2>
						</div>
						<p style="font-size: 16px;">Mohon maaf atas ketidaknyamanannya. Kami akan memberitahu Anda jika ada update lebih lanjut atau event pengganti.</p>
						<div style="text-align: center; margin-top: 30px;">
							<p style="font-size: 14px; color: #666;">Terima kasih atas pengertian Anda 🙏</p>
						</div>
					</div>
					<div style="text-align: center; padding: 20px; background-color: #f3f4f6; border-radius: 0 0 8px 8px;">
						<p style="font-size: 12px; color: #6b7280; margin: 0;">
							Email ini dikirim secara otomatis oleh <strong>GoEvent App</strong><br>
							Mohon tidak membalas email ini.
						</p>
					</div>
				</div>
			</body>
		</html>
	`, toName, eventTitle)

	textBody := fmt.Sprintf("❌ Event Dibatalkan\n\nHalo %s,\n\nKami informasikan bahwa event '%s' telah dibatalkan.\n\nMohon maaf atas ketidaknyamanannya. Kami akan memberitahu Anda jika ada update lebih lanjut atau event pengganti.\n\nTerima kasih atas pengertian Anda 🙏\n\n---\nGoEvent App\nEmail ini dikirim secara otomatis. Mohon tidak membalas email ini.", 
		toName, eventTitle)

	return s.SendEmail(to, toName, subject, htmlBody, textBody)
}

// SendUpdateEmail implements Service.
func (s *service) SendUpdateEmail(to, toName, eventTitle, updateMessage string) error {
	subject := fmt.Sprintf("📢 Update Event: %s", eventTitle)
	
	htmlBody := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff;">
					<div style="text-align: center; padding: 20px 0; background: linear-gradient(135deg, #3b82f6 0%%, #2563eb 100%%); border-radius: 8px 8px 0 0;">
						<h1 style="color: #ffffff; margin: 0; font-size: 28px;">📢 Update Event</h1>
					</div>
					<div style="padding: 30px; background-color: #eff6ff; border-radius: 0 0 8px 8px;">
						<p style="font-size: 16px;">Halo <strong>%s</strong>,</p>
						<p style="font-size: 16px;">Ada update terbaru untuk event:</p>
						<div style="background-color: #ffffff; padding: 20px; border-left: 4px solid #3b82f6; border-radius: 4px; margin: 20px 0;">
							<h2 style="color: #2563eb; margin-top: 0; font-size: 22px;">%s</h2>
						</div>
						<div style="background-color: #dbeafe; padding: 20px; border-radius: 8px; margin: 20px 0;">
							<p style="margin: 0; font-size: 16px; color: #1e40af;"><strong>Informasi Update:</strong></p>
							<p style="margin: 10px 0 0 0; font-size: 16px;">%s</p>
						</div>
						<p style="font-size: 16px;">Terima kasih atas perhatiannya.</p>
						<div style="text-align: center; margin-top: 30px;">
							<p style="font-size: 14px; color: #666;">Tetap update dengan event Anda! ✨</p>
						</div>
					</div>
					<div style="text-align: center; padding: 20px; background-color: #f3f4f6; border-radius: 0 0 8px 8px;">
						<p style="font-size: 12px; color: #6b7280; margin: 0;">
							Email ini dikirim secara otomatis oleh <strong>GoEvent App</strong><br>
							Mohon tidak membalas email ini.
						</p>
					</div>
				</div>
			</body>
		</html>
	`, toName, eventTitle, updateMessage)

	textBody := fmt.Sprintf("📢 Update Event\n\nHalo %s,\n\nAda update terbaru untuk event '%s':\n\n%s\n\nTerima kasih atas perhatiannya.\n\nTetap update dengan event Anda! ✨\n\n---\nGoEvent App\nEmail ini dikirim secara otomatis. Mohon tidak membalas email ini.", 
		toName, eventTitle, updateMessage)

	return s.SendEmail(to, toName, subject, htmlBody, textBody)
}

func NewService(cfg *config.Config) Service {
	client := mailjet.NewMailjetClient(cfg.MailjetAPIKey, cfg.MailjetAPISecret)
	
	return &service{
		client: client,
		cfg:    cfg,
	}
}
