package authModel

import (
  "golang.org/x/crypto/bcrypt"
)

func CreateHash(pass string) (string, error){
  hash, err := bcrypt.GenerateFromPassword([]byte(pass),
                                           bcrypt.DefaultCost)

  if err != nil {
    return "", err
  }

  return string(hash), nil
}

func CompareHash(hash string, pass string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
  if err != nil {
    return false
  }
  return true
}

