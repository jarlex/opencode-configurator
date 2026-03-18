package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jarlex/opencode-configurator/internal/model"
	"gopkg.in/yaml.v3"
)

// skillFrontmatter mirrors the YAML frontmatter in SKILL.md files.
type skillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	License     string `yaml:"license"`
	Metadata    struct {
		Author  string `yaml:"author"`
		Version string `yaml:"version"`
	} `yaml:"metadata"`
}

// ScanSkills scans a directory for skills by reading */SKILL.md files
// and extracting YAML frontmatter. Returns an empty slice (not error)
// if the directory does not exist.
func ScanSkills(dir string) ([]model.Skill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Skill{}, nil
		}
		return nil, err
	}

	var skills []model.Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Skip shared resource directories
		if strings.HasPrefix(entry.Name(), "_") {
			continue
		}

		skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
		skill, err := parseSkillFile(skillPath)
		if err != nil {
			// Skip skills that can't be parsed
			continue
		}
		skills = append(skills, *skill)
	}

	return skills, nil
}

// parseSkillFile reads a SKILL.md file and extracts its YAML frontmatter and content.
func parseSkillFile(path string) (*model.Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := string(data)

	frontmatter, err := extractFrontmatter(raw)
	if err != nil {
		return nil, err
	}

	// Extract markdown content after the closing --- delimiter
	content := extractContent(raw)

	var fm skillFrontmatter
	if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
		return nil, err
	}

	return &model.Skill{
		Name:        fm.Name,
		Description: strings.TrimSpace(fm.Description),
		Author:      fm.Metadata.Author,
		Version:     fm.Metadata.Version,
		Path:        path,
		Content:     content,
	}, nil
}

// extractFrontmatter extracts YAML content between --- delimiters.
func extractFrontmatter(content string) (string, error) {
	const delimiter = "---"

	// Find first delimiter
	start := strings.Index(content, delimiter)
	if start == -1 {
		return "", os.ErrNotExist
	}
	start += len(delimiter)

	// Find second delimiter
	end := strings.Index(content[start:], delimiter)
	if end == -1 {
		return "", os.ErrNotExist
	}

	return content[start : start+end], nil
}

// extractContent returns the markdown content after the closing --- frontmatter delimiter.
func extractContent(raw string) string {
	const delimiter = "---"

	// Find first delimiter
	start := strings.Index(raw, delimiter)
	if start == -1 {
		return ""
	}
	start += len(delimiter)

	// Find second (closing) delimiter
	end := strings.Index(raw[start:], delimiter)
	if end == -1 {
		return ""
	}

	// Content starts after the closing delimiter
	contentStart := start + end + len(delimiter)
	if contentStart >= len(raw) {
		return ""
	}

	return strings.TrimSpace(raw[contentStart:])
}
