package db

import "fmt"

func BuildDSN(username, password, dbname, host, port string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", username, password, dbname, host, port)
}
