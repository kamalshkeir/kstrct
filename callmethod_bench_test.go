package kstrct

import (
	"testing"
)

type BenchPerson struct {
	Name string
	Age  int
}

func (p *BenchPerson) Calculate(x, y int) (int, string) {
	return x + y, p.Name
}

func BenchmarkMethodCalls(b *testing.B) {
	person := &BenchPerson{Name: "Bob", Age: 25}
	calculate := func(p *BenchPerson, x, y int) (sum int, message string) {
		return x + y, p.Name
	}

	// Add the method using kstrct
	err := AddMethod(person, "Calculate", calculate)
	if err != nil {
		b.Fatalf("Error adding method: %v", err)
	}

	b.Run("Direct Method Call", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = person.Calculate(5, 3)
		}
	})

	b.Run("Kstrct CallMethod", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := CallMethod(person, "Calculate", 5, 3)
			if err != nil {
				b.Fatalf("Error calling method: %v", err)
			}
		}
	})
}
