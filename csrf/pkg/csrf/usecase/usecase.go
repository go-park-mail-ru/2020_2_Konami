package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"konami_backend/csrf/pkg/csrf"
	"time"
)

type CsrfUseCase struct {
	Secret     []byte
	ExpireTime int64
	CsrfRepo   csrf.Repository
}

type TokenMeta struct {
	SessionID string
	TimeStamp int64
}

func NewCsrfUseCase(secret string, expireTime int64, csrfRepo csrf.Repository) (csrf.UseCase, error) {
	key := []byte(secret)
	_, err := aes.NewCipher(key)
	if err != nil {
		return &CsrfUseCase{}, err
	}
	return &CsrfUseCase{Secret: key, ExpireTime: expireTime, CsrfRepo: csrfRepo}, nil
}

func (tk *CsrfUseCase) Create(sid string, timeStamp int64) (string, error) {
	block, err := aes.NewCipher(tk.Secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	td := &TokenMeta{SessionID: sid, TimeStamp: timeStamp}
	data, _ := json.Marshal(td)
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	res := append([]byte(nil), nonce...)
	res = append(res, ciphertext...)

	token := base64.StdEncoding.EncodeToString(res)

	return token, nil
}

func (tk *CsrfUseCase) Check(sid string, inputToken string) (bool, error) {
	block, err := aes.NewCipher(tk.Secret)
	if err != nil {
		return false, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(inputToken)
	if err != nil {
		return false, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return false, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, err
	}

	td := TokenMeta{}
	err = json.Unmarshal(plaintext, &td)
	if err != nil {
		return false, err
	}
	if time.Now().Unix()-td.TimeStamp > tk.ExpireTime {
		return false, csrf.ErrExpiredToken
	}

	expected := TokenMeta{SessionID: sid, TimeStamp: td.TimeStamp}
	err = tk.CsrfRepo.Validate(inputToken)
	if td != expected || err != nil {
		return false, nil
	}
	err = tk.CsrfRepo.Add(inputToken, tk.ExpireTime)
	if err != nil {
		return false, err
	}
	return true, nil
}
