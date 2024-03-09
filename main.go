package main

import (
	"clean-architecture/api/router"
	"clean-architecture/bootstrap"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	_gin := gin.Default()
	// Tạo server HTTP tùy chỉnh với thời gian timeout là 0
	// TODO: Chỉ dùng trong môi trường dev
	srv := &http.Server{
		Addr:         ":8080", // Cổng server
		Handler:      _gin,    // Sử dụng router Gin đã tạo
		ReadTimeout:  0,       // Tắt timeout cho thời gian đọc
		WriteTimeout: 0,       // Tắt timeout cho thời gian ghi
	}
	router.SetUp(env, timeout, db, _gin)

	// Khởi động server
	// TODO: chỉ dùng trong môi trường dev
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Xử lý lỗi nếu có
		panic(err)
	}

	err := _gin.Run(env.ServerAddress)
	if err != nil {
		return
	}
}
