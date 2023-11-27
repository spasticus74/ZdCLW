package zd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"github.com/rs/zerolog/log"
)

func encrypt(plaintext []byte) ([]byte, error) {
	log.Info().Msg("Calling encrypt")
	key := []byte("Z&27#mKkq!Gqk#saoYL6$BxML^7qmGVc")
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Err(err).Msg("couldn't create cipher")
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Err(err).Msg("couldn't create GCM")
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Err(err).Msg("couldn't read nonce")
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	log.Info().Msg("Calling decrypt")
	key := []byte("Z&27#mKkq!Gqk#saoYL6$BxML^7qmGVc")
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Err(err).Msg("couldn't create cipher")
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Err(err).Msg("couldn't create GCM")
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		err := errors.New("ciphertext too short")
		log.Err(err).Msg("invalid cipher text")
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
