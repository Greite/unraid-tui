package onboarding

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNew_InitialState(t *testing.T) {
	m := New()
	if m.step != stepWelcome {
		t.Errorf("expected stepWelcome, got %d", m.step)
	}
	if m.completed {
		t.Error("expected completed to be false")
	}
	if m.quitting {
		t.Error("expected quitting to be false")
	}
}

func TestWelcome_EnterGoesToServerName(t *testing.T) {
	m := New()
	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model := updated.(Model)
	if model.step != stepServerName {
		t.Errorf("expected stepServerName, got %d", model.step)
	}
}

func TestServerURL_EmptyShowsError(t *testing.T) {
	m := New()
	m.step = stepServerURL

	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model := updated.(Model)
	if model.err == nil {
		t.Fatal("expected error for empty URL")
	}
	if model.step != stepServerURL {
		t.Errorf("expected to stay on stepServerURL, got %d", model.step)
	}
}

func TestServerURL_ValidGoesToTestConnection(t *testing.T) {
	m := New()
	m.step = stepServerURL
	m.serverURL.SetValue("http://192.168.1.100:3001")

	updated, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model := updated.(Model)
	if model.step != stepTestConnection {
		t.Errorf("expected stepTestConnection, got %d", model.step)
	}
	if cmd == nil {
		t.Error("expected a command for connection test")
	}
}

func TestConnectionTestMsg_SuccessGoesToAPIKeyInfo(t *testing.T) {
	m := New()
	m.step = stepTestConnection

	updated, _ := m.Update(connectionTestMsg{ok: true})
	model := updated.(Model)
	if model.step != stepAPIKeyInfo {
		t.Errorf("expected stepAPIKeyInfo, got %d", model.step)
	}
	if model.err != nil {
		t.Errorf("expected no error, got %v", model.err)
	}
}

func TestConnectionTestMsg_FailureGoesBackToServerURL(t *testing.T) {
	m := New()
	m.step = stepTestConnection

	updated, _ := m.Update(connectionTestMsg{ok: false, err: errTest("timeout")})
	model := updated.(Model)
	if model.step != stepServerURL {
		t.Errorf("expected stepServerURL, got %d", model.step)
	}
	if model.err == nil {
		t.Error("expected error to be set")
	}
}

func TestAPIKeyInfo_EnterGoesToAPIKey(t *testing.T) {
	m := New()
	m.step = stepAPIKeyInfo

	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model := updated.(Model)
	if model.step != stepAPIKey {
		t.Errorf("expected stepAPIKey, got %d", model.step)
	}
}

func TestAPIKey_EmptyShowsError(t *testing.T) {
	m := New()
	m.step = stepAPIKey

	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model := updated.(Model)
	if model.err == nil {
		t.Fatal("expected error for empty key")
	}
	if model.step != stepAPIKey {
		t.Errorf("expected to stay on stepAPIKey, got %d", model.step)
	}
}

func TestAPIKey_ValidGoesToTestAPIKey(t *testing.T) {
	m := New()
	m.step = stepAPIKey
	m.serverURL.SetValue("http://192.168.1.100:3001")
	m.apiKey.SetValue("my-secret-key")

	updated, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model := updated.(Model)
	if model.step != stepTestAPIKey {
		t.Errorf("expected stepTestAPIKey, got %d", model.step)
	}
	if cmd == nil {
		t.Error("expected a command for API key test")
	}
}

func TestAPIKeyTestMsg_SuccessGoesToSaving(t *testing.T) {
	m := New()
	m.step = stepTestAPIKey

	updated, cmd := m.Update(apiKeyTestMsg{ok: true})
	model := updated.(Model)
	if model.step != stepSaving {
		t.Errorf("expected stepSaving, got %d", model.step)
	}
	if cmd == nil {
		t.Error("expected a save command")
	}
}

func TestAPIKeyTestMsg_FailureGoesBackToAPIKey(t *testing.T) {
	m := New()
	m.step = stepTestAPIKey

	updated, _ := m.Update(apiKeyTestMsg{ok: false, err: errTest("401")})
	model := updated.(Model)
	if model.step != stepAPIKey {
		t.Errorf("expected stepAPIKey, got %d", model.step)
	}
	if model.err == nil {
		t.Error("expected error to be set")
	}
}

