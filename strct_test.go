package kstrct

import (
	"fmt"
	"testing"
	"time"
)

type User struct {
	Id       uint
	FloatNum float64
	Username string
	IsAdmin  bool
	Created  time.Time
	List     []string
	Db       Database
}

type Database struct {
	DSN string
}

func BenchmarkFillSelected(b *testing.B) {
	u := User{}
	temps := time.Now()
	db := Database{
		DSN: "testdsn",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := FillFromSelected(&u, "float_num,username,is_admin,created,list,db", 3.24, "kamal", true, temps, "hello,bye", []any{"testdsn"})
		if err != nil {
			b.Error(err)
		}
	}
	if u.FloatNum != 3.24 || u.Username != "kamal" || u.Created != temps || u.Db != db {
		b.Error("failed")
	}
}

func BenchmarkFillValues(b *testing.B) {
	u := User{}
	temps := time.Now()
	db := Database{
		DSN: "testdsn",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := FillFromValues(&u, 3.24, "kamal", true, temps, "hello,bye", []any{"testdsn"})
		if err != nil {
			b.Error(err)
		}
	}
	if u.FloatNum != 3.24 || u.Username != "kamal" || u.Created != temps || u.Db != db {
		b.Error("failed")
	}
}

func BenchmarkFillFromMap(b *testing.B) {
	temps := time.Now()
	db := Database{
		DSN: "testdsn",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := User{}
		err := FillFromMap(&u, map[string]any{
			"float_num": 3.24,
			"created":   temps,
			"list":      "hello,bye",
			"db":        []any{"testdsn"},
		})
		if err != nil {
			b.Error(err)
		}
		if u.FloatNum != 3.24 || len(u.List) != 2 || u.Created != temps || u.Db != db {
			b.Log(u)
			b.Error("failed")
		}
	}
}

type KormUser struct {
	Id        int       `json:"id,omitempty"`
	Uuid      string    `json:"uuid,omitempty" korm:"size:40;iunique"`
	Email     string    `json:"email,omitempty" korm:"size:50;iunique"`
	Password  string    `json:"password,omitempty" korm:"size:150"`
	IsAdmin   bool      `json:"is_admin,omitempty" korm:"default:false"`
	Image     string    `json:"image,omitempty" korm:"size:100;default:''"`
	CreatedAt time.Time `json:"created_at,omitempty" korm:"now"`
}

func BenchmarkRange(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := KormUser{
			Id:        1,
			Uuid:      "some uuid",
			Email:     "email here",
			Password:  "passhere",
			IsAdmin:   false,
			CreatedAt: time.Now(),
		}
		u = Range(&u, func(fCtx FieldCtx) {
			if fCtx.Name == "password" {
				fCtx.Field.SetString("new something")
			}
		}, "korm")
	}
}

func BenchmarkFillFromMapS(b *testing.B) {
	temps := time.Now()
	db := Database{
		DSN: "testdsn",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u, err := FillFromMapS[User](map[string]any{
			"float_num": 3.24,
			"created":   temps,
			"list":      "hello,bye",
			"db":        []any{"testdsn"},
		})
		if err != nil {
			b.Error(err)
		}
		if u.FloatNum != 3.24 || len(u.List) != 2 || u.Created != temps || u.Db != db {
			b.Log(u)
			b.Error("failed")
		}
	}
}

func TestFillFromSelected(t *testing.T) {
	u := User{}
	temps := time.Now()
	db := Database{
		DSN: "testdsn",
	}
	err := FillFromSelected(&u, "float_num,username,created,db", 3.24, "kamal", temps, []any{"testdsn"})
	if err != nil {
		t.Error(err)
	}
	if u.FloatNum != 3.24 || u.Username != "kamal" || u.Created != temps || u.Db != db {
		t.Error("failed", u)
	}
}

func TestFillFromValues(t *testing.T) {
	u := User{}
	temps := time.Now()
	db := Database{
		DSN: "testdsn",
	}
	err := FillFromValues(&u, 3.24, "kamal", true, temps, "hello,byee", []any{"testdsn"})
	if err != nil {
		t.Error(err)
	}
	if u.FloatNum != 3.24 || u.Username != "kamal" || u.Created != temps || u.Db != db {
		t.Error("failed", u)
	}
}

type Custom struct {
	Id         *int
	Admin      *bool
	Email      *string
	FieldTime  time.Time
	FieldTime2 time.Time
	FieldTime3 time.Time
}

func TestFillFromMap(t *testing.T) {
	u := Custom{}
	err := FillFromMap(&u, map[string]any{
		"id":          nil,
		"admin":       nil,
		"email":       nil,
		"field_time":  time.Now(),
		"field_time2": time.Now().Unix(),
		"field_time3": "2023-01-06 23:08",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("--------------------")
	fmt.Println(u)
}
