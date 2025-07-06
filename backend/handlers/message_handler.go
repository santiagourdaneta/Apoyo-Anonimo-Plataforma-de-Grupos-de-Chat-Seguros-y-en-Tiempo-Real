package handlers

import (
	"encoding/json" // Para convertir mensajes a JSON
	"log"
	"net/http"
	"strconv"
	"strings" // Para la moderación de texto
	"sync"    // Para el mutex
	"time"

	"apoyo-anonimo/backend/config"
	"apoyo-anonimo/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// --- VARIABLES GLOBALES PARA WEBSOCKETS ---
// Estas deben ir aquí, después de los imports y antes de cualquier función.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// ¡IMPORTANTE! En producción, esto debe ser más estricto.
		// Solo permitir conexiones desde tu frontend.
		return true // Por ahora, permitimos todo para desarrollo.
	},
}

// WebSocketConnections guarda las conexiones activas de WebSocket por GroupID
// Es un mapa de ID de grupo a una lista de conexiones WebSocket.
var WebSocketConnections = make(map[uint][]*websocket.Conn)

// wsConnectionsMutex protege el mapa WebSocketConnections de accesos concurrentes.
// Es crucial para evitar "race conditions" cuando múltiples goroutines acceden al mapa.
var wsConnectionsMutex = &sync.Mutex{}

// Lista de palabras prohibidas para la moderación básica
var forbiddenWords = []string{"malo", "feo", "grosero", "tonto", "odio", "matar", "daño"}

// --- FUNCIONES ---

// checkModeration revisa si un mensaje contiene palabras prohibidas
func checkModeration(content string) bool {
	lowerContent := strings.ToLower(content)
	for _, word := range forbiddenWords {
		if strings.Contains(lowerContent, word) {
			return true // Contiene una palabra prohibida
		}
	}
	return false // No contiene palabras prohibidas
}

// SendMessage envía un nuevo mensaje a un grupo
func SendMessage(c *gin.Context) {
	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR: Failed to bind JSON for message: %v", err) // <--- AÑADE ESTA LÍNEA
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de mensaje inválidos", "details": err.Error()})
		return
	}

	// Validaciones básicas
	if req.GroupID == 0 || req.Content == "" || len(req.Content) > 500 { // Límite de 500 caracteres
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mensaje inválido: ID de grupo, contenido o longitud incorrectos."})
		return
	}
	if req.Username == "" {
		req.Username = "Anónimo" // Por defecto si no se proporciona
	}

	isModerated := checkModeration(req.Content) // Revisamos si tiene palabras malas

	result, err := config.DB.Exec(
		"INSERT INTO messages (group_id, username, content, is_moderated) VALUES (?, ?, ?, ?)",
		req.GroupID, req.Username, req.Content, isModerated,
	)
	if err != nil {
		log.Printf("Error al insertar mensaje en DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el mensaje", "details": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error al obtener ID del nuevo mensaje: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el mensaje"})
		return
	}

	newMessage := models.Message{
		ID: uint(id),
		GroupID: req.GroupID,
		Username: req.Username,
		Content: req.Content,
		Timestamp: time.Now(), // La base de datos lo pondrá, pero lo simulamos aquí
		IsModerated: isModerated,
	}

	if isModerated {
		log.Printf("Mensaje moderado detectado en grupo %d de %s: %s", req.GroupID, req.Username, req.Content)
	}

	c.JSON(http.StatusCreated, newMessage)

	// Enviar el mensaje a todos los que están escuchando en el WebSocket (¡la parte del chat en tiempo real!)
	go BroadcastMessage(newMessage.GroupID, newMessage) // Usamos 'go' para que no bloquee la respuesta HTTP
}

