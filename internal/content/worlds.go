package content

type SubSection struct {
	Slug  string
	Title string
}

type WorldSection struct {
	Slug        string
	Title       string
	Path        string
	SubSections []SubSection
}

type World struct {
	Slug        string
	Title       string
	Description string
	Icon        string
	Sections    []WorldSection
}

var WorldOrder = []string{
	"mundo-artistico",
	"mundo-espiritual",
	"mundo-fisico",
	"mundo-mental",
	"mundo-personal",
	"mundo-reflexivo",
}

var Worlds = map[string]World{
	"mundo-artistico": {
		Slug:        "mundo-artistico",
		Title:       "Mundo Artístico",
		Description: "Es el espacio de creatividad y expresión personal. A través del arte, la cultura y las narrativas se exploran emociones, ideas y perspectivas que enriquecen la imaginación y la sensibilidad.",
		Icon:        "/assets/image/avif/planeta1.avif",
		Sections: []WorldSection{
			{
				Slug:  "artes-visuales",
				Title: "Artes Visuales",
				Path:  "/mundos/mundo-artistico/artes-visuales",
				SubSections: []SubSection{
					{Slug: "pintura", Title: "Pintura"},
					{Slug: "fotografia", Title: "Fotografía"},
					{Slug: "narrativa-visual", Title: "Narrativa Visual"},
					{Slug: "arte-digital", Title: "Arte Digital"},
				},
			},
			{
				Slug:  "artes-narrativas",
				Title: "Artes Narrativas y Escénicas",
				Path:  "/mundos/mundo-artistico/artes-narrativas",
				SubSections: []SubSection{
					{Slug: "literatura", Title: "Literatura"},
					{Slug: "teatro", Title: "Teatro"},
					{Slug: "danza", Title: "Danza"},
					{Slug: "musica", Title: "Música"},
				},
			},
			{
				Slug:  "cultura-pensamiento",
				Title: "Cultura y Pensamiento",
				Path:  "/mundos/mundo-artistico/cultura-pensamiento",
				SubSections: []SubSection{
					{Slug: "filosofia", Title: "Filosofía"},
					{Slug: "ciencia", Title: "Ciencia"},
					{Slug: "historia", Title: "Historia"},
					{Slug: "sociedad", Title: "Sociedad"},
				},
			},
		},
	},
	"mundo-espiritual": {
		Slug:        "mundo-espiritual",
		Title:       "Mundo Espiritual",
		Description: "Se relaciona con la conexión interior y el sentido profundo de la vida. Incluye prácticas que cultivan la calma, la atención plena, el equilibrio energético y la reflexión filosófica o espiritual.",
		Icon:        "/assets/image/avif/planeta2.avif",
		Sections: []WorldSection{
			{
				Slug:  "ejercicios-espirituales",
				Title: "Ejercicios Espirituales",
				Path:  "/mundos/mundo-espiritual/ejercicios-espirituales",
				SubSections: []SubSection{
					{Slug: "meditacion", Title: "Meditacion"},
					{Slug: "respiracion-consciente", Title: "Respiracion Consciente"},
					{Slug: "practicas-cuerpo-espiritu", Title: "Practicas Cuerpo-Espiritu"},
					{Slug: "atencion-plena", Title: "Atencion Plena"},
				},
			},
			{
				Slug:  "energia",
				Title: "Renovacion de Energia",
				Path:  "/mundos/mundo-espiritual/energia",
				SubSections: []SubSection{
					{Slug: "energia-vital", Title: "Energia Vital"},
					{Slug: "equilibrio-energetico", Title: "Equilibrio Energetico"},
					{Slug: "conexion-naturaleza", Title: "Conexion con la Naturaleza"},
					{Slug: "atencion-plena", Title: "Atencion Plena"},
				},
			},
			{
				Slug:  "practicas",
				Title: "Practicas Filosoficas y Espirituales",
				Path:  "/mundos/mundo-espiritual/practicas",
				SubSections: []SubSection{
					{Slug: "taoismo", Title: "Taoismo"},
					{Slug: "estoicismo", Title: "Estoicismo"},
					{Slug: "budismo", Title: "Budismo"},
					{Slug: "zen", Title: "Zen"},
				},
			},
		},
	},
	"mundo-fisico": {
		Slug:        "mundo-fisico",
		Title:       "Mundo Fisico",
		Description: "Se centra en el cuidado y la relación con el cuerpo. Incluye la conciencia corporal, la salud, la alimentación, el descanso y la gestión de la energía física para mantener bienestar y vitalidad.",
		Icon:        "/assets/image/avif/planeta3.avif",
		Sections: []WorldSection{
			{
				Slug:  "conciencia-corporal",
				Title: "Conciencia Corporal",
				Path:  "/mundos/mundo-fisico/conciencia-corporal",
				SubSections: []SubSection{
					{Slug: "postura-corporal", Title: "Postura Corporal"},
					{Slug: "movimiento-consciente", Title: "Movimiento Consciente"},
					{Slug: "propiocepcion", Title: "Propiocepcion"},
					{Slug: "conexion-mente-cuerpo", Title: "Conexion Mente-Cuerpo"},
				},
			},
			{
				Slug:  "cuidado-corporal",
				Title: "Cuidado Corporal",
				Path:  "/mundos/mundo-fisico/cuidado-corporal",
				SubSections: []SubSection{
					{Slug: "alimentacion-consciente", Title: "Alimentacion Consciente"},
					{Slug: "nutricion", Title: "Nutricion"},
					{Slug: "autocuidado", Title: "Autocuidado"},
					{Slug: "salud", Title: "Salud"},
				},
			},
			{
				Slug:  "descanso",
				Title: "Descanso y Recuperacion",
				Path:  "/mundos/mundo-fisico/descanso",
				SubSections: []SubSection{
					{Slug: "calidad-sueno", Title: "Calidad del Sueno"},
					{Slug: "ritmos-biologicos", Title: "Ritmos Biologicos"},
					{Slug: "recuperacion-fisica", Title: "Recuperacion Fisica"},
					{Slug: "energia-diaria", Title: "Energia Diaria"},
				},
			},
		},
	},
	"mundo-mental": {
		Slug:        "mundo-mental",
		Title:       "Mundo Mental",
		Description: "Es el ámbito donde se desarrollan las capacidades cognitivas. Involucra la atención, la percepción, el pensamiento crítico y el aprendizaje, permitiendo comprender mejor la realidad y mejorar la forma de pensar.",
		Icon:        "/assets/image/avif/planeta4.avif",
		Sections: []WorldSection{
			{Slug: "aprendizaje", Title: "Aprendizaje", Path: "/mundos/mundo-mental/aprendizaje"},
			{Slug: "atencion", Title: "Atencion", Path: "/mundos/mundo-mental/atencion"},
			{Slug: "percepcion", Title: "Percepcion", Path: "/mundos/mundo-mental/percepcion"},
		},
	},
	"mundo-personal": {
		Slug:        "mundo-personal",
		Title:       "Mundo Personal",
		Description: "Se enfoca en la organización de la vida personal. Incluye el registro de experiencias, la construcción de hábitos, la planificación de proyectos y el seguimiento del crecimiento personal.",
		Icon:        "/assets/image/avif/planeta5.avif",
		Sections: []WorldSection{
			{Slug: "diario", Title: "Diario Personal", Path: "/mundos/mundo-personal/diario"},
			{Slug: "habitos", Title: "Habitos y Rutinas", Path: "/mundos/mundo-personal/habitos"},
			{Slug: "proyectos", Title: "Proyectos Personales", Path: "/mundos/mundo-personal/proyectos"},
		},
	},
	"mundo-reflexivo": {
		Slug:        "mundo-reflexivo",
		Title:       "Mundo Reflexivo",
		Description: "Es el espacio de autoconocimiento y desarrollo interior. Aquí se exploran la identidad personal, las emociones, los valores y el propósito de vida para entender quién se es, cómo se siente uno y qué dirección se quiere tomar.",
		Icon:        "/assets/image/avif/planeta6.avif",
		Sections: []WorldSection{
			{Slug: "autoconcepto", Title: "Autoconcepto", Path: "/mundos/mundo-reflexivo/autoconcepto"},
			{Slug: "gestion-emocional", Title: "Gestion de Emociones", Path: "/mundos/mundo-reflexivo/gestion-emocional"},
			{Slug: "sentido", Title: "Sentido y Direccion", Path: "/mundos/mundo-reflexivo/sentido"},
		},
	},
}

func GetWorldBySlug(slug string) (World, bool) {
	w, ok := Worlds[slug]
	return w, ok
}

func GetSectionBySlug(worldSlug, sectionSlug string) (World, WorldSection, bool) {
	w, ok := Worlds[worldSlug]
	if !ok {
		return World{}, WorldSection{}, false
	}

	for _, section := range w.Sections {
		if section.Slug == sectionSlug {
			return w, section, true
		}
	}

	return World{}, WorldSection{}, false
}

func OrderedWorlds() []World {
	items := make([]World, 0, len(WorldOrder))
	for _, slug := range WorldOrder {
		if w, ok := Worlds[slug]; ok {
			items = append(items, w)
		}
	}
	return items
}

