package services

import (
	"context"
	"log"
	"time"
	"url-shortener/models"
	"url-shortener/repository"
	"url-shortener/utils"
)

type URLService struct{
	repo *repository.URLRepository
	cacheRepo *repository.CacheRepository
}


func NewURLService(repo *repository.URLRepository, cacheRepo *repository.CacheRepository) *URLService  {
	return &URLService{
		repo: repo,
		cacheRepo: cacheRepo,
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
			oldShortCode := existingURL.ShortCode
			if err != nil{
				return nil,err
			}
			err = s.cacheRepo.Delete(ctx , "url:"+oldShortCode)
			if err != nil {
				log.Printf("Failed to delete the cached URL: %v" , err)
			}else{
				log.Printf("Cached ShortCode Deleted %s" , existingURL.ShortCode)
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
	originalURL , err := s.cacheRepo.Get(ctx , "url:"+shortCode)
	if err == nil {
		log.Printf("Cache hit: %s", shortCode)
		return originalURL , nil
	}
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
	log.Printf("Cache miss: %s" , shortCode)
	remainingTime := time.Until(expiresAt)
	err = s.cacheRepo.Set(ctx , "url:"+shortCode , originalURL , remainingTime)

	if err != nil {
		log.Printf("Failed to cache the URL: %v", err)
	}else{
		log.Printf("URL Cached %s" , shortCode)
	}


	return originalURL , nil
	 
	
}