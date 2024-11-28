package main

import (
	"fmt"
	"testing"
)

// Тест для проверки корректной работы кастомного мьютекса.
func TestMutex(t *testing.T) {
	m := Mutex{
		Count:  5,
		signal: make(chan struct{}, 5),
	}

	completed := 0
	const goroutines = 5

	for i := 0; i < goroutines; i++ {
		go func() {
			defer m.Unlock()
			completed++
		}()
	}

	m.Wait()

	if completed != goroutines {
		t.Errorf("Ожидалось %d завершённых горутин, но получили %d", goroutines, completed)
	}
}

// Бенчмарк для оценки производительности кастомного мьютекса.
func BenchmarkMutex(b *testing.B) {
	for n := 1; n <= 3; n++ { // Проверим для 1, 10 и 100 горутин
		b.Run(fmt.Sprintf("goroutines=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := Mutex{
					Count:  n,
					signal: make(chan struct{}, n),
				}

				for j := 0; j < n; j++ {
					go func() {
						defer m.Unlock()
					}()
				}

				m.Wait()
			}
		})
	}
}
