// Package main es el punto de entrada de la aplicación.
package main

import (
	"log"

	appdb "el-mundo-interior/internal/db"
	appweb "el-mundo-interior/internal/web"
)

func main() {
	// Abrimos (o creamos) el fichero de base de datos.
	// Si no existe, migrate() lo crea con las tablas vacías.
	db, err := appdb.NewDB("./datos.db")
	if err != nil {
		log.Fatalf("error abriendo base de datos: %v", err)
	}
	defer db.Close()

	// Creamos el servidor pasándole la conexión a la BD.
	server := appweb.NewServer(":8080", db)

	log.Printf("servidor escuchando en http://localhost%s", server.Addr())

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
