package main

import (
	"github.com/mavlyukaev/7S-cloud-technologies-and-backend/internal/worker"
)

func main() {
	worker.RunWorkers(3)
}
