package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/path/opencode.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if got := err.Error(); got != "config file not found: /nonexistent/path/opencode.json" {
		t.Fatalf("unexpected error message: %s", got)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "opencode.json")
	if err := os.WriteFile(path, []byte("{invalid json}"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParse_RealConfig(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}
	configPath := filepath.Join(homeDir, ".config", "opencode", "opencode.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("opencode.json not found, skipping integration test")
	}

	state, err := Parse(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Agents) == 0 {
		t.Error("expected at least one agent")
	}
	if len(state.MCPs) == 0 {
		t.Error("expected at least one MCP server")
	}
	if len(state.Providers) == 0 {
		t.Error("expected at least one provider")
	}

	// Verify agent fields are populated
	for _, a := range state.Agents {
		if a.Name == "" {
			t.Error("agent name should not be empty")
		}
		if a.Mode == "" {
			t.Errorf("agent %s: mode should not be empty", a.Name)
		}
	}
}

func TestScanSkills_MissingDir(t *testing.T) {
	skills, err := ScanSkills("/nonexistent/skills/dir")
	if err != nil {
		t.Fatalf("expected no error for missing dir, got: %v", err)
	}
	if len(skills) != 0 {
		t.Fatalf("expected empty skill list, got %d", len(skills))
	}
}

func TestScanSkills_RealDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}
	skillsDir := filepath.Join(homeDir, ".config", "opencode", "skills")
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		t.Skip("skills directory not found, skipping integration test")
	}

	skills, err := ScanSkills(skillsDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) == 0 {
		t.Error("expected at least one skill")
	}

	// Verify skill fields are populated
	for _, s := range skills {
		if s.Name == "" {
			t.Error("skill name should not be empty")
		}
		if s.Path == "" {
			t.Error("skill path should not be empty")
		}
		if s.Description == "" {
			t.Errorf("skill %s: description should not be empty", s.Name)
		}
	}
}

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid frontmatter",
			input: "---\nname: test\n---\n# Content",
			want:  "\nname: test\n",
		},
		{
			name:    "no frontmatter",
			input:   "# Just markdown",
			wantErr: true,
		},
		{
			name:    "single delimiter",
			input:   "---\nname: test\n# No closing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractFrontmatter(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
