package repository

import (
	"gorm.io/gorm"
)

type DbRepository struct {
	db *gorm.DB
}

// Tạo constructor
func NewDbRepository(db *gorm.DB) *DbRepository {
	return &DbRepository{db: db}
}
