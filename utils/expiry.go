package utils

import "time"

func GenerateExpiry() time.Time {
    return time.Now().Add(30 * 24 * time.Hour)
}



