package builder

import (
	"time"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/go-faker/faker/v4"
)

type ClientBuilder struct {
	client entity.Client
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		client: entity.Client{
			ID:                uuid.New(),
			Login:             faker.FirstName(),
			Password:          "password123",
			Fullname:          "Test User",
			Email:             "test@example.com",
			Phone:             "+7-999-999-99-99",
			Birthdate:         time.Now().AddDate(-25, 0, 0).Format("2006-01-02"),
			ClientMemberships: []entity.ClientMembership{},
			Schedules:         []entity.Schedule{},
		},
	}
}

func (b *ClientBuilder) SetID(id uuid.UUID) *ClientBuilder {
	b.client.ID = id
	return b
}

func (b *ClientBuilder) SetLogin(login string) *ClientBuilder {
	b.client.Login = login
	return b
}

func (b *ClientBuilder) SetPassword(password string) *ClientBuilder {
	b.client.Password = password
	return b
}

func (b *ClientBuilder) SetFullname(fullname string) *ClientBuilder {
	b.client.Fullname = fullname
	return b
}

func (b *ClientBuilder) SetEmail(email string) *ClientBuilder {
	b.client.Email = email
	return b
}

func (b *ClientBuilder) SetPhone(phone string) *ClientBuilder {
	b.client.Phone = phone
	return b
}

func (b *ClientBuilder) SetBirthdate(birthdate string) *ClientBuilder {
	b.client.Birthdate = birthdate
	return b
}

func (b *ClientBuilder) AddClientMembership(membership entity.ClientMembership) *ClientBuilder {
	b.client.ClientMemberships = append(b.client.ClientMemberships, membership)
	return b
}

func (b *ClientBuilder) AddSchedule(schedule entity.Schedule) *ClientBuilder {
	b.client.Schedules = append(b.client.Schedules, schedule)
	return b
}

func (b *ClientBuilder) Invalid() *ClientBuilder {
	b.client.ID = uuid.Nil              
	b.client.Login = ""                 
	b.client.Password = ""              
	b.client.Fullname = ""              
	b.client.Email = "invalid-email"    
	b.client.Phone = "invalid-phone"    
	b.client.Birthdate = "invalid-date" 

	return b
}

func (b *ClientBuilder) Build() entity.Client {
	return b.client
}
