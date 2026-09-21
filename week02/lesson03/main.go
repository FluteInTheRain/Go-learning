package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name      string
	BirthDate time.Time
	Email     string
}

func (p Person) Age() int {
	now := time.Now()
	years := now.Year() - p.BirthDate.Year()
	if now.Month() < p.BirthDate.Month() || (now.Month() == p.BirthDate.Month() && now.Day() < p.BirthDate.Day()) {
		years--
	}
	return years
}

func (p Person) isAdult() bool {
	age := p.Age()
	if age > 18 {
		return true
	}
	return false
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d tuoi)", p.Name, p.Age())
}

func (p Person) Greeting() string {
	return fmt.Sprintf("Xin chao, toi la %s", p.Name)
}

func (p *Person) Rename(newName string) {
	p.Name = newName
}

func ParsePerson(name, dateStr string) (*Person, error) {
	birthDate, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse date %q: %w", dateStr, err)
	}
	return &Person{
		Name:      name,
		BirthDate: birthDate,
	}, nil
}

func main() {
	people := []Person{
		{Name: "An", BirthDate: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		{Name: "Binh", BirthDate: time.Date(1998, 11, 30, 0, 0, 0, 0, time.UTC)},
	}
	// for _, person := range people {

	// 	fmt.Println(person)
	// }
	// for i := range people {
	// 	people[i].Name = "Changed"
	// }
	for _, person := range people {
		fmt.Println(person.Greeting())
		fmt.Printf("Is adult: %t\n", person.isAdult())
	}
	person, err := ParsePerson("An", "22/02/2006")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(person)
}
