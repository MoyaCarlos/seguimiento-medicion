package main

type config struct {
	Puerto   string
	RutaBase string
	DirWeb   string

	DemoHabilitado bool
}

func cargarConfig(getenv func(string) string) config {
	return config{
		Puerto:   valorOPorDefecto(getenv("PORT"), "8080"),
		RutaBase: valorOPorDefecto(getenv("DB_PATH"), "seguimiento.db"),
		DirWeb:   valorOPorDefecto(getenv("WEB_DIR"), "web/dist"),

		DemoHabilitado: getenv("DEMO_HABILITADO") == "true",
	}
}

func valorOPorDefecto(valor, porDefecto string) string {
	if valor == "" {
		return porDefecto
	}
	return valor
}
