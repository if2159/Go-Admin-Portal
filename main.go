package main

import(
    "fmt"
)

func main() {
	 props,err:= ReadPropertiesFile("db.properties")
	//fmt.Printf("hello, world\n")
	if(err!=nil||err==nil){
		fmt.Println(props["host"])
		fmt.Println(props["db_user"])
		fmt.Println(props["db_password"])
	}

}
