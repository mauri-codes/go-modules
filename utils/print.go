package utils

import (
	"encoding/json"
	"fmt"
)

func Pr(data any) {
	res4B, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(res4B))
}
