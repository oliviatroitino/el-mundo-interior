// Package main es el punto de entrada de la aplicación.
// Aquí se inicia el servidor HTTP y se mantiene escuchando en el puerto especificado.
package main

import (
	"log"

	// Importamos el paquete web interno que contiene la lógica del servidor
	appweb "el-mundo-interior/internal/web"
)

func main() {
	// Creamos una nueva instancia del servidor en el puerto 8080
	server := appweb.NewServer(":8080")

	// Informamos al usuario dónde se está ejecutando el servidor
	log.Printf("servidor escuchando en http://localhost%s", server.Addr())

	// Iniciamos el servidor. Si hay un error, el programa termina con log.Fatal
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}