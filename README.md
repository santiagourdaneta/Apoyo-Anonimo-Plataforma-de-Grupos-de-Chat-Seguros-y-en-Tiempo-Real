# Apoyo Anónimo – Plataforma de Grupos de Chat Seguros y en Tiempo Real

**Apoyo Anónimo** es una plataforma web robusta y segura diseñada para ofrecer un espacio confidencial donde los usuarios pueden unirse a **grupos de chat anónimos** para compartir experiencias y recibir apoyo emocional. Está pensada para tratar temas sensibles como **ansiedad, duelo, adicciones, depresión** y más, priorizando siempre el **anonimato**, la **seguridad** y la **usabilidad**.

---

## ✨ Características Principales

- **Grupos de Apoyo Anónimos** por temática
- **Chat en Tiempo Real** usando WebSockets
- **Moderación de Contenido** basada en palabras clave
- **Creación y Gestión de Grupos**
- **Experiencia de Usuario Fluida** (SPA)
- **Diseño Responsivo** para móvil y escritorio
- **Anonimato Garantizado** con nombres aleatorios
- **SEO Optimizado** con React Helmet
- **Seguridad Robusta** (validaciones, CORS, .env, SQL seguro)
- **Accesibilidad A11y** (navegación por teclado, contraste, ARIA)
- **Rendimiento Óptimo** con índices DB, gzip, WebSocket y manejo concurrente

---

## 🧱 Tecnologías Utilizadas

### Backend – Go (Golang)
- Framework: `gin-gonic`
- WebSockets: `gorilla/websocket`
- DB: MySQL (`database/sql` + `go-sql-driver/mysql`)
- Entorno: `godotenv`
- CORS: `gin-contrib/cors`
- Compresión: `gin-contrib/gzip`

### Frontend – React (SPA)
- Constructor: `Vite`
- Enrutamiento: `react-router-dom`
- SEO Dinámico: `react-helmet-async`
- Estilos: CSS puro

---

## 🧠 Arquitectura del Proyecto

apoyo-anonimo/
├── backend/
│ ├── config/
│ ├── handlers/
│ ├── models/
│ ├── routes/
│ ├── main.go
│ └── .env.example
└── frontend/
├── public/
├── src/
│ ├── assets/
│ ├── components/
│ ├── services/
│ ├── App.jsx
│ ├── main.jsx
│ └── index.css
└── vite.config.js


---

## 🚀 Instalación y Ejecución Local

### 🔧 Requisitos
- [Go](https://go.dev/) 1.18+
- [Node.js](https://nodejs.org/) y npm
- [MySQL](https://www.mysql.com/)

---

### 1. Configurar Base de Datos MySQL

```sql
CREATE DATABASE IF NOT EXISTS apoyo_anonimo_db;

USE apoyo_anonimo_db;

CREATE TABLE groups (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    topic VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT NOT NULL,
    username VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_moderated BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE INDEX idx_groups_topic ON groups (topic);
CREATE INDEX idx_groups_name ON groups (name);
CREATE INDEX idx_messages_group_id ON messages (group_id);
CREATE INDEX idx_messages_timestamp ON messages (timestamp);

Inserta datos de ejemplo en las tablas groups y messages.

2. Ejecutar Backend

cd backend
go mod tidy
cp .env.example .env
# edita el archivo .env con tus credenciales de MySQL

go run main.go

3. Ejecutar Frontend

cd frontend
npm install
npm run dev

Accede en tu navegador a: http://localhost:5173

🛡️ Seguridad y Buenas Prácticas

🔒 SQL seguro con prepared statements (?)
🔐 Variables sensibles gestionadas con .env
🧼 Validaciones robustas (frontend y backend)
⚙️ CORS restrictivo por origen
💬 Moderación básica por palabras clave

♿ Accesibilidad (A11y)

HTML5 semántico
Atributos aria-* relevantes
Navegación por teclado asegurada
Contraste de colores apto para todas las personas

📈 Optimización y Rendimiento

Gzip backend (gin-contrib/gzip)
Índices en DB para búsquedas rápidas
WebSocket en lugar de polling
Concurrencia eficiente en Go con goroutines y sync.Mutex

🤝 Contribuciones

¡Las contribuciones son bienvenidas!
Abre un Issue con sugerencias o errores
Envía un Pull Request siguiendo la guía de estilo
Propón mejoras de UI/UX, seguridad, nuevas funciones o refactorizaciones

❤️ Agradecimientos
A todas las personas que buscan apoyo o lo ofrecen, desde el anonimato y el respeto mutuo. Este proyecto es para ustedes.


