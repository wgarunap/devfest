package models

import "testing"

func TestDialNumber(t *testing.T) {
	p := Person{
		ID:        1,
		Firstname: "test fn",
		Lastname:  "test ln",
		Contactinfo: Contactinfo{
			City:    "Colombo",
			Zipcode: "+94",
			Phone:   717180171,
		},
	}

	dn := p.DialNumber()

	if dn != "+94717180171" {
		t.Error("Unexpected result ", dn)
	}
}

func TestDIalNumberTableDriven(t *testing.T) {

	type TestSource struct {
		Input  Person
		Output string
	}

	p := Person{
		ID:        1,
		Firstname: "test fn",
		Lastname:  "test ln",
		Contactinfo: Contactinfo{
			City:    "Colombo",
			Zipcode: "+94",
			Phone:   717180171,
		},
	}

	p1 := p
	p1.Zipcode = "+36"
	p1.Phone = 717180008

	testData := []TestSource{
		TestSource{
			Input:  p,
			Output: "+94717180171",
		},
		TestSource{
			Input:  p1,
			Output: "+36717180008",
		},
	}

	for _, td := range testData {
		dn := td.Input.DialNumber()

		if dn != td.Output {
			t.Error("Unexpected result ", dn)
		}
	}
}

func BenchmarkDialNumber(b *testing.B) {

	p := Person{
		ID:        1,
		Firstname: "test fn",
		Lastname:  "test ln",
		Contactinfo: Contactinfo{
			City:    "Colombo",
			Zipcode: "+94",
			Phone:   717180171,
		},
	}

	for n := 0; n < b.N; n++ {
		p.DialNumber()
	}
}
