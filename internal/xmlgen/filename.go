// Package xmlgen generates XML filename from config template.
package xmlgen

import (
	"regexp"
	"strings"

	"github.com/leokr/blogger-xml-exporter/internal/blogger"
	"github.com/leokr/blogger-xml-exporter/internal/config"
)

// unsafeFilenameChars matches characters unsafe in filenames across platforms.
var unsafeFilenameChars = regexp.MustCompile(`[\x00-\x1f/\\:*?"<>|]`)

// Filename computes the output filename from xmlCfg.Filename (a template with access to form values and {{source "..."}}).
// Returns a sanitized filename with ".xml" extension, always valid and safe.
func Filename(xmlCfg config.XMLConfig, post map[string]interface{}, values map[string]interface{}) string {
	name, err := blogger.RenderTemplate(post, "filename", xmlCfg.Filename, values)
	if err != nil {
		name = ""
	}
	return sanitizeFilename(name) + ".xml"
}

// sanitizeFilename removes unsafe characters, trims edges, ensures non-empty result.
// Returns "post" as fallback for empty/invalid inputs. Max length: 150 chars.
func sanitizeFilename(s string) string {
	s = unsafeFilenameChars.ReplaceAllString(s, "_")
	s = strings.Trim(s, " ._-")
	if s == "" || strings.Trim(s, ".") == "" {
		return "post"
	}
	const maxLen = 150
	if len(s) > maxLen {
		s = strings.TrimRight(s[:maxLen], " ._-")
		if s == "" {
			return "post"
		}
	}
	return s
}
