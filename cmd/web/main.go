// Package main es el punto de entrada de la aplicación.
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	appdb "el-mundo-interior/internal/db"
	appweb "el-mundo-interior/internal/web"
)

func main() {
	// Carga variables del fichero .env si existe (ignora el error si no hay fichero).
	_ = godotenv.Load()

	dbPath := getenv("DB_PATH", "./datos.db")
	port := getenv("PORT", "8080")
	secret := getenv("SESSION_SECRET", "")

	if secret == "" {
		log.Println("AVISO: SESSION_SECRET no definido en .env; usando sesiones sin firma")
	}

	db, err := appdb.NewDB(dbPath)
	if err != nil {
		log.Fatalf("error abriendo base de datos: %v", err)
	}
	defer db.Close()

	addr := ":" + port
	server := appweb.NewServer(addr, db)

	log.Printf("servidor escuchando en http://localhost%s", server.Addr())

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
