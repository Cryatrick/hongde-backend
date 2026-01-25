package service

import (
    "fmt"
    "time"
    "encoding/hex"
    // "encoding/hmac"
    "crypto/hmac"
    "crypto/sha256"
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
    // hostUrl := config.BaseUrl
    // return fmt.Sprintf("%s/v1/manajemen-soal/soal/%s",hostUrl,filename)
    // return fmt.Sprintf("%s/web/soal_images/%s", hostUrl, filename)

    exp := time.Now().Add(24 * time.Hour).Unix()
    secret := []byte(config.ENCRYPTION_KEY)

    payload := fmt.Sprintf("%s:%d", filename, exp)

    mac := hmac.New(sha256.New, secret)
    mac.Write([]byte(payload))
    sig := hex.EncodeToString(mac.Sum(nil))

    return fmt.Sprintf(
        "%s/v1/get-image/%s?exp=%d&sig=%s",
        config.BaseUrl,
        filename,
        exp,
        sig,
    )
}