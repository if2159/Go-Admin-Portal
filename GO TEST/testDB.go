package main

import (
    //"net/http"
    //"path/filepath"
    _ "github.com/go-sql-driver/mysql"
    //"golang.org/x/crypto/bcrypt"
    "database/sql"
    //"strings"
    "fmt"
)

var con *sql.DB

func main() {
    fmt.Println("Server Starting...")
	if(con == nil){
        conn, err := sql.Open("mysql", "hydrogen:@Catscatscats1@tcp(51.15.113.23:3306)/portal?charset=utf8")
        con = conn
        checkErr(err)
    }
	
	rows, err := con.Query("SELECT * FROM portal.USERS");
	checkErr(err)
	
	cols,err := rows.Columns();
	checkErr(err)
	fmt.Println(cols)
	
	for rows.Next() {
        var s1,s2,s3,s4 string
        err = rows.Scan(&s1,&s2,&s3,&s4)
        checkErr(err)
		fmt.Println(s1+s2+s3+s4)
    }
	
	
	
	
}

func checkErr(err error) {

    if err != nil {
        panic(err)
    }
}