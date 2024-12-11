package main

import v1 "github.com/sachatarba/course-db/internal/api/v1"

func main() {
	api := v1.ApiServer{}
	api.Run()
}
