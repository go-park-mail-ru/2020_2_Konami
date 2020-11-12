package cors_init

import "github.com/rs/cors"

func InitCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:80"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "X-Content-Type-Options", "Csrf-Token"},
		Debug:            false,
	})
}
