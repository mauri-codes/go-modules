package logger

import (
	"encoding/json"
	"fmt"
)

func Pretty(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshaling: %v", err)
	}
	return string(b)
}
