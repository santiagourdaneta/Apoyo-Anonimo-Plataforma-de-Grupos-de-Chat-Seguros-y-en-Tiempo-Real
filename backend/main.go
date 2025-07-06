package main

import (
	"log"
	"os" // Para leer secretos de tu computadora

	"apoyo-anonimo/backend/config"  // Nuestra configuración del cuaderno mágico
	"apoyo-anonimo/backend/routes"  // Nuestros caminos para hablar
	"github.com/gin-gonic/gin"      // Nuestro ayudante para hablar por internet
	"github.com/gin-contrib/cors"   // Para que la cara y los cerebros se entiendan
	"github.com/joho/godotenv"      // Para leer secretos de un archivo .env
	"github.com/gin-contrib/gzip" // Importa esto
)

func main() {
	// 1. Cargamos nuestros secretos (como la contraseña del cuaderno mágico)
	// Esto busca un archivo llamado .env en la misma carpeta.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env. Asegúrate de tener uno con tus secretos.")
	}

	// 2. Conectamos con nuestro cuaderno mágico (la base de datos)
	config.ConnectDB()

	// 3. Creamos nuestro "hablador" principal (el router de Gin)
	router := gin.Default() // gin.Default() es como un router que ya sabe hacer cosas básicas.
	router.Use(gzip.Gzip(gzip.DefaultCompression)) // <--- Añade esta línea DESPUÉS de gin.Default()

	// 4. Decimos que la cara (frontend) y los cerebros (backend) pueden hablar
	// Esto es importante para que tu página web pueda pedir cosas a tu cerebro.
	corsConfig := cors.DefaultConfig()
	// ¡IMPORTANTE! Aquí le decimos desde dónde puede hablar nuestra "cara".
	// Para probar, usamos localhost:5173 (que es donde React suele vivir).
	// Si tu React vive en otro lugar, cámbialo aquí.
	frontendURL := os.Getenv("FRONTEND_URL")


	    // --- AÑADE ESTA LÍNEA PARA DEPURAR ---
		log.Printf("DEBUG: FRONTEND_URL leída del entorno: '%s'", frontendURL)
	    // --- FIN LÍNEA DE DEPURACIÓN ---

	if frontendURL == "" {
		frontendURL = "http://localhost:5173" // Valor por defecto para desarrollo

		   // --- AÑADE ESTA OTRA LÍNEA DE DEPURACIÓN ---
        log.Printf("DEBUG: FRONTEND_URL usando valor por defecto: '%s'", frontendURL)
        // --- FIN LÍNEA DE DEPURACIÓN ---
	}
	corsConfig.AllowOrigins = []string{frontendURL}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"} // Authorization para cuando tengamos usuarios de verdad
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 300 // Cuánto tiempo se guarda la información de CORS
	router.Use(cors.New(corsConfig))

	// 5. Configuramos los caminos para hablar (nuestras rutas de API)
	routes.SetupAPIRoutes(router)

	// 6. Encendemos los cerebros para que empiecen a escuchar
	port := os.Getenv("PORT") // Leemos el puerto de los secretos
	if port == "" {
		port = "8000" // Si no hay secreto, usamos el puerto 8000
	}
	log.Printf("Servidor Go escuchando en el puerto :%s", port)
	log.Fatal(router.Run(":" + port)) // ¡A correr! Si hay un error, se detiene.
}