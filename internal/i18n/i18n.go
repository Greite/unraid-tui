package i18n

import (
	"os"
	"strings"
)

var currentLang = "en"

var translations = map[string]map[string]string{
	// === General ===
	"loading":       {"en": "Loading...", "fr": "Chargement..."},
	"waiting":       {"en": "Waiting...", "fr": "En attente..."},
	"error":         {"en": "Error", "fr": "Erreur"},
	"none":          {"en": "None", "fr": "Aucun"},
	"quit":          {"en": "quit", "fr": "quitter"},
	"pages":         {"en": "pages", "fr": "pages"},
	"next":          {"en": "next", "fr": "suivant"},
	"servers":       {"en": "servers", "fr": "serveurs"},
	"navigate":      {"en": "navigate", "fr": "naviguer"},
	"refresh":       {"en": "refresh", "fr": "rafraichir"},
	"sort":          {"en": "sort", "fr": "trier"},
	"back":          {"en": "back", "fr": "retour"},
	"scroll":        {"en": "scroll", "fr": "scroll"},
	"follow":        {"en": "follow", "fr": "follow"},
	"start_end":     {"en": "start/end", "fr": "debut/fin"},
	"execute":       {"en": "execute", "fr": "executer"},
	"select":        {"en": "select", "fr": "selectionner"},
	"connect":       {"en": "connect", "fr": "connecter"},
	"default":       {"en": "default", "fr": "defaut"},
	"delete":        {"en": "delete", "fr": "supprimer"},
	"close":         {"en": "close", "fr": "fermer"},
	"archive":       {"en": "archive", "fr": "archiver"},
	"archive_all":   {"en": "archive all", "fr": "archiver tout"},
	"begin":         {"en": "begin", "fr": "commencer"},
	"continue":      {"en": "continue", "fr": "continuer"},
	"validate":      {"en": "validate", "fr": "valider"},
	"enter_key":     {"en": "enter key", "fr": "saisir la cle"},
	"launch_dash":   {"en": "launch dashboard", "fr": "lancer le dashboard"},
	"line":          {"en": "line", "fr": "ligne"},
	"lang":          {"en": "language", "fr": "langue"},

	// === Pages ===
	"page_dashboard":     {"en": "Dashboard", "fr": "Dashboard"},
	"page_docker":        {"en": "Docker", "fr": "Docker"},
	"page_vms":           {"en": "VMs", "fr": "VMs"},
	"page_notifications": {"en": "Notifs", "fr": "Notifs"},
	"page_shares":        {"en": "Shares", "fr": "Shares"},

	// === Dashboard ===
	"system":       {"en": "System", "fr": "Systeme"},
	"hostname":     {"en": "Hostname", "fr": "Hostname"},
	"uptime":       {"en": "Uptime", "fr": "Uptime"},
	"cpu":          {"en": "CPU", "fr": "CPU"},
	"cpu_cores":    {"en": "CPU Cores", "fr": "CPU Cores"},
	"memory":       {"en": "Memory", "fr": "Memoire"},
	"network":      {"en": "Network", "fr": "Reseau"},
	"disks":        {"en": "Disks", "fr": "Disques"},
	"hardware":     {"en": "Hardware", "fr": "Materiel"},
	"parity":       {"en": "Parity", "fr": "Parite"},
	"array":        {"en": "Array", "fr": "Array"},
	"total":        {"en": "total", "fr": "total"},
	"running":      {"en": "running", "fr": "en cours"},
	"exited":       {"en": "exited", "fr": "arrete"},
	"paused":       {"en": "paused", "fr": "pause"},
	"devices":      {"en": "devices", "fr": "peripheriques"},

	// === Docker ===
	"containers":      {"en": "Containers", "fr": "Containers"},
	"loading_docker":  {"en": "Loading containers...", "fr": "Chargement des containers..."},
	"docker_disabled": {"en": "Docker is not enabled on this server.", "fr": "Docker n'est pas active sur ce serveur."},
	"docker_enable":   {"en": "Enable it in Settings > Docker.", "fr": "Activez-le dans Settings > Docker."},
	"logs":            {"en": "logs", "fr": "logs"},
	"console":         {"en": "console", "fr": "console"},
	"webui":           {"en": "WebUI", "fr": "WebUI"},
	"start":           {"en": "start", "fr": "demarrer"},
	"stop":            {"en": "stop", "fr": "arreter"},
	"pause":           {"en": "pause", "fr": "pause"},
	"unpause":         {"en": "unpause", "fr": "reprendre"},
	"update":          {"en": "update", "fr": "mettre a jour"},
	"update_all":      {"en": "update all", "fr": "tout mettre a jour"},
	"no_webui":        {"en": "No WebUI for %s", "fr": "Pas de WebUI pour %s"},
	"webui_opened":    {"en": "WebUI opened for %s", "fr": "WebUI ouvert pour %s"},
	"not_running":     {"en": "%s is not running", "fr": "%s n'est pas running"},
	"console_done":    {"en": "Console finished", "fr": "Console terminee"},
	"console_error":   {"en": "Console finished with error", "fr": "Console terminee avec erreur"},
	"connected_to":    {"en": "Connected to %s via SSH", "fr": "Connecte a %s via SSH"},
	"logs_error":      {"en": "Logs error: %s", "fr": "Erreur logs: %s"},
	"action_ok":       {"en": "%s %s OK", "fr": "%s %s OK"},
	"action_error":    {"en": "Error %s %s: %s", "fr": "Erreur %s %s: %s"},
	"follow_on":       {"en": "FOLLOW", "fr": "SUIVI"},
	"follow_off":      {"en": "PAUSE", "fr": "PAUSE"},

	// === VMs ===
	"loading_vms":  {"en": "Loading VMs...", "fr": "Chargement des VMs..."},
	"no_vms":       {"en": "No VMs configured", "fr": "Aucune VM configuree"},
	"vms_disabled": {"en": "VMs are not enabled on this server.", "fr": "Les VMs ne sont pas activees sur ce serveur."},
	"vms_enable":   {"en": "Enable them in Settings > VM Manager.", "fr": "Activez-les dans Settings > VM Manager."},
	"reboot":       {"en": "reboot", "fr": "redemarrer"},
	"force_stop":   {"en": "force stop", "fr": "forcer l'arret"},
	"resume":       {"en": "resume", "fr": "reprendre"},

	// === Notifications ===
	"loading_notifs": {"en": "Loading notifications...", "fr": "Chargement des notifications..."},
	"no_notifs":      {"en": "No unread notifications", "fr": "Aucune notification non lue"},
	"notifications":  {"en": "Notifications", "fr": "Notifications"},

	// === Shares ===
	"loading_shares": {"en": "Loading shares...", "fr": "Chargement des shares..."},
	"no_shares":      {"en": "No shares configured", "fr": "Aucun share configure"},
	"shares":         {"en": "Shares", "fr": "Shares"},

	// === Onboarding ===
	"onboarding_title":   {"en": "UNRAID TUI — Configuration", "fr": "UNRAID TUI — Configuration"},
	"welcome":            {"en": "Welcome!", "fr": "Bienvenue !"},
	"welcome_desc":       {"en": "This wizard will help you configure the connection\nto your Unraid server in a few steps:", "fr": "Cet assistant va vous aider a configurer la connexion\na votre serveur Unraid en quelques etapes :"},
	"step_enter_name":    {"en": "1. Name your server", "fr": "1. Nommer votre serveur"},
	"step_enter_url":     {"en": "2. Enter your server address", "fr": "2. Saisir l'adresse de votre serveur"},
	"step_test":          {"en": "3. Test the connection", "fr": "3. Tester la connexion"},
	"step_api_key":       {"en": "4. Configure your API key", "fr": "4. Configurer votre cle API"},
	"step_save":          {"en": "5. Save the configuration", "fr": "5. Sauvegarder la configuration"},
	"config_saved_in":    {"en": "Config will be saved in ~/.unraid-tui/config.yaml", "fr": "Le fichier sera sauvegarde dans ~/.unraid-tui/config.yaml"},
	"server_name_title":  {"en": "Server name", "fr": "Nom du serveur"},
	"server_name_desc":   {"en": "Give your server a name (e.g. NAS, Backup, Media).", "fr": "Donnez un nom a votre serveur (ex: NAS, Backup, Media)."},
	"server_name_hint":   {"en": "This name identifies the server in the list.", "fr": "Ce nom permet d'identifier le serveur dans la liste."},
	"server_name_empty":  {"en": "Server name cannot be empty", "fr": "Le nom du serveur ne peut pas etre vide"},
	"server_url_title":   {"en": "Unraid server address", "fr": "Adresse du serveur Unraid"},
	"server_url_desc":    {"en": "Enter the URL of your Unraid server (with port).", "fr": "Entrez l'URL de votre serveur Unraid (avec le port)."},
	"server_url_hint":    {"en": "By default, the Unraid API listens on port 3001.", "fr": "Par defaut, l'API Unraid ecoute sur le port 3001."},
	"server_url_empty":   {"en": "Server URL cannot be empty", "fr": "L'URL du serveur ne peut pas etre vide"},
	"testing_connection": {"en": "Testing connection to", "fr": "Test de la connexion a"},
	"testing_api_key":    {"en": "Verifying API key...", "fr": "Verification de la cle API..."},
	"saving_config":      {"en": "Saving configuration...", "fr": "Sauvegarde de la configuration..."},
	"api_key_title":      {"en": "Enter API key", "fr": "Saisir la cle API"},
	"api_key_desc":       {"en": "Paste your Unraid API key below.", "fr": "Collez votre cle API Unraid ci-dessous."},
	"api_key_hint":       {"en": "The key is masked for security.", "fr": "La cle est masquee pour des raisons de securite."},
	"api_key_empty":      {"en": "API key cannot be empty", "fr": "La cle API ne peut pas etre vide"},
	"api_key_info_title": {"en": "Create an API key", "fr": "Creer une cle API"},
	"api_key_howto":      {"en": "How to get an API key:", "fr": "Comment obtenir une cle API :"},
	"api_step_1":        {"en": "1. Open the Unraid web interface", "fr": "1. Ouvrez l'interface web de votre serveur Unraid"},
	"api_step_2":        {"en": "2. Go to Settings > Management Access", "fr": "2. Allez dans Settings > Management Access"},
	"api_step_3":        {"en": "3. Enable Developer Options", "fr": "3. Activez Developer Options"},
	"api_step_4":        {"en": "4. Open Apollo GraphQL Studio", "fr": "4. Ouvrez Apollo GraphQL Studio"},
	"api_step_5":        {"en": "5. Execute this mutation:", "fr": "5. Executez cette mutation :"},
	"api_step_6":        {"en": "6. Copy the returned key", "fr": "6. Copiez la cle retournee"},
	"config_done":        {"en": "Configuration complete!", "fr": "Configuration terminee !"},
	"config_saved_at":    {"en": "Your configuration has been saved in:", "fr": "Votre configuration a ete sauvegardee dans :"},
	"server_label":       {"en": "Server", "fr": "Serveur"},
	"api_key_label":      {"en": "API key", "fr": "Cle API"},
	"api_key_saved":      {"en": "********** (saved)", "fr": "********** (sauvegardee)"},
	"dash_will_launch":   {"en": "The dashboard will now launch.", "fr": "Le dashboard va maintenant se lancer."},
	"test_connection":    {"en": "test connection", "fr": "tester la connexion"},

	// === Server picker ===
	"server_picker_title": {"en": "Servers", "fr": "Serveurs"},
	"add_server":          {"en": "+ Add a server...", "fr": "+ Ajouter un serveur..."},

	// === Progress ===
	"step_name":       {"en": "Name", "fr": "Nom"},
	"step_url":        {"en": "URL", "fr": "URL"},
	"step_connection": {"en": "Connection", "fr": "Connexion"},
	"step_api":        {"en": "API Key", "fr": "Cle API"},
	"step_done":       {"en": "Done", "fr": "Termine"},
}

// SetLang sets the current language ("en" or "fr").
func SetLang(lang string) {
	lang = strings.ToLower(lang)
	if lang == "fr" || strings.HasPrefix(lang, "fr_") || strings.HasPrefix(lang, "fr-") {
		currentLang = "fr"
	} else {
		currentLang = "en"
	}
}

// DetectLang detects language from environment.
func DetectLang() {
	for _, env := range []string{"UNRAID_LANG", "LANG", "LC_ALL", "LC_MESSAGES"} {
		if v := os.Getenv(env); v != "" {
			SetLang(v)
			return
		}
	}
}

// Lang returns the current language code.
func Lang() string {
	return currentLang
}

// T translates a key to the current language.
func T(key string) string {
	if m, ok := translations[key]; ok {
		if v, ok := m[currentLang]; ok {
			return v
		}
		if v, ok := m["en"]; ok {
			return v
		}
	}
	return key
}
