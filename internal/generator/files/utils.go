package files

import "fmt"

func EnvVarInDocker(key string) string {
	return fmt.Sprintf("$${%s}", key)
}

func EnvVar(key string) string {
	return fmt.Sprintf("${%s}", key)
}
