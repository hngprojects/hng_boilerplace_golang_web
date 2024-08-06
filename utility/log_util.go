package utility

import (
	"fmt"
)

func LogAndPrint(logger *Logger, data interface{}, args ...interface{}) {
	if len(args) < 1 {
		fmt.Println(data)
		logger.Info(data)
		return
	}
	fmt.Println(data, args)
	logger.Info(data, args)
}
