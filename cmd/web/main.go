package main

import (
	"log"

	"el-mundo-interior/internal/web"
)

func main() {
	srv := web.NewServer(":8080")
	log.Println("server listening on http://localhost:8080")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
