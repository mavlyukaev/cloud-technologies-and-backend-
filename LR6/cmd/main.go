package main

import (
	"github.com/mavlyukaev/cloud-technologies-and-backend/internal/worker"
)

func main() {
	worker.RunWorkers(3)
}
