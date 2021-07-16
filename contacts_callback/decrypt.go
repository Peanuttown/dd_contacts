package contacts_callback

import(
    "encoding/base64"
    "encoding/json"
    "crypto/aes"
    // "encoding/hex"
)



func decryptDingDingRequest(
    encrypted []byte,
    aesKey []byte,
)([]byte,error){
    b64Decoded,err := base64.StdEncoding.DecodeString(string(encrypted))
    if err != nil{
        return nil,err
    }
    c,err:= aes.NewCipher([]byte(aesKey))
    if err != nil{
        return nil,err
    }
    decryptBytes := make([]byte,len(b64Decoded))
    c.Decrypt(decryptBytes,b64Decoded)
    return decryptBytes,nil
}

func UnmarshalRequest(encrypted []byte,aesKey []byte)(map[string]interface{},error){
    decryptedByes, err := decryptDingDingRequest(encrypted,aesKey)
    if err != nil{
        return  nil,err
    }
    ret := make(map[string]interface{})
    err = json.Unmarshal(decryptedByes,&ret)
    if err != nil{
        return nil,err
    }
    return ret,nil
    
}
