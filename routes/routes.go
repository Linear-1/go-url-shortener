package routes

import ("net/http"
		"url-shortener/handlers"


)

func Register(mux *http.ServeMux , urlHandler *handlers.URLHandler) {
	mux.HandleFunc("GET /health", handlers.HealthCheck)

	mux.HandleFunc("POST /api/v1/urls", urlHandler.CreateURL)

	mux.HandleFunc("GET /{shortCode}" , urlHandler.RedirectURL)

}

