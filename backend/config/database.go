package config

import (
	"database/sql" // Para hablar con la base de datos
	"fmt"          // Para imprimir mensajes
	"log"          // Para errores importantes
	"os"           // Para leer los secretos del archivo .env

	_ "github.com/go-sql-driver/mysql" // El "traductor" para MySQL
)

var DB *sql.DB // Aquí guardamos la conexión a la base de datos para usarla en otros lugares

// ConnectDB es la función que conecta nuestros cerebros con el cuaderno mágico
func ConnectDB() {
	dsn := os.Getenv("MYSQL_DSN") // Leemos el secreto de la conexión
	if dsn == "" {
		log.Fatal("Error: La variable de entorno MYSQL_DSN no está configurada. Revisa tu archivo .env")
	}

	var err error
	// Abrimos la conexión a la base de datos
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a la base de datos: %v", err)
	}

	// Hacemos una pequeña prueba para ver si la conexión funciona
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v. Revisa tus credenciales y si MySQL está encendido.", err)
	}

	fmt.Println("¡Conectado exitosamente a nuestro cuaderno mágico (base de datos MySQL)!")
}