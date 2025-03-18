package main

import (
	"fmt"

	"github.com/whatsmynameagain/go-blog-aggregator/internal/config"
)

func main() {
	conf, _ := config.Read()
	conf.SetUser("Potato")
	fmt.Printf("%v\n", conf.DBUrl)
	fmt.Printf("%v\n", conf.CurrentUserName)
}
