package kstrct

import (
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
	u := User{}
	temps := time.Now()
	db := Database{
		DSN: "testdsn",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := FillFromMap(&u, map[string]any{
			"float_num": 3.24,
			"created":   temps,
			"list":      "hello,bye",
			"db":        []any{"testdsn"},
		})
		if err != nil {
			b.Error(err)
		}
	}
	if u.FloatNum != 3.24 || len(u.List) != 2 || u.Created != temps || u.Db != db {
		b.Log(u)
		b.Error("failed")
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
