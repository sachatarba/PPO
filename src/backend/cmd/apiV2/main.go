package main

import "github.com/sachatarba/course-db/internal/api/v2"

func main() {
	api := v2.ApiServer{}
	api.Run()
}
