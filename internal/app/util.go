package app

import "fmt"

func toHumanReadable(x int) string {
	if x < 1024 {
		return fmt.Sprintf("%d B", x)
	}
	if x < 1024*1024 {
		return fmt.Sprintf("%d KB", x/1024)
	}
	if x < 1024*1024*1024 {
		return fmt.Sprintf("%d MB", x/1024/1024)
	}

	return fmt.Sprintf("%d GB", x/1024/1024/1024)
}
