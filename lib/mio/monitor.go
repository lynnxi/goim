package mio

type Monitor struct {
}

func (m *Monitor) Start(f func(args ...interface{})) {

}

func (m *Monitor) Stop() {

}

func NewMonitor() *Monitor {

	return &Monitor{}
}
