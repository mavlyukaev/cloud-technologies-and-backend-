package mutex

type Mutex struct {
	Count  int
	signal chan struct{}
}

func New(count int) *Mutex {
	return &Mutex{
		Count:  count,
		signal: make(chan struct{}, count),
	}
}

func (m *Mutex) Unlock() {
	m.signal <- struct{}{}
}

func (m *Mutex) Wait() {
	for i := 0; i < m.Count; i++ {
		<-m.signal
	}
}