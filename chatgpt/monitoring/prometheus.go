package monitoring

type Prometheus struct{}

func NewPrometheus() *Prometheus {
    return &Prometheus{}
}

func (p *Prometheus) Export() {
    // Mock implementation
}