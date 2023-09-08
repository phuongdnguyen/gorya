package utils

import "fmt"

func GetComputeFilter(key, val string) string {
	return fmt.Sprintf("labels.%s=%s", key, val)
}

func GetCloudSqlFilter(key, val string) string {
	return fmt.Sprintf("settings.userLabels.%s=%s", key, val)
}
