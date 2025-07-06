# Apoyo-Anonimo-Plataforma-de-Grupos-de-Chat-Seguros-y-en-Tiempo-Real

Descubre "Apoyo Anónimo", una plataforma web robusta y segura diseñada para ofrecer un espacio confidencial donde los usuarios pueden unirse a grupos de chat anónimos para compartir experiencias y recibir apoyo. Desarrollada con un backend en Go, una base de datos MySQL y un frontend en React.

Apoyo Anónimo es una plataforma web innovadora que proporciona un espacio seguro y confidencial para que los usuarios se conecten en grupos de chat anónimos. Diseñada para facilitar el apoyo mutuo en temas sensibles como ansiedad, duelo o adicciones, la aplicación prioriza la privacidad, la seguridad y una experiencia de usuario fluida.

Este proyecto demuestra un enfoque integral en el desarrollo full-stack, aplicando las mejores prácticas en UI/UX, SEO, seguridad, validaciones y optimizaciones de rendimiento, utilizando exclusivamente herramientas de código abierto y gratuitas.

Características Principales

Grupos de Apoyo Anónimos: Los usuarios pueden unirse a grupos basados en temas específicos (ansiedad, duelo, adicciones, etc.) sin revelar su identidad real.
Chat en Tiempo Real: Comunicación instantánea dentro de los grupos gracias a la implementación de WebSockets.
Moderación de Contenido: Un sistema de moderación básico (basado en palabras clave) ayuda a mantener un ambiente seguro y respetuoso.
Gestión de Grupos: Los usuarios pueden ver los grupos existentes y crear nuevos grupos con nombre, descripción y tema.
Experiencia de Usuario Fluida: Interfaz reactiva y dinámica que permite la navegación y la interacción sin recargar la página.
Diseño Responsivo: Adaptabilidad completa a diferentes tamaños de pantalla (escritorio, tablet, móvil) para una experiencia consistente.
Anonimato Garantizado: Generación automática de nombres de usuario anónimos para cada sesión de chat.
Optimización SEO: Meta etiquetas dinámicas, URLs amigables y contenido semántico para una mejor visibilidad en motores de búsqueda.
Seguridad Robusta: Validaciones en frontend y backend, protección contra inyección SQL, configuración de CORS estricta y uso de variables de entorno para secretos.
Accesibilidad (A11y): Implementación de HTML semántico, atributos ARIA, navegación por teclado y contraste de colores adecuado para usuarios con diversas capacidades.
Rendimiento Optimizado: Uso de índices en la base de datos, compresión Gzip en el backend y manejo eficiente de conexiones para respuestas rápidas.

Tecnologías Utilizadas

Este proyecto ha sido construido utilizando un stack de tecnologías moderno y gratuito:

Backend: Go (Golang)
Framework Web: Gin Gonic para construir la API RESTful.
WebSockets: Gorilla WebSocket para la comunicación en tiempo real.
Base de Datos: Conexión a MySQL con database/sql y el driver github.com/go-sql-driver/mysql.
Variables de Entorno: Godotenv para la gestión segura de secretos.
CORS: Gin-contrib/cors para la seguridad de comunicación entre frontend y backend.
Compresión: Gin-contrib/gzip para respuestas HTTP más rápidas.
Base de Datos: MySQL
Almacenamiento persistente de grupos y mensajes.
Uso de índices para optimizar el rendimiento de las consultas.
Frontend: React
Constructor: Vite para un desarrollo rápido y eficiente.
Enrutamiento: React Router DOM para URLs amigables y navegación SPA (Single Page Application).
SEO Dinámico: React Helmet Async para la gestión dinámica de meta etiquetas en el <head>.
Estilos: CSS puro para un control total sobre la UI/UX.

Arquitectura del Proyecto

La aplicación sigue una arquitectura de cliente-servidor:

Frontend (React): Una aplicación de una sola página (SPA) que se ejecuta en el navegador del usuario. Se encarga de la interfaz de usuario, la lógica de presentación y las interacciones del usuario. Se comunica con el backend a través de peticiones HTTP y conexiones WebSocket.
Backend (Go): Un servidor API RESTful que gestiona la lógica de negocio, la interacción con la base de datos y la retransmisión de mensajes en tiempo real a través de WebSockets.
Base de Datos (MySQL): Almacena de forma persistente los datos de los grupos y los mensajes.

