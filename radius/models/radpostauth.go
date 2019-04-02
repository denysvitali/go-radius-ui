package models

import (
	"encoding/json"
	"time"
)

type RadPostAuth struct {
	id       int32
	Username string `json:"username"`
	pass     string
	Reply    string    `json:"reply"`
	AuthDate time.Time `json:"authdate"`
}

func (u *RadPostAuth) MarshalJSON() ([]byte, error) {
	type Alias RadPostAuth
	return json.Marshal(&struct {
		AuthDate int64 `json:"authdate"`
		AuthDateString string `json:"authdate_formatted"`
		*Alias
	}{
		AuthDate: u.AuthDate.Unix(),
		AuthDateString: u.AuthDate.Format(time.UnixDate),
		Alias:    (*Alias)(u),
	})
}

func (r *RadPostAuth) SetId(id int32) {
	r.id = id
}
func (r *RadPostAuth) SetPass(pass string) {
	r.pass = pass
}
