package benchmark

type Series struct {
	Name         string        `json:"name"`
	Measurements []Measurement `json:"measurements"`
}

type Measurement struct {
	Name       string   `json:"-"`
	CommitHash string   `json:"commit_hash"`
	Metrics    []Metric `json:"metrics"`
}

type Metric struct {
	Unit   string    `json:"unit"`
	Values []float64 `json:"values"`
}
