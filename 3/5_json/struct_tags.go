package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int `json:"user_id,omitempty,string"`
	Username string
	Address  string `json:",omitempty"`
	Comnpany string `json:"-"`
}

func main() {
	u := &User{
		ID:       1,
		Username: "rvasily",
		Address:  "",
		Comnpany: "Mail.Ru Group",
	}
	result, _ := json.Marshal(u)
	fmt.Printf("json string: %s\n", string(result))
}
