package models

import "fmt"

type Person struct {
	ID          int    `json:"id,omitempty"`
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	ContactInfo `json:"contactinfo,omitempty"`
}

type ContactInfo struct {
	City     string `json:"city,omitempty"`
	AreaCode string `json:"areacode,omitempty"`
	Phone    int    `json:"phone,omitempty"`
}

// GetDialNumber retrns the dialing number
func (p *Person) GetDialNumber() string {
	return p.AreaCode + fmt.Sprintf("%d", p.Phone)
}
