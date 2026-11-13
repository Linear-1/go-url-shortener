package validation

import (
	"url-shortener/models"
	"net/url"
	"errors"
	"unicode"
    "strings"

)
	
func ValidationCreateURL(req *models.CreateURLRequest) error {
	req.URL = strings.TrimSpace(req.URL)
	
	if req.URL == "" {
		return errors.New("URL cannot be empty")
	}
	if len(req.URL) > 2048 {
		return errors.New("URL exceeds maximum length of 2048 characters")
	}

	for _,char := range req.URL {
		if unicode.IsSpace(char) {
			return errors.New("URL cannot contain whitespace")
		}
	}

	parsedURL , err := url.Parse(req.URL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must use HTTP or HTTPS")
	}

	if parsedURL.Hostname() == "" {
		return  errors.New("URL must contain a hostname")
	}

	if parsedURL.User != nil {
		return errors.New("URL cannot contain embedded credential")
	}

	return nil

	}





