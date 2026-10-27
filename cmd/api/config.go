package main

type config struct {
	Puerto   string
	RutaBase string
	DirWeb   string
}

func cargarConfig(getenv func(string) string) config {
	return config{
		Puerto:   valorOPorDefecto(getenv("PORT"), "8080"),
		RutaBase: valorOPorDefecto(getenv("DB_PATH"), "seguimiento.db"),
		DirWeb:   valorOPorDefecto(getenv("WEB_DIR"), "web/dist"),
	}
}

func valorOPorDefecto(valor, porDefecto string) string {
	if valor == "" {
		return porDefecto
	}
	return valor
}
