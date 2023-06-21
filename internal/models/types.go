package models

import "time"

type Config struct {
	Host string
	Port string
}

type HumoPayment struct {
	Id        int        `json:"-" gorm:"column:id; autoIncrement;primaryKey"`
	AgentId   int        `json:"agent_id" gorm:"column:agent_id"`
	Amount    int        `json:"amount" gorm:"column:amount; not nul"`
	Currency  string     `json:"currency" gorm:"column:currency; not nul"`
	Account   string     `json:"account" gorm:"column:account; not nul"`
	Status    string     `json:"status" gorm:"column:status"`
	CreatedAt *time.Time `json:"-" gorm:"column:created_at; autoCreateTime"`
	UpdatedAt *time.Time `json:"-" gorm:"column:updated_at"`
}

type PartnerPayment struct {
	Id        int        `json:"-" gorm:"column:id;autoIncrement; primaryKey"`
	Amount    int        `json:"amount" gorm:"column:amount; not nul"`
	Currency  string     `json:"currency" gorm:"column:currency; not nul"`
	Account   string     `json:"account" gorm:"column:account; not nul"`
	Bank      string     `json:"bank" gorm:"column:bank; not nul"`
	Status    string     `json:"status" gorm:"column:status"`
	CreatedAt *time.Time `json:"-" gorm:"column:created_at; autoCreateTime"`
	UpdatedAt *time.Time `json:"-" gorm:"column:updated_at"`
}

type Output struct {
	Id        int
	AgentId   int
	Amount    int
	Currency  string
	Account   string
	Bank      string
	Status    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Partner struct {
	Id              int    `gorm:"column:id; autoIncrement; primaryKey"`
	Name            string `gorm:"column:name"`
	IntegrationDate string `gorm:"column:integration_data"`
	IntervalType    string `gorm:"column:interval_date"`
}

type Reestr struct {
	Id        int     `gorm:"column:id; autoIncrement; primaryKey"`
	PartnerId int     `gorm:"column:partner_id"`
	FileName  string  `gorm:"column:file_name"`
	Code      int     `gorm:"column:code; default:206"`
	Comment   *string `gorm:"column:comment"`
}
