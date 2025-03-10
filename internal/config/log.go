package config

import "fmt"

func PrintLog(log *Log, args interface{}, data ...interface{}) {

	if len(data) < 1 {
		log.Info(args, data)
		fmt.Println(args)
		return
	}
	fmt.Println(args, data)
	log.Info(args, data)
}
