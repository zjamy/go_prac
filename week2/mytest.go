package main

import (
	"log"
	"errors"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func TestErrWrap(a int) (string, error) {
	// db 是一个 sql.DB 类型的对象
	// 该对象线程安全，且内部已包含了一个连接池
	// 连接池的选项可以在 sql.DB 的方法中设置，这里为了简单省略了
	db, err := sql.Open("mysql",
		"dog:123456@tcp(localhost:3306)/my_db")
	//db, err := sql.Open("mysql",
	//	"zjuan@tcp(localhost:3306)/my_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		//id int
		name string
	)
	err = db.QueryRow("select name from users where id = ?", a).Scan(&name)
	if err != nil {
		fmt.Println("query error")
		if err == sql.ErrNoRows {
			//log.Println("sql.ErrNoRows")
			return "", fmt.Errorf("%w, data is nil", err)
		}
		log.Fatal(err)
	}

//	defer rows.Close()

	// 必须要把 rows 里的内容读完，或者显式调用 Close() 方法，
	// 否则在 defer 的 rows.Close() 执行之前，连接永远不会释放
//	for rows.Next() {
//	rows.Next()
//	err = rows.Scan(&id, &name)
//	if err != nil {
//		fmt.Println("as expected")
//		if err == sql.ErrNoRows {
//			//log.Println("sql.ErrNoRows")
//			return "", fmt.Errorf("%w, data is nil", err)
//		}
//		log.Fatal(err)
//		return "", err
//	}
//	log.Println(id, name)
////	}
//
//	err = rows.Err()
//	if err != nil {
//		log.Fatal(err)
//
//	}
//
	return "", err
}

func main() {
	_, err := TestErrWrap(71)
//	_, err = TestErrWrap(7)
	err = errors.Unwrap(err)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("we got it.")
	}
}
