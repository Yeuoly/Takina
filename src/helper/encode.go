package helper

//b64
import (
	"encoding/base64"
)

func Base64Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}