apoyo-anonimo/
├── backend/
│   ├── config/             # Configuración de la base de datos
│   ├── handlers/           # Lógica de los controladores de API (grupos, mensajes, WebSockets)
│   ├── models/             # Definición de estructuras de datos (Group, Message)
│   ├── routes/             # Definición de rutas de la API
│   ├── main.go             # Punto de entrada del servidor Go
│   ├── .env.example        # Ejemplo de variables de entorno (NO SUBIR .env real)
│   └── go.mod, go.sum      # Módulos Go
└── frontend/
    ├── public/             # Archivos estáticos (index.html, favicon.ico)
    ├── src/
    │   ├── assets/         # Activos como imágenes, iconos
    │   ├── components/     # Componentes React (GroupList, ChatRoom)
    │   ├── services/       # Servicios para la comunicación con la API (api.js)
    │   ├── App.jsx         # Componente principal de React
    │   ├── main.jsx        # Punto de entrada de la aplicación React
    │   └── index.css       # Estilos globales
    └── package.json, vite.config.js # Configuración de dependencias y build de Vite

Configuración y Ejecución Local
Sigue estos pasos para levantar la aplicación en tu entorno de desarrollo.

Prerequisitos
Go (Golang) (versión 1.18 o superior)

Node.js y npm

1. Configuración de la Base de Datos (MySQL)
Asegúrate de que tu servidor MySQL esté corriendo.

-- Conéctate a tu MySQL CLI o usa un cliente como MySQL Workbench
-- Crea la base de datos si no existe
CREATE DATABASE IF NOT EXISTS apoyo_anonimo_db;

-- Usa la base de datos
USE apoyo_anonimo_db;

