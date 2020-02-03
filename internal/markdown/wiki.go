package markdown

import (
	"bytes"
	"fmt"

	"github.com/moutend/go-backlog/pkg/types"
)

type WikiFrontMatter struct {
	Name    string `fm:"name"`
	Project string `fm:"project"`
	Created string `fm:"created"`
	Updated string `fm:"updated"`
	Content string `fm:"content"`
}

type Wiki struct {
	Project *types.Project `json:"project"`
	Wiki    *types.Wiki    `json:"wiki"`
}

func (i *Wiki) Marshal() ([]byte, error) {
	buffer := &bytes.Buffer{}

	fmt.Fprintln(buffer, "---")

	if i.Project != nil {
		fmt.Fprintf(buffer, "project: %s\n", i.Project.ProjectKey)
	}
	if i.Wiki == nil {
		goto END_WIKI
	}

	fmt.Fprintf(buffer, "name: %q\n", i.Wiki.Name)

	if i.Wiki.Created != nil {
		fmt.Fprintf(buffer, "created: %s\n", i.Wiki.Created)
	}
	if i.Wiki.Updated != nil {
		fmt.Fprintf(buffer, "updated: %s\n", i.Wiki.Updated)
	}

	fmt.Fprintln(buffer, "---")
	fmt.Fprintln(buffer, i.Wiki.Content)

END_WIKI:

	return buffer.Bytes(), nil
}
