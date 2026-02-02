package lg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func P(data any) {
	strLog, _ := json.MarshalIndent(data, "", "\t")
	if IsAWS() {
		strLog, _ = json.Marshal(data)
	}
	fmt.Println(string(strLog))
}

func IsAWS() bool {
	env := os.Getenv("AWS_EXECUTION_ENV")
	return strings.Contains(env, "ECS") || strings.Contains(env, "EC2") || os.Getenv("ECS_CONTAINER_METADATA_URI_V4") != ""
}
