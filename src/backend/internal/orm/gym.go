package orm

import "github.com/google/uuid"

type Gym struct {
	ID              uuid.UUID  `gorm:"type:UUID;primaryKey"`
	Name            string     `gorm:"type:TEXT"`
	Phone           string     `gorm:"type:TEXT;check:check_valid_international_phone, (phone ~ '^\\+[0-9]+-[0-9]+-[0-9]+-[0-9]+-[0-9]+')"`
	City            string     `gorm:"type:TEXT"`
	Addres          string     `gorm:"type:TEXT"`
	IsChain         string     `gorm:"type:BOOLEAN"`
	Trainers        []*Trainer `gorm:"many2many:gym_trainers"`
	Equipments      []Equipment
	MembershipTypes []MembershipType
}
  