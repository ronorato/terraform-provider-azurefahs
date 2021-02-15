package helpers

import (
	"fmt"
	"log"
)

func DebugLog(f string, v ...interface{}) {
	// if os.Getenv("TF_LOG") == "" {
	// 	return
	// }

	// if os.Getenv("TF_ACC") != "" {
	// 	return
	// }

	log.Printf("[STG-DEBUG] %s", fmt.Sprintf(f, v...))
}
