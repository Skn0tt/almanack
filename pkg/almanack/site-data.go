package almanack

import "fmt"

const (
	EditorsPicksLoc = "data/editorsPicks.json"
	SidebarLoc      = "data/sidebar.json"
	SiteParamsLoc   = "config/_default/params.json"
)

var messageForLoc = map[string]string{
	EditorsPicksLoc: "Setting homepage configuration",
	SidebarLoc:      "Setting sidebar configuration",
	SiteParamsLoc:   "Setting site parameters",
}

func MessageForLoc(loc string) string {
	msg := messageForLoc[loc]
	if msg == "" {
		return fmt.Sprintf("Updating %s", loc)
	}
	return msg
}
