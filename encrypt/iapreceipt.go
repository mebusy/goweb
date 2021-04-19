package encrypt

import (
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"log"
	// "strings"
)

/*
ASN.1 INTEGER | int,32,64
ASN.1 BIT STRING | BitString
ASN.1 OCTET STRING | []byte
ASN.1 OBJECT IDENTIFIER  | ObjectIdentifier
ASN.1 ENUMERATED | Enumerated.
ASN.1 UTCTIME or GENERALIZEDTIME | time.Time.
ASN.1 PrintableString, IA5String, or NumericString | string.
Any of the above ASN.1 values | interface{}.

An ASN.1 SEQUENCE OF x or SET OF x can be written to a slice if an x can be written to the slice's element type.
An ASN.1 SEQUENCE or SET can be written to a struct if each of the elements in the sequence can be written to the corresponding element in the struct.

The following tags on struct fields have special meaning to Unmarshal:

application specifies that an APPLICATION tag is used
private     specifies that a PRIVATE tag is used
default:x   sets the default value for optional integer fields (only used if optional is also present)
explicit    specifies that an additional, explicit tag wraps the implicit one
optional    marks the field as ASN.1 OPTIONAL
set         causes a SET, rather than a SEQUENCE type to be expected
tag:x       specifies the ASN.1 tag number; implies ASN.1 CONTEXT SPECIFIC

If the type of the first field of a structure is RawContent then the raw ASN1 contents of the struct will be stored in it.

If the name of a slice type ends with "SET" then it's treated as if the "set" tag was set on it. This results in interpreting the type as a SET OF x rather than a SEQUENCE OF x. This can be used with nested slices where a struct tag cannot be given.
//*/

/*
 There are 4 classes of ASN.1 tags: UNIVERSAL, APPICATION, PRIVATE, and context-specific. The [0] is a context-specific tag
 context-specific的Tag只能出现在SEQUENCE、SET和CHOICE类型的组件中
//*/

type ReceiptAttribute_t struct {  // SEQUENCE
    Type int
    Version int
    Value []byte
}


type ReceiptFile_t  struct {
    SignedData asn1.ObjectIdentifier
    Data []struct {
        I23 int
        Set26 struct {
            // Seq28 struct {
            //     Sha246 asn1.ObjectIdentifier
            //     N41 interface{}
            // }
        } `asn1:"set,optional"`
        ContentSeq43 struct {
            Pkcs7_data asn1.ObjectIdentifier
            ReceiptData [][]byte `asn1:"tag:0"`
        }
        /*
        Certificate383 []struct {
            Seq391 struct {
                A396 []int   `asn1:"tag:0"`
            }
        } `asn1:"tag:0"`
        //*/
        // Signature ?

    } `asn1:"tag:0"`  // whole data
}

type IAPReceipt_t struct {
    BundleID string
    BundleVersion string
    CreationDate string
    ExpirationData string `json:",omitempty"`
    Receipt struct {
        Quantity  int
        ProductID string
        TransactionID string
        OriginTransactionID string  `json:",omitempty"`
        PurchaseData string
        OriginPurchaseData string  `json:",omitempty"`
    }
}


func ans1UnmarshalString( data []byte) (string, error) {
    var val string
    _, err := asn1.Unmarshal( data, &val )
    if err != nil {
        log.Println( err )
        return "", err
    }
    return val, nil
}

func ans1UnmarshalInt( data []byte) (int, error) {
    var val int
    _, err := asn1.Unmarshal( data, &val )
    if err != nil {
        log.Println( err )
        return 0, err
    }
    return val, nil
}

func ParseIAPReceipt( datab64 string ) (string, error) {
    //*
    data, err := base64.StdEncoding.DecodeString( datab64 )
    if err != nil {
        log.Println( err )
        return "" , nil
    }
    _ = data
    //*/

    var v ReceiptFile_t
    rest, err := asn1.Unmarshal( []byte(data), &v )
    if err != nil {
        log.Println( err )
        return "" , nil
    }

    // log.Printf("unmarshaled: %+v", v)
    log.Printf("rest: %+v", string(rest) )
    var all_iap_receipts []IAPReceipt_t
    for _ , data_block := range v.Data {
        for _, receiptdata :=range data_block.ContentSeq43.ReceiptData {
            // log.Println( receiptdata )
            var payload []ReceiptAttribute_t
            rest, err := asn1.UnmarshalWithParams( receiptdata, &payload, "set" )
            if err != nil {
                log.Println( err )
                continue
            }
            // log.Println("rest:", rest)
            // log.Printf( "%+v", payload )
            _ = rest

            var iap IAPReceipt_t
            for _, attr := range payload {
                switch attr.Type {
                    case 2:
                        bid, err := ans1UnmarshalString( attr.Value )
                        if err != nil {
                            continue
                        }
                        // log.Printf( "bid:%s len(%d)", bid, len(bid) )
                        iap.BundleID = bid
                    case 3:
                        bver, err := ans1UnmarshalString( attr.Value )
                        if err != nil {
                            continue
                        }
                        iap.BundleVersion = bver
                    case 12:
                        createdata, err := ans1UnmarshalString( attr.Value )
                        if err != nil {
                            continue
                        }
                        iap.CreationDate = createdata
                    case 21:
                        expireData, err := ans1UnmarshalString( attr.Value )
                        if err != nil {
                            continue
                        }
                        iap.ExpirationData = expireData
                    // case 4: // An opaque value used, with other data, to compute the SHA-1 hash during validation.
                    // case 5: // A SHA-1 hash, used to validate the receipt.
                    case 17: // The receipt for an in-app purchase.
                        var rpld []ReceiptAttribute_t
                        _, err := asn1.UnmarshalWithParams( attr.Value, &rpld, "set" )
                        if err != nil {
                            log.Println( err )
                            continue
                        }
                        // log.Printf("%+v", rpld)
                        for _, rattr  := range rpld {
                            switch rattr.Type {
                            case 1701:
                                quantity, err := ans1UnmarshalInt( rattr.Value )
                                if err != nil {
                                    continue
                                }
                                iap.Receipt.Quantity = quantity
                            case 1702:
                                pid, err := ans1UnmarshalString( rattr.Value )
                                if err != nil {
                                    continue
                                }
                                iap.Receipt.ProductID = pid
                            case 1703:
                                tid, err := ans1UnmarshalString( rattr.Value )
                                if err != nil {
                                    continue
                                }
                                iap.Receipt.TransactionID = tid
                            case 1705:
                                otid, err := ans1UnmarshalString( rattr.Value )
                                if err != nil {
                                    continue
                                }
                                iap.Receipt.OriginTransactionID = otid
                            case 1704:
                                purchaseData, err := ans1UnmarshalString( rattr.Value )
                                if err != nil {
                                    continue
                                }
                                iap.Receipt.PurchaseData = purchaseData
                            case 1706:
                                opurchaseData, err := ans1UnmarshalString( rattr.Value )
                                if err != nil {
                                    continue
                                }
                                iap.Receipt.OriginPurchaseData = opurchaseData
                            }
                        }

                }
            } // payload
            all_iap_receipts = append(all_iap_receipts, iap)
        } // end ContentSeq43.ReceiptData
    } // end v.Data
    // log.Printf( "%+v", all_iap_receipts )

    b, err := json.Marshal( &all_iap_receipts )
    if err != nil {
        log.Println(err)
        return "", err
    }

    ret := string(b)
    // log.Println(ret)
    return ret , nil
}



