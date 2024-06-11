package seedwork

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func PrintIndentedLog(data bytes.Buffer) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data.Bytes(), "", "  "); err != nil {
		log.Fatalf("error indenting JSON: %v", err)
	}
	fmt.Println(prettyJSON.String())
}
