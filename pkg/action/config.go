package action

type Config struct {
	GithubPages struct {
		Branch    string `fig:"branch" default:"gh-pages"`
		Directory string `fig:"directory" default:"chronos"`
	} `fig:"github_pages"`

	Storage struct {
		// MeasurementsCapacity specifies how many latests measurements
		// (or commits) for one individual series will be preserved.
		MeasurementsCapacity int `fig:"measurements_capacity" default:"25"`
	} `fig:"storage"`

	Units []struct {
		// Name stands for metric's name.
		// For example `ns/op`, `B/op` and others.
		Name string `fig:"name" validate:"required"`
		// Better must be one of `higher` or `lower`.
		Better string `fig:"better" validate:"required"`
		// Threshold is a value in (0; 1).
		Threshold float64 `fig:"threshold" default:"0.05"`
	} `fig:"units"`
}
