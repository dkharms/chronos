package action

type Config struct {
	GithubPages struct {
		Branch    string `fig:"branch" default:"gh-pages"`
		Directory string `fig:"directory" default:"chronos"`
	} `fig:"github-pages"`

	Storage struct {
		// Capacity specifies how many latests measurements
		// (or commits) for one individual series will be preserved.
		Capacity int `fig:"capacity" default:"25"`
	} `fig:"storage"`

	Units []struct {
		// Name stands for metric's name.
		// For example `ns/op`, `B/op` and others.
		Name string `fig:"name" validate:"required"`
		// Better must be one of `higher` or `lower`.
		Better string `fig:"better" validate:"required"`
		// Threshold is a value in (0; 1).
		Threshold float64 `fig:"threshold" default:"0.05"`
		// Reduction is a function which takes a vector as an input
		// and returns scalar value. It must be one of 'min', 'max', 'median', 'mean'.
		Reduction string `fig:"reduction" default:"mean"`
	} `fig:"units"`
}
