package models

import "fmt"

type Person struct {
	ID          int    `json:"id,omitempty"`
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	Contactinfo `json:"contactinfo,omitempty"`
}
type Contactinfo struct {
	City    string `json:"city,omitempty"`
	Zipcode string `json:"zipcode,omitempty"`
	Phone   int    `json:"phone,omitempty"`
}

// DialNumber returns the dialable telephone number
func (p *Person) DialNumber() string {
	return p.Zipcode + fmt.Sprintf("%d", p.Phone)
}