-- Tabla para los grupos
CREATE TABLE IF NOT EXISTS groups (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    topic VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabla para los mensajes
CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT NOT NULL,
    username VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_moderated BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- Índices para búsquedas rápidas (¡optimización!)
CREATE INDEX idx_groups_topic ON groups (topic);
CREATE INDEX idx_groups_name ON groups (name);
CREATE INDEX idx_messages_group_id ON messages (group_id);
CREATE INDEX idx_messages_timestamp ON messages (timestamp);

-- (Opcional) Inserta datos de ejemplo en la tabla 'groups' para empezar
INSERT INTO groups (name, description, topic) VALUES
('Superando la Ansiedad Diaria', 'Un espacio seguro para compartir experiencias y estrategias para manejar la ansiedad en el día a día.', 'Ansiedad'),
('Círculo de Resiliencia ante el Duelo', 'Apoyo mutuo para quienes están procesando la pérdida de un ser querido. Juntos encontramos fuerza.', 'Duelo'),
('Camino a la Recuperación (Adicciones)', 'Grupo de apoyo confidencial para personas en proceso de recuperación de adicciones. Un paso a la vez.', 'Adicciones'),
('Manejo del Estrés Laboral', 'Compartimos técnicas y experiencias para lidiar con el estrés en el entorno de trabajo.', 'Estrés'),
('Apoyo para Padres (Adolescentes)', 'Un lugar para padres de adolescentes que buscan orientación y apoyo en esta etapa.', 'Familia'),
('Viviendo con Depresión', 'Espacio para compartir y entender la depresión, buscando herramientas para una vida plena.', 'Depresión'),
('Nuevos Comienzos (Separación/Divorcio)', 'Ayuda y orientación para quienes atraviesan una separación o divorcio.', 'Duelo'),
('Ansiedad Social y Timidez', 'Superando juntos la ansiedad en situaciones sociales y la timidez.', 'Ansiedad'),
('Gestión de la Ira y Frustración', 'Aprende a identificar y manejar la ira de forma constructiva.', 'Estrés'),
('Apoyo a Cuidadores de Adultos Mayores', 'Un espacio para cuidadores que necesitan apoyo emocional y compartir experiencias.', 'Familia');

-- (Opcional) Inserta datos de ejemplo en la tabla 'messages'
-- Asegúrate de que los group_id coincidan con los IDs generados por MySQL (normalmente 1, 2, 3...)
INSERT INTO messages (group_id, username, content) VALUES
(1, 'AlmaValiente', 'Hola a todos, es mi primera vez aquí. Me siento muy ansiosa hoy.'),
(1, 'LuzInterior', 'Bienvenida, AlmaValiente. Es un paso valiente estar aquí. Recuerda respirar profundo.'),
(2, 'EsperanzaRenovada', 'Es difícil, pero cada día es un poco menos pesado. Fuerza a todos.'),
(3, 'FenixRising', 'Un día más sobrio. La lucha es real, pero la recompensa es mayor.'),
(4, 'TrabajadorZen', 'Hola, ¿alguien tiene consejos para el estrés de las reuniones?'),
(5, 'PadrePreocupado', 'Mis hijos adolescentes no me hablan. ¿Es normal?');

2. Configuración y Ejecución del Backend (Go)
Navega a la carpeta backend:

cd backend

Instala las dependencias de Go:

go mod tidy

Crea un archivo .env en la carpeta backend/ con tus credenciales de MySQL y la URL del frontend. 

apoyo-anonimo/backend/.env

MYSQL_DSN="root:tu_contraseña_mysql@tcp(127.0.0.1:3306)/apoyo_anonimo_db?charset=utf8mb4&parseTime=True&loc=Local"
PORT=8000
FRONTEND_URL="http://localhost:5173"

Asegúrate de reemplazar tu_contraseña_mysql con tu contraseña real. Si no tienes contraseña para root, usa root:@.

Configura las variables de entorno en tu terminal (antes de ejecutar go run main.go).

En Linux/macOS (Bash/Zsh) o Git Bash (Windows):

export MYSQL_DSN="root:tu_contraseña_mysql@tcp(127.0.0.1:3306)/apoyo_anonimo_db?charset=utf8mb4&parseTime=True&loc=Local"
export PORT=8000
export FRONTEND_URL="http://localhost:5173"

En Windows CMD:

set MYSQL_DSN="root:tu_contraseña_mysql@tcp(127.0.0.1:3306)/apoyo_anonimo_db?charset=utf8mb4&parseTime=True&loc=Local"
set PORT=8000
set FRONTEND_URL="http://localhost:5173"

En Windows PowerShell:

$env:MYSQL_DSN="root:tu_contraseña_mysql@tcp(127.0.0.1:3306)/apoyo_anonimo_db?charset=utf8mb4&parseTime=True&loc=Local"
$env:PORT=8000
$env:FRONTEND_URL="http://localhost:5173"

Ejecuta el servidor Go:

go run main.go

Verás un mensaje indicando que el servidor está escuchando en el puerto 8000. Deja esta terminal abierta.

3. Configuración y Ejecución del Frontend (React)
Abre una nueva ventana de terminal y navega a la carpeta frontend:

cd frontend

Asegúrate de que tus dependencias de React estén en la versión 18.2.0 en package.json para evitar conflictos con react-helmet-async.

"dependencies": {
  "react": "^18.2.0",     
  "react-dom": "^18.2.0",
  "react-router-dom": "^7.6.3",
  "react-helmet-async": "^1.3.0"
},

Realiza una limpieza profunda y reinstala las dependencias (¡importante!):

rm -rf node_modules
rm package-lock.json
rm -rf .vite
npm cache clean --force
npm install

Inicia el servidor de desarrollo de Vue:

npm run dev

El servidor frontend se iniciará, generalmente en http://localhost:5173/.

¡Listo!
Abre tu navegador y visita http://localhost:5173/ para ver la aplicación "Apoyo Anónimo" en acción.

Mejoras Implementadas y Buenas Prácticas

Este proyecto ha sido desarrollado prestando especial atención a las siguientes áreas:

UI/UX (Interfaz de Usuario / Experiencia de Usuario)

Diseño Responsivo: La interfaz se adapta fluidamente a cualquier tamaño de pantalla (móvil, tablet, escritorio) utilizando CSS Grid y Flexbox.

Feedback Visual: Indicadores de carga ("Cargando grupos..."), mensajes de éxito (alertas simples por ahora) y deshabilitación de botones durante el envío de mensajes (isSending) para evitar clics duplicados.

Animaciones Sutiles: Transiciones CSS para una experiencia más fluida y atractiva.

Consistencia Visual: Estilos CSS unificados para tipografía, colores y espaciado, garantizando una apariencia profesional.

Scroll Automático en Chat: La ventana de chat se desplaza automáticamente al final para mostrar los mensajes más recientes.

SEO (Optimización para Motores de Búsqueda)

Meta Etiquetas Dinámicas: Integración de react-helmet-async para que el título y la descripción de la página (y Open Graph/Twitter Cards) se actualicen dinámicamente según el grupo de chat visitado, mejorando la indexación y la compartibilidad en redes sociales.

URLs Amigables: Uso de react-router-dom para crear rutas legibles y significativas (ej. /groups/1 en lugar de /?groupId=1).

Contenido Semántico: Empleo de etiquetas HTML5 semánticas (<header>, <main>, <section>, <h1> a <h6>) para estructurar el contenido de manera comprensible para los motores de búsqueda.

Compresión Gzip: El backend Go comprime las respuestas HTTP, reduciendo el tiempo de carga de la página, un factor clave para el SEO.

Validaciones

Validación de Cliente (Frontend): Comprobaciones iniciales en el navegador (ej. campos requeridos, longitud máxima de mensajes) para proporcionar feedback instantáneo al usuario y reducir solicitudes inválidas al servidor.

Validación de Servidor (Backend): Validaciones robustas en Go (ej. nombre de grupo no vacío, longitud de mensaje, tema válido) para asegurar la integridad y seguridad de los datos antes de que se almacenen en la base de datos.

Seguridad

Variables de Entorno: Uso de archivos .env para gestionar credenciales sensibles (ej. MYSQL_DSN).

CORS Estricto: Configuración explícita de Cross-Origin Resource Sharing en el backend para permitir solo solicitudes desde el dominio del frontend, previniendo ataques de sitios cruzados.

Prevención de Inyección SQL: Las consultas a la base de datos en Go utilizan parámetros (?), lo que automáticamente previene la inyección SQL.

Moderación de Contenido: Implementación de una lista de palabras prohibidas para filtrar mensajes y mantener un ambiente seguro en los grupos.

Anonimato: Diseño centrado en no almacenar información personal identificable del usuario.

Accesibilidad (A11y)

HTML Semántico: Uso de elementos HTML con su propósito correcto para mejorar la comprensión por parte de tecnologías de asistencia.

Atributos ARIA: Inclusión de aria-label en elementos interactivos para proporcionar descripciones significativas a los lectores de pantalla.

Navegación por Teclado: Asegura que todos los elementos interactivos sean accesibles y operables usando solo el teclado.

Contraste de Colores: Elección de una paleta de colores que garantiza un contraste suficiente entre texto y fondo para una legibilidad óptima.

Mensajes Claros: Textos de error y éxito concisos y fáciles de entender para todos los usuarios.

Optimizaciones de Backend y Base de Datos

Índices de Base de Datos: Creación de índices en las tablas groups y messages para acelerar las operaciones de búsqueda y filtrado.

Manejo de Concurrencia (Go): Uso de goroutines para tareas no bloqueantes (como BroadcastMessage) y sync.Mutex para proteger el acceso a recursos compartidos (mapas de conexiones WebSocket), garantizando la estabilidad bajo carga.

Pool de Conexiones a DB: Go gestiona automáticamente un pool de conexiones a la base de datos, mejorando la eficiencia y el rendimiento al reutilizar conexiones existentes.

WebSockets: Implementación de WebSockets para el chat en tiempo real, reduciendo la latencia y la carga del servidor en comparación con el polling HTTP tradicional.

Contribución
¡Las contribuciones son bienvenidas! Siéntete libre de abrir un issue o enviar un pull request.
