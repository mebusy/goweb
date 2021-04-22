package iap

import (
    "crypto/sha1"
    "encoding/json"
    "encoding/hex"
    "log"
    "io"
    "fmt"
    "bytes"
)
func VerifyIAPReceipt( receiptB64 string, hex_uuid string, bundleID, prodcutID string ) bool {
    retjson, err := ParseIAPReceiptBer( receiptB64 )
    if err != nil {
        log.Println(err)
        return false
    }

    var all_iap_receipts []IAPReceipt_t
    err = json.Unmarshal( []byte(retjson), &all_iap_receipts )
    if err != nil {
        log.Println( "retjson Unmarshal failed:", err)
        return false
    }

    for _, receipt := range all_iap_receipts {
        // check bundle id
        if bundleID != receipt.BundleID {
            log.Printf( "bundleID %s not fit this app", receipt.BundleID )
            continue
        }
        // check productID
        if prodcutID != receipt.Receipt.ProductID {
            log.Printf( "unmatched productID:%s", receipt.Receipt.ProductID )
            continue
        }

        h := sha1.New()
        //*
        d, err := hex.DecodeString( hex_uuid + receipt.OpaqueHex + receipt.BundleIDHex )
        /*/
        d, err := hex.DecodeString( `c68bce287e27494aa082acecb932817157FFAEFF040000000C18636F6D2E74656D706F726172792E69706174657374617070` )
        //*/
        if err !=nil {
            log.Println(err)
            continue
        }
        io.Copy( h, bytes.NewReader(d) )

        // log.Printf("%x", h.Sum(nil))
        sha1_calc := fmt.Sprintf( "%x", h.Sum(nil) )
        log.Printf( "sha1_calc:%s -- sha1_receipt:%s", sha1_calc, receipt.SHA1Hex  )
        if sha1_calc == receipt.SHA1Hex {
            return true
        }
    }

    return false
}

