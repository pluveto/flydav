package app

import (
	"fmt"

	"example.com/m/pkg/logger"
)

func Run(args Args, conf Conf) {
	logger.Warn("initialize this project")
	if args.Seperately {
		for _, name := range args.Names {
			fmt.Print("Hi ", name)
		}
	} else {
		fmt.Print("Hi ")
		for i, name := range args.Names {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(name)
		}
	}
	fmt.Println("!")
}
