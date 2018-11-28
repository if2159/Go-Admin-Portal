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


func getStoredUser(username string) (User){
    startConnection()
    rows, err := con.Query("SELECT UID, USERNAME, ROLE_ID, HASHED_PASS FROM USERS WHERE USERNAME = ?", username)
    var user User
    if(rows != nil){
        for rows.Next() {
            var uid string
            var username string
            var roleId int
            var hashedPass []byte
            err = rows.Scan(&uid, &username, &roleId, &hashedPass)
            checkErr(err)
            user = User{uid, username, roleId, hashedPass}
        }
    }
    return user
}

func startConnection(){
    if(con == nil){
        user, password, host := loadDbProperties()
        url := user + ":" + password + "@tcp(" + host + ":3306)/portal?charset=utf8"
        conn, err := sql.Open("mysql", url)
        con = conn
        checkErr(err)
    }
}

func CheckSessionId(username, sessionId string) (bool){
    startConnection()
    rows, err := con.Query("SELECT SESSION_ID, UID FROM portal.SESSIONS WHERE UID=(SELECT UID FROM USERS WHERE USERNAME=?) AND SESSION_ID=?", username, sessionId);
    checkErr(err)

    for rows.Next() {
        var sessionId string
        var uid string
        err = rows.Scan(&sessionId, &uid)
        checkErr(err)

        return true
    }
    return false
}

func updatePassword(){
    startConnection()

    pass := []byte("password123")

   // Hashing the password with the default cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
    checkErr(err)
    fmt.Println("Update Password.")
    fmt.Println(hashedPassword)
    _, errr := con.Query("UPDATE USERS SET HASHED_PASS = ? WHERE USERNAME=?", hashedPassword, "CHRISTIAN")
    checkErr(errr)

}


func comparePassword(username string, password []byte) (bool){
    //updatePassword()

    user := getStoredUser(strings.ToUpper(username))
    storedHash := user.hashedPass
    return (bcrypt.CompareHashAndPassword(storedHash, password) == nil)
}

func createSessionId(username string) (string){
    fmt.Println("Creating Session ID")
    startConnection()

    con.Query("INSERT INTO SESSIONS(UID, SESSION_ID) VALUES((SELECT UID FROM USERS WHERE USERNAME=?), UUID())", username)
    rows, err := con.Query("SELECT SESSION_ID, UID FROM portal.SESSIONS WHERE UID=(SELECT UID FROM USERS WHERE USERNAME=?) AND LOGIN_TIME = (SELECT MAX(LOGIN_TIME) FROM SESSIONS WHERE UID = (SELECT UID FROM USERS WHERE USERNAME=?))", username, username);
    checkErr(err)
    if(rows != nil){
        for rows.Next() {
            var sessionId string
            var uid string
            err = rows.Scan(&sessionId, &uid)
            checkErr(err)

            return sessionId
        }
    }
    return ""
}


func checkErr(err error) {

    if err != nil {
        panic(err)
    }
}

func GetRolesForUser(username string) ([]string){
    var allRoles = GetUserRolesForUser(username)
    var roles []string
    for _, role := range allRoles {
        if(role.HasAccess){
            roles = append(roles, role.RoleName)
        }
    }
    return roles
}

func GetAllRoles()([]string){
    startConnection()
    rows, err := con.Query("SELECT ROLE_NAME FROM ROLES_REF");
    checkErr(err)

    var roles []string

    for rows.Next() {
        var roleName string
        err = rows.Scan(&roleName)
        checkErr(err)

        roles = append(roles, roleName)
    }
    return roles

}

func AdjustUserRoles(username string, roles []string){
    startConnection()
    con.Query("DELETE FROM USER_ROLE_MAP WHERE UID=(SELECT UID FROM USERS WHERE USERNAME=?)", username)
    for _, role := range roles{
        fmt.Println(role)
        _, err := con.Query("INSERT INTO USER_ROLE_MAP(UID, ROLE_ID) VALUES ((SELECT UID FROM USERS WHERE USERNAME=?),(SELECT ID FROM ROLES_REF WHERE ROLE_NAME=?))", username, role)
        checkErr(err)
    }


}

func GetAllLinks()([]string){
    startConnection()
    rows, err := con.Query("SELECT LINK_NAME FROM LINKS_REF");
    checkErr(err)

    var links []string

    for rows.Next() {
        var linkName string
        err = rows.Scan(&linkName)
        checkErr(err)

        links = append(links, linkName)
    }
    return links

}

func AdjustRoleLinks(roleName string, links []string){
    startConnection()
    con.Query("DELETE FROM LINK_ROLE_MAP WHERE ROLE_ID=(SELECT ID FROM ROLES_REF WHERE ROLE_NAME=?)", roleName)
    for _, link := range links{
        fmt.Println(link)
        fmt.Println("Inserting Links")
        _, err := con.Query("INSERT INTO LINK_ROLE_MAP(ROLE_ID, LINK_ID) VALUES ((SELECT ID FROM ROLES_REF WHERE ROLE_NAME=?), (SELECT ID FROM LINKS_REF WHERE LINK_NAME=?))", roleName, link)
        checkErr(err)
    }
}

func GetRoleLinkForRole(roleName string) ([]RoleLink) {
    startConnection()
    rows, err := con.Query("SELECT links.ID, LINK_NAME, ROLE_ID FROM LINKS_REF links " +
	                           "LEFT JOIN (SELECT ROLE_ID, LINK_ID FROM LINK_ROLE_MAP WHERE " +
                               "ROLE_ID=(SELECT ID FROM ROLES_REF WHERE ROLE_NAME=?)) lrm " +
                               "ON lrm.LINK_ID=links.ID", roleName);
    checkErr(err)

    var links []RoleLink

    for rows.Next() {
        var id int
        var linkName string
        var roleId sql.NullInt64
        err = rows.Scan(&id, &linkName, &roleId)
        checkErr(err)
        fmt.Println(roleName)
        link := RoleLink{
    			LinkName: linkName,
                HasAccess: roleId.Valid,
            }
        links = append(links, link)
    }
    return links

}



func GetUserRolesForUser(username string) ([]UserRole){
    startConnection()
    rows, err := con.Query("SELECT ID, ROLE_NAME, UID FROM ROLES_REF " +
	                       "LEFT JOIN " +
                           "(SELECT UID, ROLE_ID FROM USER_ROLE_MAP WHERE " +
                           "UID=(SELECT UID FROM USERS WHERE USERNAME=?)) urm " +
                           "ON urm.ROLE_ID=ID", username);
    checkErr(err)

    var roles []UserRole

    for rows.Next() {
        var id int
        var roleName string
        var uid sql.NullString
        err = rows.Scan(&id, &roleName, &uid)
        checkErr(err)
        role := UserRole{
    			RoleName: roleName,
                HasAccess: uid.Valid,
            }
        roles = append(roles, role)
    }
    return roles

}
