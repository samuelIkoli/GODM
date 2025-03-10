package models

type User struct{
	Email string
	Hash string 
	Salt string 
	Username string 
}


// type User struct{
// 	Email string `"json": "email"`
// 	Hash string `"json": "hash"`
// 	Salt string `"json": "salt"`
// 	Username string `"json": "username"`
// }