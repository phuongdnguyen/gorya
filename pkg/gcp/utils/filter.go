package utils

import "fmt"

func GetFilter(key, val string) string {
	return fmt.Sprintf("labels.%s=%s", key, val)
}
