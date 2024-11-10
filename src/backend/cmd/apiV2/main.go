package main

import (
	"fmt"

	"github.com/sachatarba/course-db/internal/api/v2"
)

func main() {
	api := v2.ApiServer{}
	fmt.Println("kek")
	api.Run()
}
