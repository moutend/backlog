package markdown

import (
	"bytes"
	"fmt"

	"github.com/ericaro/frontmatter"
	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/go-backlog/pkg/types"
)

type WikiFrontMatter struct {
	Project string `fm:"project"`
	Wiki    uint64 `fm:"wiki"`
	Name    string `fm:"name"`
	Content string `fm:"content"`
}

type Wiki struct {
	Project *types.Project `json:"project"`
	Wiki    *types.Wiki    `json:"wiki"`
}

func (v *Wiki) Marshal() ([]byte, error) {
	buffer := &bytes.Buffer{}

	fmt.Fprintln(buffer, "---")

	if v.Project != nil {
		fmt.Fprintf(buffer, "project: %s\n", v.Project.ProjectKey)
	}
	if v.Wiki == nil {
		goto END_WIKI
	}

	fmt.Fprintf(buffer, "wiki: %d\n", v.Wiki.Id)
	fmt.Fprintf(buffer, "name: %q\n", v.Wiki.Name)

	if v.Wiki.Created != nil {
		fmt.Fprintf(buffer, "created: %s\n", v.Wiki.Created)
	}
	if v.Wiki.Updated != nil {
		fmt.Fprintf(buffer, "updated: %s\n", v.Wiki.Updated)
	}

	fmt.Fprintf(buffer, "url: https://%s/alias/wiki/%d\n", backlog.SpaceName(), v.Wiki.Id)
	fmt.Fprintln(buffer, "---")
	fmt.Fprintln(buffer, v.Wiki.Content)

END_WIKI:

	return buffer.Bytes(), nil
}

func (v *Wiki) unmarshal(data []byte) error {
	var fo WikiFrontMatter

	if err := frontmatter.Unmarshal(data, &fo); err != nil {
		return err
	}
	if fo.Project == "" {
		return fmt.Errorf("markdown: project is required")
	}

	project, err := backlog.GetProject(fo.Project)

	if err != nil {
		return err
	}

	v.Project = project

	if fo.Wiki != 0 {
		wiki, err := backlog.GetWiki(fo.Wiki)

		if err != nil {
			return err
		}

		v.Wiki = wiki
	} else {
		v.Wiki = &types.Wiki{}
	}

	v.Wiki.ProjectId = v.Project.Id
	v.Wiki.Name = fo.Name
	v.Wiki.Content = fo.Content

	return nil
}

func (v *Wiki) Unmarshal(data []byte) error {
	return v.unmarshal(data)
}
