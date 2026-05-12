package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

var Questions = []string{
	"¿Qué emoción quieres transmitir con esta expresión?",
	"¿Qué momento de la realidad estás eligiendo capturar y por qué merece ser observado?",
	"¿Qué historia puede entenderse sin necesidad de palabras?",
	"¿Qué herramienta o técnica te permitiría expresar mejor la idea que tienes ahora?",
	"¿Qué parte de ti mismo quieres explorar hoy?",
	"¿Qué te detiene y qué te impulsa en este momento?",
	"¿Qué símbolo o imagen representa mejor lo que sientes ahora?",
	"¿Qué cambiarías de la situación que tienes en mente?",
	"¿Qué necesitas soltar para avanzar?",
	"¿Qué patrón repites que ya no te sirve?",
	"¿Cómo describirías tu estado interno con un solo color?",
	"¿Qué conversación pendiente llevas contigo?",
	"¿Qué versión de ti mismo quieres honrar hoy?",
	"¿Qué te costaría más admitir en voz alta?",
	"¿Cuándo fue la última vez que te sentiste completamente presente?",
	"¿Qué límite necesitas poner o quitar?",
	"¿Qué historia te cuentas sobre ti que ya no es cierta?",
	"¿Qué haría la versión más valiente de ti en este momento?",
	"¿Qué aspecto de tu vida necesita más atención y ternura?",
	"¿De qué manera el miedo está dando forma a tus decisiones ahora mismo?",
}

// ApiGetQuestion maneja GET /api/questions y devuelve una pregunta aleatoria.
func ApiGetQuestion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := Questions[rand.Intn(len(Questions))]
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"question": q})
	}
}
