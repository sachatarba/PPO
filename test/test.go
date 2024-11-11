package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/lib/pq"
// )

// func main() {
// 	// Настройки подключения
// 	connStr := "host=172.17.0.1 port=5432 user=postgres password=postgres sslmode=disable"

// 	// Открываем соединение с базой данных
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Ошибка подключения к базе данных:", err)
// 	}
// 	defer db.Close()

// 	// Проверка соединения
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal("Не удалось подключиться к базе данных:", err)
// 	}

// 	fmt.Println("Успешное подключение к базе данных!")
// }
