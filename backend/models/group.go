package models

import "time"

// Group representa un grupo de apoyo
type Group struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Topic       string    `json:"topic"` // Tema del grupo (ej. ansiedad, duelo)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JoinGroupRequest es lo que la cara enviará cuando alguien quiera unirse a un grupo
type JoinGroupRequest struct {
	GroupID uint `json:"group_id"`
	// En un futuro, aquí iría la ID del usuario real, pero por ahora es anónimo.
}
