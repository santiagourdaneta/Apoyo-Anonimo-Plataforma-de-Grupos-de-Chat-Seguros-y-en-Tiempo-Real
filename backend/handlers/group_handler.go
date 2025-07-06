package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"log"

	"apoyo-anonimo/backend/config"
	"apoyo-anonimo/backend/models"
	"github.com/gin-gonic/gin"
	"strings" //
)

// CreateGroup crea un nuevo grupo de apoyo
func CreateGroup(c *gin.Context) {
	var newGroup models.Group
	if err := c.ShouldBindJSON(&newGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de grupo inválidos", "details": err.Error()})
		return
	}

	// Validaciones básicas (¡importante para la seguridad y calidad de datos!)
	if newGroup.Name == "" || len(newGroup.Name) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre del grupo es obligatorio y debe tener al menos 3 letras."})
		return
	}
	if newGroup.Topic == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El tema del grupo es obligatorio."})
		return
	}
	// Puedes añadir más validaciones aquí (ej. longitud máxima, caracteres permitidos)

	result, err := config.DB.Exec(
		"INSERT INTO groups (name, description, topic) VALUES (?, ?, ?)",
		newGroup.Name, newGroup.Description, newGroup.Topic,
	)
	if err != nil {
		log.Printf("Error al insertar grupo en DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el grupo", "details": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error al obtener ID del nuevo grupo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el grupo"})
		return
	}
	newGroup.ID = uint(id)
	c.JSON(http.StatusCreated, newGroup)
}

// GetGroups obtiene todos los grupos disponibles
func GetGroups(c *gin.Context) {
	// Puedes añadir filtros aquí (ej. por tema, por nombre)
	topic := c.Query("topic")
	search := c.Query("search")

	log.Printf("DEBUG GRUPO: Solicitud GetGroups - Tema: '%s', Búsqueda: '%s'", topic, search) // <--- NUEVO LOG

	query := "SELECT id, name, description, topic, created_at, updated_at FROM groups"
	args := []interface{}{}
	whereClauses := []string{}

	if topic != "" {
		whereClauses = append(whereClauses, "topic = ?")
		args = append(args, topic)
	}
	if search != "" {
		whereClauses = append(whereClauses, "name LIKE ? OR description LIKE ?")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ") // Necesitarás importar "strings"
	}

	query += " ORDER BY created_at DESC" // Ordenar por los más nuevos

	log.Printf("DEBUG GRUPO: Ejecutando consulta SQL: '%s' con args: %v", query, args) // <--- NUEVO LOG

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Printf("ERROR GRUPO: Error al ejecutar la consulta: %v", err) // <--- LOG MEJORADO
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener grupos", "details": err.Error()})
		return
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.Topic, &g.CreatedAt, &g.UpdatedAt); err != nil {
			log.Printf("ERROR GRUPO: Error al escanear fila de grupo: %v", err) // <--- LOG MEJORADO
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar grupos", "details": err.Error()})
			return
		}
		groups = append(groups, g)
	}

	if err = rows.Err(); err != nil {
		log.Printf("ERROR GRUPO: Error durante la iteración de grupos: %v", err) // <--- LOG MEJORADO

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar grupos"})
		return
	}
	log.Printf("DEBUG GRUPO: Se encontraron %d grupos.", len(groups)) // <--- NUEVO LOG
	c.JSON(http.StatusOK, groups)
}

// GetGroupByID obtiene un grupo específico por su ID
func GetGroupByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR GRUPO: ID de grupo inválido para GetGroupByID: %v", err) // <--- LOG MEJORADO
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de grupo inválido"})
		return
	}

	log.Printf("DEBUG GRUPO: Solicitud GetGroupByID para ID: %d", id) // <--- NUEVO LOG

	var group models.Group
	err = config.DB.QueryRow(
		"SELECT id, name, description, topic, created_at, updated_at FROM groups WHERE id = ?",
		id,
	).Scan(&group.ID, &group.Name, &group.Description, &group.Topic, &group.CreatedAt, &group.UpdatedAt)

	if err == sql.ErrNoRows {
		log.Printf("DEBUG GRUPO: Grupo con ID %d no encontrado.", id) // <--- NUEVO LOG
		c.JSON(http.StatusNotFound, gin.H{"error": "Grupo no encontrado"})
		return
	}
	if err != nil {
		log.Printf("ERROR GRUPO: Error al obtener grupo por ID de DB: %v", err) // <--- LOG MEJORADO
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el grupo", "details": err.Error()})
		return
	}
	log.Printf("DEBUG GRUPO: Grupo %d ('%s') encontrado.", group.ID, group.Name) // <--- NUEVO LOG
	c.JSON(http.StatusOK, group)
}