package onboarding

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/Greite/unraid-tui/internal/config"
	"github.com/Greite/unraid-tui/internal/tui/common"
)

type step int

const (
	stepWelcome step = iota
	stepServerName
	stepServerURL
	stepTestConnection
	stepAPIKeyInfo
	stepAPIKey
	stepTestAPIKey
	stepSaving
	stepDone
)

// connectionTestMsg is sent after testing server connectivity.
type connectionTestMsg struct {
	ok  bool
	err error
}

// apiKeyTestMsg is sent after testing the API key.
type apiKeyTestMsg struct {
	ok  bool
	err error
}

// saveResultMsg is sent after saving config.
type saveResultMsg struct {
	err error
}

type Model struct {
	step       step
	serverName textinput.Model
	serverURL  textinput.Model
	apiKey     textinput.Model
	spinner    spinner.Model
	err        error
	width      int
	height     int
	config     *config.Config
	quitting   bool
	completed  bool
}

func New() Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "NAS"
	nameInput.Prompt = "  > "
	nameInput.CharLimit = 64

	serverInput := textinput.New()
	serverInput.Placeholder = "http://192.168.1.100"
	serverInput.Prompt = "  > "
	serverInput.CharLimit = 256

	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "unraid-api-key-xxxx..."
	apiKeyInput.Prompt = "  > "
	apiKeyInput.CharLimit = 512
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.EchoCharacter = '*'

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(common.ColorPrimary)

	return Model{
		step:       stepWelcome,
		serverName: nameInput,
		serverURL:  serverInput,
		apiKey:     apiKeyInput,
		spinner:    s,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.step == stepWelcome || m.step == stepServerName {
				m.quitting = true
				return m, tea.Quit
			}
			if m.step > stepWelcome && m.step < stepDone {
				m.err = nil
				m.step--
				if m.step == stepTestConnection {
					m.step = stepServerURL
				}
				if m.step == stepTestAPIKey {
					m.step = stepAPIKey
				}
				if m.step == stepSaving {
					m.step = stepAPIKey
				}
				m = m.focusCurrentInput()
				return m, nil
			}
		}

	case connectionTestMsg:
		if msg.ok {
			m.err = nil
			m.step = stepAPIKeyInfo
			return m, nil
		}
		m.err = msg.err
		m.step = stepServerURL
		m = m.focusCurrentInput()
		return m, nil

	case apiKeyTestMsg:
		if msg.ok {
			m.err = nil
			m.step = stepSaving
			return m, m.saveConfig()
		}
		m.err = msg.err
		m.step = stepAPIKey
		m = m.focusCurrentInput()
		return m, nil

	case saveResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.step = stepAPIKey
			return m, nil
		}
		m.step = stepDone
		m.completed = true
		return m, nil
	}

	// Step-specific key handling
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch m.step {
		case stepWelcome:
			if keyMsg.String() == "enter" {
				m.step = stepServerName
				m = m.focusCurrentInput()
				return m, m.serverName.Focus()
			}
		case stepServerName:
			if keyMsg.String() == "enter" {
				name := strings.TrimSpace(m.serverName.Value())
				if name == "" {
					m.err = fmt.Errorf("le nom du serveur ne peut pas etre vide")
					return m, nil
				}
				m.err = nil
				m.step = stepServerURL
				m = m.focusCurrentInput()
				return m, m.serverURL.Focus()
			}
		case stepServerURL:
			if keyMsg.String() == "enter" {
				url := strings.TrimSpace(m.serverURL.Value())
				if url == "" {
					m.err = fmt.Errorf("l'URL du serveur ne peut pas etre vide")
					return m, nil
				}
				url = normalizeURL(url)
				m.serverURL.SetValue(url)
				m.err = nil
				m.step = stepTestConnection
				return m, m.testConnection(url)
			}
		case stepAPIKeyInfo:
			if keyMsg.String() == "enter" {
				m.step = stepAPIKey
				m = m.focusCurrentInput()
				return m, m.apiKey.Focus()
			}
		case stepAPIKey:
			if keyMsg.String() == "enter" {
				key := strings.TrimSpace(m.apiKey.Value())
				if key == "" {
					m.err = fmt.Errorf("la cle API ne peut pas etre vide")
					return m, nil
				}
				m.err = nil
				m.step = stepTestAPIKey
				return m, m.testAPIKey(m.serverURL.Value(), key)
			}
		case stepDone:
			if keyMsg.String() == "enter" {
				return m, tea.Quit
			}
		}
	}

	// Update active input
	var cmd tea.Cmd
	switch m.step {
	case stepServerName:
		m.serverName, cmd = m.serverName.Update(msg)
		return m, cmd
	case stepServerURL:
		m.serverURL, cmd = m.serverURL.Update(msg)
		return m, cmd
	case stepAPIKey:
		m.apiKey, cmd = m.apiKey.Update(msg)
		return m, cmd
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	var s strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(common.ColorPrimary).
		MarginBottom(1).
		Render("  UNRAID CLI — Configuration")

	s.WriteString("\n" + title + "\n\n")

	// Progress
	s.WriteString(m.renderProgress())
	s.WriteString("\n\n")

	// Step content
	switch m.step {
	case stepWelcome:
		s.WriteString(m.renderWelcome())
	case stepServerName:
		s.WriteString(m.renderServerName())
	case stepServerURL:
		s.WriteString(m.renderServerURL())
	case stepTestConnection:
		s.WriteString(m.renderTestConnection())
	case stepAPIKeyInfo:
		s.WriteString(m.renderAPIKeyInfo())
	case stepAPIKey:
		s.WriteString(m.renderAPIKey())
	case stepTestAPIKey:
		s.WriteString(m.renderTestAPIKey())
	case stepSaving:
		s.WriteString(m.renderSaving())
	case stepDone:
		s.WriteString(m.renderDone())
	}

	// Error
	if m.err != nil {
		errStyle := lipgloss.NewStyle().
			Foreground(common.ColorDanger).
			Bold(true).
			MarginTop(1)
		s.WriteString("\n" + errStyle.Render("  ✗ "+m.err.Error()) + "\n")
	}

	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}

