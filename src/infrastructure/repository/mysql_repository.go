package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"src/entity"
	"strings"
)

// Возвращает все поля типа в массиве
func getFieldValues(c MySQLConvertable) []string {
	s := reflect.ValueOf(&c).Elem()

	var fields []string
	for i := 0; i < s.NumField(); i++ {
		fields = append(fields, fmt.Sprintf("%v", s.Field(i).Interface()))
	}

	return fields
}

func getFieldNames(c MySQLConvertable) []string {
	s := reflect.ValueOf(&c).Elem()

	var fieldNames []string
	for i := 0; i < s.NumField(); i++ {
		fieldNames = append(fieldNames, s.Type().Field(i).Name)
	}

	return fieldNames
}

type MySQLConvertable interface {
	GetTableName() string
	GetID() entity.ID
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQL(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}

func (r *MySQLRepository) Create(c MySQLConvertable) (entity.ID, error) {
	stmt, err := r.db.Prepare(
		fmt.Sprintf(
			`insert into %s (%s) values (%s)`,
			c.GetTableName(),
			strings.Join(getFieldNames(c), ", "),
			strings.Join(getFieldValues(c), ", "),
		),
	)

	if err != nil {
		return c.GetID(), err
	}
}
