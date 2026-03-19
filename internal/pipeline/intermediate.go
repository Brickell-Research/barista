package pipeline

type Guarantee struct {
	Name       string
	Threshold  float64
	WindowDays int
}

type Intermediate struct {
	ServiceName  string
	ProviderName string
	SourceURL    string
	Guarantees   []Guarantee
}
