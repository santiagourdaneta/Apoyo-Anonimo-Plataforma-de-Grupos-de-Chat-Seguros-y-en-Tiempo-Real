package models

import "time"

// Message representa un mensaje en un grupo
type Message struct {
	ID        uint      `json:"id"`
	GroupID   uint      `json:"group_id"`
	// UserID    uint      `json:"user_id"` // No lo usaremos por ahora para mantener el anonimato real
	Username  string    `json:"username"`  // Un nombre de usuario anónimo (ej. "UsuarioAnónimo123")
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	IsModerated bool `json:"is_moderated"` // Si el mensaje fue marcado por moderación
}

// SendMessageRequest es lo que la cara enviará cuando alguien escriba un mensaje
type SendMessageRequest struct {
	GroupID  uint   `json:"group_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}