// GetMessagesByGroup obtiene todos los mensajes de un grupo
func GetMessagesByGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		log.Printf("ERROR MENSAJE: ID de grupo inválido para GetMessagesByGroup: %v", err) // <--- LOG MEJORADO
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de grupo inválido"})
		return
	}

	log.Printf("DEBUG MENSAJE: Solicitud GetMessagesByGroup para ID: %d", groupID) // <--- NUEVO LOG

	rows, err := config.DB.Query(
		"SELECT id, group_id, username, content, timestamp, is_moderated FROM messages WHERE group_id = ? ORDER BY timestamp ASC",
		groupID,
	)
	if err != nil {
		log.Printf("ERROR MENSAJE: Error al obtener mensajes de DB para grupo %d: %v", groupID, err) // <--- LOG MEJORADO
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener mensajes", "details": err.Error()})
		return
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.GroupID, &m.Username, &m.Content, &m.Timestamp, &m.IsModerated); err != nil {
			log.Printf("ERROR MENSAJE: Error al escanear mensaje para grupo %d: %v", groupID, err) // <--- LOG MEJORADO
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar mensajes", "details": err.Error()})
			return
		}
		messages = append(messages, m)
	}

	if err = rows.Err(); err != nil {
		log.Printf("ERROR MENSAJE: Error durante la iteración de mensajes para grupo %d: %v", groupID, err) // <--- LOG MEJORADO
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar mensajes"})
		return
	}
	log.Printf("DEBUG MENSAJE: Se encontraron %d mensajes para el grupo %d.", len(messages), groupID) // <--- NUEVO LOG
	c.JSON(http.StatusOK, messages)
}

// HandleWebSocket maneja las conexiones WebSocket para el chat en tiempo real
func HandleWebSocket(c *gin.Context) {
	log.Printf("DEBUG WS: Intentando actualizar a WebSocket para URL: %s", c.Request.URL.String())
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ERROR WS: Error al actualizar a WebSocket para URL %s: %v", c.Request.URL.String(), err)
		return
	}
	defer conn.Close()

	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		log.Printf("ERROR WS: ID de grupo inválido en WebSocket: %v", err)
		return
	}
	groupUintID := uint(groupID)

	// Añadir la conexión al mapa de conexiones del grupo, protegido por mutex
	wsConnectionsMutex.Lock()
	WebSocketConnections[groupUintID] = append(WebSocketConnections[groupUintID], conn)
	wsConnectionsMutex.Unlock()
	log.Printf("DEBUG WS: Nuevo WebSocket conectado al grupo %d. Total conexiones: %d", groupUintID, len(WebSocketConnections[groupUintID]))

	// Remover la conexión cuando se cierre, protegido por mutex
	defer func() {
		wsConnectionsMutex.Lock()
		for i, existingConn := range WebSocketConnections[groupUintID] {
			if existingConn == conn {
				WebSocketConnections[groupUintID] = append(WebSocketConnections[groupUintID][:i], WebSocketConnections[groupUintID][i+1:]...)
				break
			}
		}
		wsConnectionsMutex.Unlock()
		log.Printf("DEBUG WS: WebSocket desconectado del grupo %d. Conexiones restantes: %d", groupUintID, len(WebSocketConnections[groupUintID]))
	}()

	// Bucle para leer mensajes del cliente (mantener la conexión abierta)
	for {
		// ReadMessage bloqueará hasta que se reciba un mensaje o la conexión se cierre.
		// Esto mantiene la goroutine viva y la conexión abierta.
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("DEBUG WS: WebSocket cerrado normalmente para grupo %d. Tipo de cierre: %v", groupUintID, err)
			} else {
				log.Printf("ERROR WS: Error de lectura WebSocket para grupo %d: %v", groupUintID, err)
			}
			break
		}
		// Si el cliente envía mensajes, puedes procesarlos aquí.
		// Para este chat, los mensajes se envían vía POST /messages.
	}
	log.Printf("DEBUG WS: Bucle de lectura de WebSocket terminado para grupo %d.", groupUintID)
}

// BroadcastMessage envía un mensaje a todas las conexiones WebSocket de un grupo
func BroadcastMessage(groupID uint, message models.Message) {
	wsConnectionsMutex.Lock() // Bloquear el mutex antes de acceder al mapa
	connections := WebSocketConnections[groupID]
	wsConnectionsMutex.Unlock() // Desbloquear el mutex después de obtener las conexiones

	if connections == nil {
		return // No hay conexiones para este grupo
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar mensaje para broadcast: %v", err)
		return
	}

	for _, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Error al enviar mensaje a WebSocket del grupo %d: %v", groupID, err)
			// Podrías añadir lógica para remover conexiones rotas aquí
		}
	}
}