package encrypt

import (
	"encoding/base64"
	"encoding/asn1"
	"log"
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
        } `asn1:"set"`
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
    for _ , data_block := range v.Data {
        for _, receiptdata :=range data_block.ContentSeq43.ReceiptData {
            log.Println( receiptdata )
            var payload []ReceiptAttribute_t
            rest, err := asn1.UnmarshalWithParams( receiptdata, &payload, "set" )
            if err != nil {
                log.Println( err )
                continue
            }
            log.Println("rest:", rest)
            log.Printf( "%+v", payload )
        }
    }

    return "" , nil
}
