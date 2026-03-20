package tui

import (
	"fmt"
	"strings"

	"github.com/jarlex/opencode-configurator/internal/model"
)

// RenderAgentDetail renders a styled detail view for an Agent.
func RenderAgentDetail(a *model.Agent) string {
	s := renderAgentHeader(a)
	s += renderAgentTools(a)
	s += renderAgentPermissions(a)
	s += renderAgentPrompt(a)
	return s
}

func renderAgentHeader(a *model.Agent) string {
	s := DetailTitle.Render(a.Name) + "\n\n"
	s += DetailLabel.Render("Mode: ") + DetailValue.Render(a.Mode) + "\n"
	if a.Model != "" {
		s += DetailLabel.Render("Model: ") + DetailValue.Render(a.Model) + "\n"
	}
	if a.Description != "" {
		s += DetailLabel.Render("Description: ") + DetailValue.Render(a.Description) + "\n"
	}
	if a.MaxSteps > 0 {
		s += DetailLabel.Render("Max Steps: ") + DetailValue.Render(fmt.Sprintf("%d", a.MaxSteps)) + "\n"
	}
	if a.Temperature > 0 {
		s += DetailLabel.Render("Temperature: ") + DetailValue.Render(fmt.Sprintf("%.1f", a.Temperature)) + "\n"
	}
	if a.TopP > 0 {
		s += DetailLabel.Render("Top P: ") + DetailValue.Render(fmt.Sprintf("%.1f", a.TopP)) + "\n"
	}
	if a.Hidden {
		s += DetailLabel.Render("Hidden: ") + DetailValue.Render("yes") + "\n"
	}
	return s
}

func renderAgentTools(a *model.Agent) string {
	if len(a.Tools) == 0 {
		return ""
	}
	s := "\n" + DetailLabel.Render("Tools:") + "\n"
	for tool, enabled := range a.Tools {
		status := "  \u2713 "
		if !enabled {
			status = "  \u2717 "
		}
		s += status + tool + "\n"
	}
	return s
}

func renderAgentPermissions(a *model.Agent) string {
	s := "\n" + DetailLabel.Render("Permissions:") + "\n"
	if a.Permission.Edit != "" {
		s += "  Edit: " + a.Permission.Edit + "\n"
	}
	if a.Permission.Bash != "" {
		s += "  Bash: " + a.Permission.Bash + "\n"
	}
	if a.Permission.Webfetch != "" {
		s += "  Webfetch: " + a.Permission.Webfetch + "\n"
	}
	if a.Permission.DoomLoop != "" {
		s += "  DoomLoop: " + a.Permission.DoomLoop + "\n"
	}
	if a.Permission.ExternalDirectory != "" {
		s += "  External Directory: " + a.Permission.ExternalDirectory + "\n"
	}
	if len(a.Permission.Task) > 0 {
		s += "  Task:\n"
		for pattern, perm := range a.Permission.Task {
			s += "    " + pattern + ": " + perm + "\n"
		}
	}
	return s
}

func renderAgentPrompt(a *model.Agent) string {
	if a.Prompt == "" {
		return ""
	}
	s := "\n" + DetailLabel.Render("Prompt:") + "\n"
	s += DetailMuted.Render(a.Prompt) + "\n"
	return s
}

// RenderSkillDetail renders a styled detail view for a Skill.
func RenderSkillDetail(sk *model.Skill) string {
	s := DetailTitle.Render(sk.Name) + "\n\n"
	if sk.Description != "" {
		s += DetailLabel.Render("Description: ") + DetailValue.Render(sk.Description) + "\n"
	}
	if sk.Author != "" {
		s += DetailLabel.Render("Author: ") + DetailValue.Render(sk.Author) + "\n"
	}
	if sk.Version != "" {
		s += DetailLabel.Render("Version: ") + DetailValue.Render(sk.Version) + "\n"
	}
	s += DetailLabel.Render("Path: ") + DetailMuted.Render(sk.Path) + "\n"

	// Skill content / prompt — displayed in full (scrollable via viewport)
	if sk.Content != "" {
		s += "\n" + DetailLabel.Render("Skill Prompt:") + "\n"
		s += DetailMuted.Render(sk.Content) + "\n"
	}

	return s
}

// RenderMCPDetail renders a styled detail view for an MCPServer.
func RenderMCPDetail(m *model.MCPServer) string {
	s := DetailTitle.Render(m.Name) + "\n\n"
	s += DetailLabel.Render("Type: ") + DetailValue.Render(m.Type) + "\n"

	enabled := "yes"
	if !m.Enabled {
		enabled = "no"
	}
	s += DetailLabel.Render("Enabled: ") + DetailValue.Render(enabled) + "\n"

	if m.Status != "" {
		s += DetailLabel.Render("Status: ") + DetailValue.Render(m.Status) + "\n"
	}
	if m.URL != "" {
		s += DetailLabel.Render("URL: ") + DetailValue.Render(m.URL) + "\n"
	}
	if len(m.Command) > 0 {
		s += DetailLabel.Render("Command: ") + DetailMuted.Render(strings.Join(m.Command, " ")) + "\n"
	}
	if m.Error != "" {
		s += "\n" + StatusOffline.Render("Error: "+m.Error) + "\n"
	}
	if len(m.Environment) > 0 {
		s += "\n" + DetailLabel.Render("Environment:") + "\n"
		for k, v := range m.Environment {
			s += "  " + k + "=" + v + "\n"
		}
	}
	if m.Timeout > 0 {
		s += DetailLabel.Render("Timeout: ") + DetailValue.Render(fmt.Sprintf("%ds", m.Timeout)) + "\n"
	}

	return s
}

// RenderProviderDetail renders a styled detail view for a Provider.
func RenderProviderDetail(p *model.Provider) string {
	s := DetailTitle.Render(p.Name) + "\n\n"
	if p.NPM != "" {
		s += DetailLabel.Render("Package: ") + DetailValue.Render(p.NPM) + "\n"
	}
	if p.BaseURL != "" {
		s += DetailLabel.Render("Base URL: ") + DetailValue.Render(p.BaseURL) + "\n"
	}
	if len(p.Models) > 0 {
		s += "\n" + DetailLabel.Render("Models:") + "\n"
		for _, m := range p.Models {
			launch := ""
			if m.Launch {
				launch = " \u2605"
			}
			s += "  \u2022 " + m.Name + launch + "\n"
		}
	}
	return s
}
