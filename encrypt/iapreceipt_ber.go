package encrypt

import (
    "github.com/go-asn1-ber/asn1-ber"
    "log"
    "encoding/base64"
    "encoding/json"
)


func parseReceipt (pket *ber.Packet) []IAPReceipt_t {
    var all_iap_receipts []IAPReceipt_t
    // log.Println( pket.Tag )
    if pket.Tag == 16 {  // SEQUENCE and SEQUENCE OF
        for _, p0 := range pket.Children {
            if p0.Tag == 16 {  // care about only SEQUENCE , ContentSeq43
                // log.Println(  p.Tag )
                for _, d := range p0.Children {
                    // care about array only
                    if d.Tag == 0 {
                        for _,v := range d.Children { // ReceiptData [][]byte `asn1:"tag:0"`
                            // log.Println( i,v.Tag )
                            if v.Tag == 4 {  //  OCTET STRING
                                var receiptdata  []byte = nil
                                for c:=0; c<3; c++ { // protective
                                    v = ber.DecodePacket( v.Data.Bytes() )
                                    // log.Println( v.Tag, v.TagType)
                                    if v.Tag == 4 && v.TagType == 0  {  // TagType==0, primitive
                                        receiptdata = v.Data.Bytes()
                                        break // stop iteract
                                    }
                                } // end for
                                if receiptdata != nil {
                                    iap, err := parseReceiptAttributes( receiptdata )
                                    if err != nil {
                                        continue
                                    }
                                    // log.Println( iap, err )
                                    all_iap_receipts = append(all_iap_receipts, iap)
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    return all_iap_receipts
}


func ParseIAPReceiptBer( datab64 string ) (string, error) {
    //*
    data, err := base64.StdEncoding.DecodeString( datab64 )
    if err != nil {
        log.Println( err )
        return "" , nil
    }
    _ = data
    //*/

    var all_iap_receipts []IAPReceipt_t

    pket := ber.DecodePacket( data )
    if pket.Tag == 16 {  // SEQUENCE and SEQUENCE OF
        for _, p0 := range pket.Children {
            // log.Println( p0.Tag)
            if p0.Tag == 0 {  // Data []struct 
                for _, d := range p0.Children {
                    all_iap_receipts =  append(all_iap_receipts, parseReceipt( d )... )
                }
            }
        }
    }

    b, err := json.Marshal( &all_iap_receipts )
    if err != nil {
        log.Println(err)
        return "", err
    }

    return string(b), nil
}
