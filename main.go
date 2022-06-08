package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main() {
	// 此处192.168.67.92是容器所在主机的ip，1234是容器在主机上的映射端口
	// docker run --name mysql -d -p 1234:3306 -e MYSQL_ROOT_PASSWORD=root docker.io/mysql:latest
	db, err := sql.Open("mysql", "root:root@tcp(192.168.67.92:1234)/gosql")
	if err != nil {
		log.Panic(err)
	}
	// 使用Ping判断连接是否正常
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	defer db.Close()
	// 新建table
	query := `
	CREATE TABLE users (
	   id INT AUTO_INCREMENT,
	   username TEXT NOT NULL,
	   password TEXT NOT NULL,
	   created_at DATETIME,
	   PRIMARY KEY (id)
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Panic(err)
	}
	// 插入两个user记录
	_, err = db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, "foo", "foo", time.Now())
	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, "bar", "bar", time.Now())
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("最后插入的id是：%d\n", lastInsertId)
	fmt.Printf("insert后影响的行数：%d\n", rowsAffected)
	// 多行查询
	rows, _ := db.Query("select id, username from users")
	for rows.Next() {
		var name string
		var id int
		rows.Scan(&id, &name)
		fmt.Printf("id=%d的用户名是：%s\n", id, name)
	}
	// 单行指定查询，获取id=1的username
	var nameSpecified string
	db.QueryRow("select username from users where id= ?", 1).Scan(&nameSpecified)
	fmt.Printf("指定id=1获取的用户名是：%s\n", nameSpecified)
}
