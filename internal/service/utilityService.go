package service

import (
    "fmt"
    "math/rand"

    "hongde_backend/internal/config"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func BuildImageURL(filename string) string {
    hostUrl := config.BaseUrl
    return fmt.Sprintf("%s/web/soal_images/%s", hostUrl, filename)
}