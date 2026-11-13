package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"url-shortener/models"
	"url-shortener/services"
	"url-shortener/validation"
)

type URLHandler struct{
	
	service *services.URLService
	baseURL string

}

func NewURLHandler(service *services.URLService, baseURL string) *URLHandler  {
	return &URLHandler{
		service: service,
		baseURL: baseURL,
	}
}


func HealthCheck(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintln(w, "URL shortener is running!")
}

func (h *URLHandler) CreateURL(w http.ResponseWriter , r *http.Request){
	var req models.CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w , "Invalid JSON request body" , http.StatusBadRequest)
		return 
	}

	err = validation.ValidationCreateURL(&req)
	if err != nil{
		http.Error(w ,err.Error() , http.StatusBadRequest)
		return
	}

	createdURL, err := h.service.CreateShortURL(r.Context(), req.URL)

	if err != nil{
		log.Printf("CreateShortURL failed: %v", err)
		http.Error(w, "Failed to create short URL" , http.StatusInternalServerError)
		return
	}

	shortURL := strings.TrimRight(h.baseURL, "/") + "/" + createdURL.ShortCode

	response := models.CreateURLResponse{
		ShortURL: shortURL,
		OriginalURL: createdURL.OriginalURL,
		ExpiresAt: createdURL.ExpiresAt,
	}
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(response); err != nil{
		http.Error(w, "Failed to encode response" , http.StatusInternalServerError)
		return 
	}
	


	w.Header().Set("Content-Type" , "application/json")
	w.WriteHeader(http.StatusCreated)

	if _ , err := w.Write(buf.Bytes()); err != nil{
		log.Printf("failed to write response: %v", err)
		return
	}

}

func (h *URLHandler) RedirectURL(w http.ResponseWriter , r *http.Request)  {
	
	shortCode := r.PathValue("shortCode")

	originalURL , err := h.service.RedirectURL(r.Context(),shortCode)
	
	if errors.Is(err , services.ErrURLNotFound){
		http.Error(w , err.Error(), http.StatusNotFound)
		return
	}

	if errors.Is(err , services.ErrURLExpired){
		http.Error(w, err.Error(), http.StatusGone)
		return
	}

	if err != nil {
		http.Error(w ,"Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w ,r, originalURL, http.StatusFound)

}

























