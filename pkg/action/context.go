package action

type Input struct {
	Token      string
	Owner      string
	Repository string
	CommitHash string

	LanguageTool string

	BenchmarksFilepath string
	BranchStorage      string
}
