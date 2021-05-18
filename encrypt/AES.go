package encrypt

import (
    "crypto/aes"
    "crypto/cipher"
)

// golang crypt包的AES加密函数的使用


// GCM实现算法不需要pad。
func AES_GCM_Encrypt(plaintext, key, nonce []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
    return ciphertext, nil
}
func AES_GCM_Decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }
    return plaintext, nil
}

