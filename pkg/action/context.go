package action

type Context struct {
	CommitHash string

	Token      string
	Owner      string
	Repository string

	InputFilepath string
	BranchStorage string
	BranchPages   string
}
