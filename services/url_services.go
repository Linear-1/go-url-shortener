package services

import (
	"context"
	"time"
	"url-shortener/models"
	"url-shortener/repository"
	"url-shortener/utils"
)

type URLService struct{
	repo *repository.URLRepository
}


func NewURLService(repo *repository.URLRepository) *URLService  {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) CreateShortURL(ctx context.Context , OriginalURL string) (*models.URL , error)  {
	existingURL , err := s.repo.FindByOgURL(ctx , OriginalURL)     //existingURL is whole struct with all the feilds
	
	if err != nil {
		return nil , err
	}

	if existingURL != nil {
		if time.Now().Before(existingURL.ExpiresAt){
			return existingURL, nil
		}else{
			shortCode,err := utils.GenerateShortCode(7)
			if err != nil {
				return nil , err
			}

			expiresAt := utils.GenerateExpiry()
			renewalURL := &models.URL{
				ID : existingURL.ID,
				OriginalURL: existingURL.OriginalURL,
				ShortCode: shortCode,
				ExpiresAt: expiresAt,
			}

			renewedURL , err := s.repo.RenewUrl(ctx , renewalURL)

			if err != nil{
				return nil,err
			}
			return renewedURL,nil
		}

	}

	shortCode,err := utils.GenerateShortCode(7)
	if err != nil{
		return nil, err
	}
	expiresAt := utils.GenerateExpiry()

	newURL := &models.URL{
		OriginalURL: OriginalURL,
		ShortCode: shortCode,
		ExpiresAt: expiresAt,
	} 

	createdURL , err := s.repo.CreateURL(ctx , newURL)
	if err != nil {
		return nil,err
	}
	return createdURL,nil	
	
}

func (s *URLService) RedirectURL(ctx context.Context, shortCode string) (string, error) {
	originalURL ,expiresAt, err := s.repo.GetURLByShortCode(ctx,shortCode)
	if err != nil{
		return "",err
	}
	
	if originalURL == ""{
		return "", ErrURLNotFound
	}
	if !time.Now().Before(expiresAt){
		return "",ErrURLExpired
	}
	return originalURL,nil
	 
	
}