package encrypt

import (
    "encoding/pem"
    "crypto/sha256"
    "crypto/rsa"
    "crypto"
    "crypto/x509"
    "crypto/rand"
    "log"
    "errors"
    "encoding/base64"
)

/*
    fullPrivateKey:
-----BEGIN xxx PRIVATE KEY-----
...
-----END xxx PRIVATE KEY-----
    bPKCS8:
        PKCS8 | PKCS1
*/
func SignSHA256withRSA( signString,  fullPrivateKey string, bPKCS8 bool ) ( string, error ) {

    block, _ := pem.Decode(  []byte( fullPrivateKey  )   )
    if block == nil {
        return "", errors.New("private key error")
    }

    // combo:  Public/Private + PKCS8/PKCS1
    var pri *rsa.PrivateKey
    if bPKCS8 {
        privateInterface, err  := x509.ParsePKCS8PrivateKey( block.Bytes )
        if err != nil {
            log.Println(err)
            return "", err
        }
        pri = privateInterface.(*rsa.PrivateKey)
    } else {
        privateInterface, err  := x509.ParsePKCS1PrivateKey( block.Bytes )
        if err != nil {
            log.Println(err)
            return "", err
        }
        pri = privateInterface
    }

    h := sha256.New()
    h.Write([]byte( signString  ))
    d := h.Sum(nil)

    signature, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, d)
    if err != nil {
        log.Println(err)
        return "",err
    }
    return base64.StdEncoding.EncodeToString(signature), nil
}

