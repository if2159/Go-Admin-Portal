package main

import (
    _ "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
    "strings"
    "fmt"
)

var con *sql.DB

type User struct {
    uid string
    username string
    roleId int
    hashedPass []byte
}


/*
func main() {

    user, password, host := loadDbProperties()
    con, err := sql.Open("mysql", user + ":@" + password + "@tcp(" + host + ":3306)/portal?charset=utf8")
    checkErr(err)



    pass := []byte("MyDarkSecret")

   // Hashing the password with the default cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

    rows, err := con.Query("UPDATE USERS SET HASHED_PASS = ? WHERE USERNAME=?", hashedPassword, "IAN")
    checkErr(err)


    for rows.Next() {
        var uid string
        var username string
        var role_id int
        var hashed_pass string
        err = rows.Scan(&uid, &username, &role_id, &hashed_pass)
        checkErr(err)
        fmt.Println(uid)
        fmt.Println(username)
        fmt.Println(role_id)
        fmt.Println(hashed_pass)
    }

    con.Close()

    fmt.Println("Done with DB")

    pass := []byte("MyDarkSecret")

    // Hashing the password with the default cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(hashedPassword))


    // Comparing the password with the hash
    err = bcrypt.CompareHashAndPassword(hashedPassword, pass)
    pass := []byte("MyDarkSecret1")

    if (comparePassword("IAN", pass)) {
        fmt.Println("Right Password!")
    } else{
        fmt.Println("Wrong password")
    }

}
*/
func getStoredUser(username string) (User){
    startConnection()
    rows, err := con.Query("SELECT UID, USERNAME, ROLE_ID, HASHED_PASS FROM USERS WHERE USERNAME = ?", username)
    var user User
    for rows.Next() {
        var uid string
        var username string
        var roleId int
        var hashedPass []byte
        err = rows.Scan(&uid, &username, &roleId, &hashedPass)
        checkErr(err)
        user = User{uid, username, roleId, hashedPass}
    }

    return user
}

func startConnection(){
    fmt.Println("start")
    if(con == nil){
        user, password, host := loadDbProperties()

        conn, err := sql.Open("mysql", user + ":@" + password + "@tcp(" + host + ":3306)/portal?charset=utf8")
        con = conn
        checkErr(err)
    }
    fmt.Println("startEnd")
}


func updatePassword(){

    pass := []byte("password123")

   // Hashing the password with the default cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

    con.Query("UPDATE USERS SET HASHED_PASS = ? WHERE USERNAME=?", hashedPassword, "CHRISTIAN")
    checkErr(err)
}


func comparePassword(username string, password []byte) (bool){
    user := getStoredUser(strings.ToUpper(username))
    updatePassword()
    storedHash := user.hashedPass
    return (bcrypt.CompareHashAndPassword(storedHash, password) == nil)
}

func createSessionId(username string) (string){
    startConnection()
    fmt.Println("11")

    con.Query("INSERT INTO SESSIONS(UID, SESSION_ID) VALUES((SELECT UID FROM USERS WHERE USERNAME=?), UUID())", username)
    rows, err := con.Query("SELECT SESSION_ID, UID FROM portal.SESSIONS WHERE UID=(SELECT UID FROM USERS WHERE USERNAME=?) AND LOGIN_TIME = (SELECT MAX(LOGIN_TIME) FROM SESSIONS WHERE UID = (SELECT UID FROM USERS WHERE USERNAME=?))", username, username);
    checkErr(err)

    fmt.Println("12")
    for rows.Next() {
        fmt.Println("15")
        var sessionId string
        var uid string
        err = rows.Scan(&sessionId, &uid)
        checkErr(err)
        fmt.Println("sessionId")

        return sessionId
    }
fmt.Println("14")
    return ""
}


func checkErr(err error) {

    if err != nil {
        panic(err)
    }
}