// Completed returns true if onboarding finished successfully.
func (m Model) Completed() bool {
	return m.completed
}

// Quitting returns true if user pressed Ctrl+C.
func (m Model) Quitting() bool {
	return m.quitting
}

// --- Rendering ---

func (m Model) renderProgress() string {
	steps := []struct {
		label string
		s     step
	}{
		{"Nom", stepServerName},
		{"URL", stepServerURL},
		{"Connexion", stepTestConnection},
		{"Cle API", stepAPIKey},
		{"Termine", stepDone},
	}

	var parts []string
	for _, st := range steps {
		style := lipgloss.NewStyle().Foreground(common.ColorMuted)
		marker := "○"
		if m.step > st.s || (m.step == stepDone && st.s == stepDone) {
			style = lipgloss.NewStyle().Foreground(common.ColorSuccess)
			marker = "●"
		} else if m.step >= st.s {
			style = lipgloss.NewStyle().Foreground(common.ColorPrimary).Bold(true)
			marker = "◉"
		}
		parts = append(parts, style.Render(marker+" "+st.label))
	}
	return "  " + strings.Join(parts, "  —  ")
}

func (m Model) renderWelcome() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(common.ColorPrimary).
		Padding(1, 3).
		MarginLeft(2).
		Width(60)

	content := "Bienvenue !\n\n"
	content += "Cet assistant va vous aider a configurer la connexion\n"
	content += "a votre serveur Unraid en quelques etapes :\n\n"
	content += "  1. Saisir l'adresse de votre serveur\n"
	content += "  2. Tester la connexion\n"
	content += "  3. Configurer votre cle API\n"
	content += "  4. Sauvegarder la configuration\n\n"
	content += "Le fichier sera sauvegarde dans ~/.unraid-tui/config.yaml"

	return box.Render(content) + "\n\n" + actionHint("enter", "commencer") + "  " + escHint()
}

func (m Model) renderServerName() string {
	var s strings.Builder
	s.WriteString(stepTitle("Etape 1/4", "Nom du serveur"))
	s.WriteString("\n")
	s.WriteString("  Donnez un nom a votre serveur (ex: NAS, Backup, Media).\n")
	s.WriteString(common.StyleSubtle.Render("  Ce nom permet d'identifier le serveur dans la liste.") + "\n\n")
	s.WriteString(m.serverName.View() + "\n\n")
	s.WriteString(actionHint("enter", "continuer") + "  " + escHint())
	return s.String()
}

func (m Model) renderServerURL() string {
	var s strings.Builder
	s.WriteString(stepTitle("Etape 2/4", "Adresse du serveur Unraid"))
	s.WriteString("\n")
	s.WriteString("  Entrez l'URL de votre serveur Unraid (avec le port).\n")
	s.WriteString(common.StyleSubtle.Render("  Par defaut, l'API Unraid ecoute sur le port 3001.") + "\n\n")
	s.WriteString(m.serverURL.View() + "\n\n")
	s.WriteString(actionHint("enter", "tester la connexion") + "  " + escHint())
	return s.String()
}

func (m Model) renderTestConnection() string {
	return "  " + m.spinner.View() + " Test de la connexion a " + common.StyleTitle.Render(m.serverURL.Value()) + "..."
}

