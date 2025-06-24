package sys

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func Pr(val ...interface{}) {
	if len(val) > 1 {
		for _, v := range val {
			switch v.(type) {
			case []uint8:
				fmt.Println("[]uint8 ori: ", v)
				fmt.Printf("[]uint8 str: %s-\n\n", v)
				continue
			default:
				bytes, _ := json.MarshalIndent(v, "", "    ")
				fmt.Printf("%T : %s-\n", v, bytes)
			}
		}
	} else {
		bytes, _ := json.MarshalIndent(val, "", "    ")
		fmt.Printf("%T : %s-\n", val, bytes)
	}
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
