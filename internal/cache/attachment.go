package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/moutend/go-backlog/pkg/types"
)

func SaveWikiAttachment(wiki *types.Wiki, attachment *types.Attachment) error {
	if wiki == nil {
		return fmt.Errorf("cache: can't save nil wiki")
	}
	if attachment == nil {
		return fmt.Errorf("cache: can't save nil attachment")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheWikiAttachmentPath, 0755)

	data, err := json.Marshal(attachment)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheWikiAttachmentPath, fmt.Sprint(wiki.Id), fmt.Sprintf("%d.json", attachment.Id))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func SaveWikiAttachments(wiki *types.Wiki, attachments []*types.Attachment) error {
	for _, attachment := range attachments {
		if err := SaveWikiAttachment(wiki, attachment); err != nil {
			return err
		}
	}

	return nil
}

func LoadWikiAttachments(wikiId uint64) ([]*types.Attachment, error) {
	attachments := []*types.Attachment{}

	err := filepath.Walk(filepath.Join(cacheWikiAttachmentPath, fmt.Sprint(wikiId)), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var attachment *types.Attachment

		if err := json.Unmarshal(data, &attachment); err != nil {
			return err
		}

		attachments = append(attachments, attachment)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return attachments, nil
}
