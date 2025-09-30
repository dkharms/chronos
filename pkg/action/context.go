package action

type Context struct {
	Token      string
	Owner      string
	Repository string
	CommitHash string

	InputFilepath string
	BranchStorage string
	BranchPages   string
}
