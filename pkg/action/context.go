package action

type Input struct {
	Token      string
	Owner      string
	Repository string
	CommitHash string

	BenchmarksFilepath string
	BranchStorage      string
}
