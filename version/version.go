package version

var (
	gitTag string
)

// Version returns the version number
func Version() string {
	if len(gitTag) == 0 {
		return "dev"
	}
	return gitTag
}