func TestSaveResultMsg_SuccessGoesToDone(t *testing.T) {
	m := New()
	m.step = stepSaving

	updated, _ := m.Update(saveResultMsg{err: nil})
	model := updated.(Model)
	if model.step != stepDone {
		t.Errorf("expected stepDone, got %d", model.step)
	}
	if !model.completed {
		t.Error("expected completed to be true")
	}
}

func TestSaveResultMsg_FailureGoesBackToAPIKey(t *testing.T) {
	m := New()
	m.step = stepSaving

	updated, _ := m.Update(saveResultMsg{err: errTest("permission denied")})
	model := updated.(Model)
	if model.step != stepAPIKey {
		t.Errorf("expected stepAPIKey, got %d", model.step)
	}
	if model.err == nil {
		t.Error("expected error to be set")
	}
}

func TestCtrlC_Quits(t *testing.T) {
	m := New()
	updated, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}))
	model := updated.(Model)
	if !model.quitting {
		t.Error("expected quitting to be true")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestEsc_QuitsOnWelcome(t *testing.T) {
	m := New()
	updated, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))
	model := updated.(Model)
	if !model.quitting {
		t.Error("expected quitting on esc at welcome")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestEsc_QuitsOnServerName(t *testing.T) {
	m := New()
	m.step = stepServerName
	updated, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))
	model := updated.(Model)
	if !model.quitting {
		t.Error("expected quitting on esc at serverName")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestEsc_GoesBackFromAPIKey(t *testing.T) {
	m := New()
	m.step = stepAPIKey

	updated, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))
	model := updated.(Model)
	if model.step != stepAPIKeyInfo {
		t.Errorf("expected stepAPIKeyInfo after esc, got %d", model.step)
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"192.168.1.100:3001", "http://192.168.1.100:3001"},
		{"http://tower:3001", "http://tower:3001"},
		{"https://secure.local:3001/", "https://secure.local:3001"},
		{"http://192.168.1.100:3001/", "http://192.168.1.100:3001"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeURL(tt.input)
			if got != tt.want {
				t.Errorf("normalizeURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestView_WelcomeContainsExpectedContent(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 30
	v := m.View()
	if !strings.Contains(v.Content, "Bienvenue") {
		t.Error("expected 'Bienvenue' in welcome view")
	}
	if !strings.Contains(v.Content, "commencer") {
		t.Error("expected 'commencer' hint in welcome view")
	}
}

func TestView_DoneContainsSuccess(t *testing.T) {
	m := New()
	m.step = stepDone
	m.width = 80
	m.height = 30
	m.serverURL.SetValue("http://tower:3001")
	v := m.View()
	if !strings.Contains(v.Content, "terminee") {
		t.Error("expected 'terminee' in done view")
	}
	if !strings.Contains(v.Content, "tower") {
		t.Error("expected server URL in done view")
	}
}

func TestView_ServerURLStep(t *testing.T) {
	m := New()
	m.step = stepServerURL
	m.width = 80
	m.height = 30
	v := m.View()
	if !strings.Contains(v.Content, "Etape 2/4") {
		t.Error("expected 'Etape 2/4' in server URL view")
	}
	if !strings.Contains(v.Content, "3001") {
		t.Error("expected port hint in server URL view")
	}
}

func TestView_APIKeyInfoStep(t *testing.T) {
	m := New()
	m.step = stepAPIKeyInfo
	m.width = 80
	m.height = 30
	v := m.View()
	if !strings.Contains(v.Content, "mutation") {
		t.Error("expected GraphQL mutation in API key info view")
	}
	if !strings.Contains(v.Content, "unraid-tui") {
		t.Error("expected 'unraid-tui' key name in mutation")
	}
}

func TestView_ErrorDisplayed(t *testing.T) {
	m := New()
	m.step = stepServerURL
	m.err = errTest("something went wrong")
	m.width = 80
	m.height = 30
	v := m.View()
	if !strings.Contains(v.Content, "something went wrong") {
		t.Error("expected error message in view")
	}
}

func TestDone_EnterQuits(t *testing.T) {
	m := New()
	m.step = stepDone
	_, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	if cmd == nil {
		t.Error("expected quit command on enter at done step")
	}
}

// helper
type testError string

func (e testError) Error() string { return string(e) }

func errTest(msg string) error { return testError(msg) }