func (m Model) renderAPIKeyInfo() string {
	var s strings.Builder
	s.WriteString(stepTitle("Etape 3/4", "Creer une cle API"))
	s.WriteString("\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(common.ColorBorder).
		Padding(1, 2).
		MarginLeft(2).
		Width(64)

	instructions := "Comment obtenir une cle API :\n\n"
	instructions += "  1. Ouvrez l'interface web de votre serveur Unraid\n"
	instructions += "  2. Allez dans Settings > Management Access\n"
	instructions += "  3. Activez Developer Options\n"
	instructions += "  4. Ouvrez Apollo GraphQL Studio\n"
	instructions += "  5. Executez cette mutation :\n\n"
	instructions += "     mutation {\n"
	instructions += "       apiKey {\n"
	instructions += "         create(input: {\n"
	instructions += `           name: "unraid-tui"` + "\n"
	instructions += "           roles: [ADMIN]\n"
	instructions += "         }) { key }\n"
	instructions += "       }\n"
	instructions += "     }\n\n"
	instructions += "  6. Copiez la cle retournee\n"

	s.WriteString(box.Render(instructions) + "\n\n")
	s.WriteString(actionHint("enter", "saisir la cle") + "  " + escHint())
	return s.String()
}

func (m Model) renderAPIKey() string {
	var s strings.Builder
	s.WriteString(stepTitle("Etape 4/4", "Saisir la cle API"))
	s.WriteString("\n")
	s.WriteString("  Collez votre cle API Unraid ci-dessous.\n")
	s.WriteString(common.StyleSubtle.Render("  La cle est masquee pour des raisons de securite.") + "\n\n")
	s.WriteString(m.apiKey.View() + "\n\n")
	s.WriteString(actionHint("enter", "valider") + "  " + escHint())
	return s.String()
}

func (m Model) renderTestAPIKey() string {
	return "  " + m.spinner.View() + " Verification de la cle API..."
}

func (m Model) renderSaving() string {
	return "  " + m.spinner.View() + " Sauvegarde de la configuration..."
}

func (m Model) renderDone() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(common.ColorSuccess).
		Padding(1, 3).
		MarginLeft(2).
		Width(60)

	content := "Configuration terminee !\n\n"
	content += "Votre configuration a ete sauvegardee dans :\n"
	content += "  " + config.FilePath() + "\n\n"
	content += "Serveur : " + m.serverURL.Value() + "\n"
	content += "Cle API : ********** (sauvegardee)\n\n"
	content += "Le dashboard va maintenant se lancer."

	return box.Render(content) + "\n\n" + actionHint("enter", "lancer le dashboard")
}

// --- Commands ---

func (m Model) testConnection(url string) tea.Cmd {
	return func() tea.Msg {
		endpoint := url + "/graphql"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(`{"query":"{ __typename }"}`))
		if err != nil {
			return connectionTestMsg{ok: false, err: fmt.Errorf("URL invalide : %w", err)}
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return connectionTestMsg{ok: false, err: fmt.Errorf("impossible de joindre le serveur : %s", cleanHTTPError(err))}
		}
		resp.Body.Close()

		// Any HTTP response means the server is reachable (even 401)
		return connectionTestMsg{ok: true}
	}
}

func (m Model) testAPIKey(url, key string) tea.Cmd {
	return func() tea.Msg {
		endpoint := url + "/graphql"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := `{"query":"{ info { os { hostname } } }"}`
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(body))
		if err != nil {
			return apiKeyTestMsg{ok: false, err: fmt.Errorf("erreur requete : %w", err)}
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+key)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return apiKeyTestMsg{ok: false, err: fmt.Errorf("erreur connexion : %s", cleanHTTPError(err))}
		}
		resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			return apiKeyTestMsg{ok: true}
		case http.StatusUnauthorized, http.StatusForbidden:
			return apiKeyTestMsg{ok: false, err: fmt.Errorf("cle API invalide ou permissions insuffisantes (HTTP %d)", resp.StatusCode)}
		default:
			return apiKeyTestMsg{ok: false, err: fmt.Errorf("reponse inattendue du serveur (HTTP %d)", resp.StatusCode)}
		}
	}
}

func (m Model) saveConfig() tea.Cmd {
	return func() tea.Msg {
		name := strings.TrimSpace(m.serverName.Value())
		if name == "" {
			name = "default"
		}
		cfg := &config.Config{
			ServerURL: m.serverURL.Value(),
			APIKey:    m.apiKey.Value(),
		}
		err := config.SaveServer(name, cfg)
		return saveResultMsg{err: err}
	}
}

// --- Helpers ---

func (m Model) focusCurrentInput() Model {
	m.serverName.Blur()
	m.serverURL.Blur()
	m.apiKey.Blur()
	return m
}

func normalizeURL(url string) string {
	url = strings.TrimRight(url, "/")
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return url
}

func cleanHTTPError(err error) string {
	msg := err.Error()
	// Trim verbose Go HTTP error wrapping
	if idx := strings.LastIndex(msg, ": "); idx != -1 {
		short := msg[idx+2:]
		if len(short) > 5 {
			return short
		}
	}
	return msg
}

func stepTitle(number, title string) string {
	num := lipgloss.NewStyle().Foreground(common.ColorMuted).Render(number)
	ttl := lipgloss.NewStyle().Bold(true).Foreground(common.ColorText).Render(title)
	return "  " + num + " — " + ttl + "\n"
}

func actionHint(key, desc string) string {
	k := lipgloss.NewStyle().Bold(true).Foreground(common.ColorText).Render(key)
	d := lipgloss.NewStyle().Foreground(common.ColorMuted).Render(desc)
	return "  " + k + " " + d
}

func escHint() string {
	return lipgloss.NewStyle().Foreground(common.ColorMuted).Render("esc retour")
}
