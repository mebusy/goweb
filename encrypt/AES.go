package encrypt

import (
    "crypto/aes"
    "crypto/cipher"
)

// golang crypt包的AES加密函数的使用


// GCM实现算法不需要pad。
func AES_GCM_Encrypt(plaintext, key, nonce []byte) []byte {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err.Error())
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        panic(err.Error())
    }
    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
    return ciphertext
}
func AES_GCM_Decrypt(ciphertext, key, nonce []byte) []byte {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err.Error())
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        panic(err.Error())
    }
    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        panic(err.Error())
    }
    return plaintext
}

