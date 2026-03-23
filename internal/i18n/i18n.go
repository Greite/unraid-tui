package i18n

import (
	"os"
	"strings"
)

var currentLang = "en"

// SupportedLanguages lists all supported language codes.
var SupportedLanguages = []string{"en", "fr", "zh", "hi", "es", "ar"}

var translations = map[string]map[string]string{
	// === General ===
	"loading":              {"en": "Loading...", "fr": "Chargement...", "zh": "加载中...", "hi": "लोड हो रहा है...", "es": "Cargando...", "ar": "...جارٍ التحميل"},
	"waiting":              {"en": "Waiting...", "fr": "En attente...", "zh": "等待中...", "hi": "प्रतीक्षा हो रही है...", "es": "Esperando...", "ar": "...في الانتظار"},
	"error":                {"en": "Error", "fr": "Erreur", "zh": "错误", "hi": "त्रुटि", "es": "Error", "ar": "خطأ"},
	"none":                 {"en": "None", "fr": "Aucun", "zh": "无", "hi": "कोई नहीं", "es": "Ninguno", "ar": "لا شيء"},
	"quit":                 {"en": "quit", "fr": "quitter", "zh": "退出", "hi": "बाहर", "es": "salir", "ar": "خروج"},
	"pages":                {"en": "pages", "fr": "pages", "zh": "页面", "hi": "पृष्ठ", "es": "páginas", "ar": "صفحات"},
	"next":                 {"en": "next", "fr": "suivant", "zh": "下一个", "hi": "अगला", "es": "siguiente", "ar": "التالي"},
	"servers":              {"en": "servers", "fr": "serveurs", "zh": "服务器", "hi": "सर्वर", "es": "servidores", "ar": "خوادم"},
	"navigate":             {"en": "navigate", "fr": "naviguer", "zh": "导航", "hi": "नेविगेट", "es": "navegar", "ar": "تنقل"},
	"refresh":              {"en": "refresh", "fr": "rafraîchir", "zh": "刷新", "hi": "रिफ्रेश", "es": "actualizar", "ar": "تحديث"},
	"sort":                 {"en": "sort", "fr": "trier", "zh": "排序", "hi": "क्रमबद्ध", "es": "ordenar", "ar": "ترتيب"},
	"back":                 {"en": "back", "fr": "retour", "zh": "返回", "hi": "वापस", "es": "volver", "ar": "رجوع"},
	"scroll":               {"en": "scroll", "fr": "scroll", "zh": "滚动", "hi": "स्क्रॉल", "es": "scroll", "ar": "تمرير"},
	"follow":               {"en": "follow", "fr": "follow", "zh": "跟踪", "hi": "फॉलो", "es": "seguir", "ar": "متابعة"},
	"start_end":            {"en": "start/end", "fr": "début/fin", "zh": "开始/结束", "hi": "शुरू/अंत", "es": "inicio/fin", "ar": "بداية/نهاية"},
	"execute":              {"en": "execute", "fr": "exécuter", "zh": "执行", "hi": "चलाएं", "es": "ejecutar", "ar": "تنفيذ"},
	"select":               {"en": "select", "fr": "sélectionner", "zh": "选择", "hi": "चुनें", "es": "seleccionar", "ar": "اختيار"},
	"connect":              {"en": "connect", "fr": "connecter", "zh": "连接", "hi": "कनेक्ट", "es": "conectar", "ar": "اتصال"},
	"default":              {"en": "default", "fr": "défaut", "zh": "默认", "hi": "डिफ़ॉल्ट", "es": "defecto", "ar": "افتراضي"},
	"delete":               {"en": "delete", "fr": "supprimer", "zh": "删除", "hi": "हटाएं", "es": "eliminar", "ar": "حذف"},
	"close":                {"en": "close", "fr": "fermer", "zh": "关闭", "hi": "बंद", "es": "cerrar", "ar": "إغلاق"},
	"archive":              {"en": "archive", "fr": "archiver", "zh": "归档", "hi": "संग्रह", "es": "archivar", "ar": "أرشفة"},
	"archive_all":          {"en": "archive all", "fr": "archiver tout", "zh": "全部归档", "hi": "सभी संग्रहित करें", "es": "archivar todo", "ar": "أرشفة الكل"},
	"begin":                {"en": "begin", "fr": "commencer", "zh": "开始", "hi": "शुरू करें", "es": "comenzar", "ar": "بدء"},
	"continue":             {"en": "continue", "fr": "continuer", "zh": "继续", "hi": "जारी रखें", "es": "continuar", "ar": "متابعة"},
	"validate":             {"en": "validate", "fr": "valider", "zh": "验证", "hi": "सत्यापित करें", "es": "validar", "ar": "تحقق"},
	"enter_key":            {"en": "enter key", "fr": "saisir la clé", "zh": "输入密钥", "hi": "कुंजी दर्ज करें", "es": "ingresar clave", "ar": "إدخال المفتاح"},
	"launch_dash":          {"en": "launch dashboard", "fr": "lancer le dashboard", "zh": "启动仪表盘", "hi": "डैशबोर्ड लॉन्च करें", "es": "abrir panel", "ar": "فتح لوحة التحكم"},
	"line":                 {"en": "line", "fr": "ligne", "zh": "行", "hi": "पंक्ति", "es": "línea", "ar": "سطر"},
	"reboot_server":        {"en": "reboot server", "fr": "redémarrer le serveur", "zh": "重启服务器", "hi": "सर्वर रीबूट", "es": "reiniciar servidor", "ar": "إعادة تشغيل الخادم"},
	"shutdown_server":      {"en": "shutdown server", "fr": "arrêter le serveur", "zh": "关闭服务器", "hi": "सर्वर बंद करें", "es": "apagar servidor", "ar": "إيقاف الخادم"},
	"confirm_reboot":       {"en": "Press Enter to confirm REBOOT, Esc to cancel", "fr": "Appuyez sur Entrée pour confirmer le REDÉMARRAGE, Echap pour annuler", "zh": "按 Enter 确认重启，Esc 取消", "hi": "रीबूट की पुष्टि के लिए Enter दबाएं, रद्द करने के लिए Esc", "es": "Presione Enter para confirmar REINICIO, Esc para cancelar", "ar": "اضغط Enter للتأكيد، Esc للإلغاء"},
	"confirm_shutdown":     {"en": "Press Enter to confirm SHUTDOWN, Esc to cancel", "fr": "Appuyez sur Entrée pour confirmer l'ARRÊT, Echap pour annuler", "zh": "按 Enter 确认关机，Esc 取消", "hi": "शटडाउन की पुष्टि के लिए Enter दबाएं, रद्द करने के लिए Esc", "es": "Presione Enter para confirmar APAGADO, Esc para cancelar", "ar": "اضغط Enter للتأكيد، Esc للإلغاء"},
	"server_rebooting":     {"en": "Server is rebooting...", "fr": "Le serveur redémarre...", "zh": "服务器正在重启...", "hi": "सर्वर रीबूट हो रहा है...", "es": "El servidor se está reiniciando...", "ar": "...الخادم يعيد التشغيل"},
	"server_shutting_down": {"en": "Server is shutting down...", "fr": "Le serveur s'arrête...", "zh": "服务器正在关闭...", "hi": "सर्वर बंद हो रहा है...", "es": "El servidor se está apagando...", "ar": "...الخادم يتوقف"},
	"syslog":               {"en": "syslog", "fr": "syslog", "zh": "系统日志", "hi": "सिसलॉग", "es": "syslog", "ar": "سجل النظام"},
	"lang":                 {"en": "language", "fr": "langue", "zh": "语言", "hi": "भाषा", "es": "idioma", "ar": "اللغة"},

	// === Pages ===
	"page_dashboard":     {"en": "Dashboard", "fr": "Dashboard", "zh": "仪表盘", "hi": "डैशबोर्ड", "es": "Panel", "ar": "لوحة التحكم"},
	"page_docker":        {"en": "Docker", "fr": "Docker", "zh": "Docker", "hi": "Docker", "es": "Docker", "ar": "Docker"},
	"page_vms":           {"en": "VMs", "fr": "VMs", "zh": "虚拟机", "hi": "VMs", "es": "VMs", "ar": "VMs"},
	"page_notifications": {"en": "Notifs", "fr": "Notifs", "zh": "通知", "hi": "सूचनाएं", "es": "Notifs", "ar": "إشعارات"},
	"page_shares":        {"en": "Shares", "fr": "Shares", "zh": "共享", "hi": "शेयर", "es": "Shares", "ar": "مشاركات"},
	"page_syslog":        {"en": "Syslog", "fr": "Syslog", "zh": "系统日志", "hi": "Syslog", "es": "Syslog", "ar": "Syslog"},

	// === Dashboard ===
	"system":    {"en": "System", "fr": "Système", "zh": "系统", "hi": "सिस्टम", "es": "Sistema", "ar": "النظام"},
	"hostname":  {"en": "Hostname", "fr": "Hostname", "zh": "主机名", "hi": "होस्टनाम", "es": "Hostname", "ar": "اسم المضيف"},
	"uptime":    {"en": "Uptime", "fr": "Uptime", "zh": "运行时间", "hi": "अपटाइम", "es": "Tiempo activo", "ar": "وقت التشغيل"},
	"cpu":       {"en": "CPU", "fr": "CPU", "zh": "CPU", "hi": "CPU", "es": "CPU", "ar": "CPU"},
	"cpu_cores": {"en": "CPU Cores", "fr": "CPU Cores", "zh": "CPU 核心", "hi": "CPU कोर", "es": "Núcleos CPU", "ar": "أنوية CPU"},
	"memory":    {"en": "Memory", "fr": "Mémoire", "zh": "内存", "hi": "मेमोरी", "es": "Memoria", "ar": "الذاكرة"},
	"network":   {"en": "Network", "fr": "Réseau", "zh": "网络", "hi": "नेटवर्क", "es": "Red", "ar": "الشبكة"},
	"disks":     {"en": "Disks", "fr": "Disques", "zh": "磁盘", "hi": "डिस्क", "es": "Discos", "ar": "الأقراص"},
	"hardware":  {"en": "Hardware", "fr": "Matériel", "zh": "硬件", "hi": "हार्डवेयर", "es": "Hardware", "ar": "العتاد"},
	"parity":    {"en": "Parity", "fr": "Parité", "zh": "校验", "hi": "पैरिटी", "es": "Paridad", "ar": "التماثل"},
	"array":     {"en": "Array", "fr": "Array", "zh": "阵列", "hi": "ऐरे", "es": "Array", "ar": "المصفوفة"},
	"total":     {"en": "total", "fr": "total", "zh": "总计", "hi": "कुल", "es": "total", "ar": "إجمالي"},
	"running":   {"en": "running", "fr": "en cours", "zh": "运行中", "hi": "चल रहा है", "es": "en ejecución", "ar": "قيد التشغيل"},
	"exited":    {"en": "exited", "fr": "arrêté", "zh": "已停止", "hi": "बंद हो गया", "es": "detenido", "ar": "متوقف"},
	"paused":    {"en": "paused", "fr": "pause", "zh": "已暂停", "hi": "रुका हुआ", "es": "pausado", "ar": "متوقف مؤقتاً"},
	"devices":   {"en": "devices", "fr": "périphériques", "zh": "设备", "hi": "डिवाइस", "es": "dispositivos", "ar": "أجهزة"},

	// === Parity History ===
	"parity_history":    {"en": "Parity History", "fr": "Historique parité", "zh": "校验历史", "hi": "पैरिटी इतिहास", "es": "Historial de paridad", "ar": "سجل التماثل"},
	"parity_date":       {"en": "Date", "fr": "Date", "zh": "日期", "hi": "तारीख", "es": "Fecha", "ar": "التاريخ"},
	"parity_status":     {"en": "Status", "fr": "Statut", "zh": "状态", "hi": "स्थिति", "es": "Estado", "ar": "الحالة"},
	"parity_duration":   {"en": "Duration", "fr": "Durée", "zh": "持续时间", "hi": "अवधि", "es": "Duración", "ar": "المدة"},
	"parity_speed":      {"en": "Speed", "fr": "Vitesse", "zh": "速度", "hi": "गति", "es": "Velocidad", "ar": "السرعة"},
	"parity_errors":     {"en": "Errors", "fr": "Erreurs", "zh": "错误", "hi": "त्रुटियां", "es": "Errores", "ar": "الأخطاء"},
	"parity_start":      {"en": "start check", "fr": "lancer le check", "zh": "开始检查", "hi": "जांच शुरू", "es": "iniciar check", "ar": "بدء الفحص"},
	"parity_pause":      {"en": "pause", "fr": "pause", "zh": "暂停", "hi": "रोकें", "es": "pausar", "ar": "إيقاف مؤقت"},
	"parity_resume":     {"en": "resume", "fr": "reprendre", "zh": "恢复", "hi": "जारी रखें", "es": "reanudar", "ar": "استئناف"},
	"parity_cancel":     {"en": "cancel", "fr": "annuler", "zh": "取消", "hi": "रद्द करें", "es": "cancelar", "ar": "إلغاء"},
	"no_parity_history": {"en": "No parity checks recorded", "fr": "Aucun check de parité enregistré", "zh": "无校验记录", "hi": "कोई पैरिटी जांच दर्ज नहीं", "es": "Sin verificaciones de paridad", "ar": "لا توجد سجلات تماثل"},

	// === Docker ===
	"containers":      {"en": "Containers", "fr": "Containers", "zh": "容器", "hi": "कंटेनर", "es": "Contenedores", "ar": "حاويات"},
	"loading_docker":  {"en": "Loading containers...", "fr": "Chargement des containers...", "zh": "加载容器中...", "hi": "कंटेनर लोड हो रहे हैं...", "es": "Cargando contenedores...", "ar": "...جارٍ تحميل الحاويات"},
	"docker_disabled": {"en": "Docker is not enabled on this server.", "fr": "Docker n'est pas activé sur ce serveur.", "zh": "此服务器未启用 Docker。", "hi": "इस सर्वर पर Docker सक्षम नहीं है।", "es": "Docker no está habilitado en este servidor.", "ar": ".Docker غير مفعل على هذا الخادم"},
	"docker_enable":   {"en": "Enable it in Settings > Docker.", "fr": "Activez-le dans Settings > Docker.", "zh": "请在 Settings > Docker 中启用。", "hi": "इसे Settings > Docker में सक्षम करें।", "es": "Actívelo en Settings > Docker.", "ar": ".Settings > Docker قم بتفعيله من"},
	"logs":            {"en": "logs", "fr": "logs", "zh": "日志", "hi": "लॉग", "es": "logs", "ar": "سجلات"},
	"console":         {"en": "console", "fr": "console", "zh": "控制台", "hi": "कंसोल", "es": "consola", "ar": "وحدة التحكم"},
	"webui":           {"en": "WebUI", "fr": "WebUI", "zh": "WebUI", "hi": "WebUI", "es": "WebUI", "ar": "WebUI"},
	"stats":           {"en": "stats", "fr": "stats", "zh": "统计", "hi": "आंकड़े", "es": "stats", "ar": "إحصائيات"},
	"start":           {"en": "start", "fr": "démarrer", "zh": "启动", "hi": "शुरू करें", "es": "iniciar", "ar": "تشغيل"},
	"stop":            {"en": "stop", "fr": "arrêter", "zh": "停止", "hi": "रोकें", "es": "detener", "ar": "إيقاف"},
	"pause":           {"en": "pause", "fr": "pause", "zh": "暂停", "hi": "रोकें", "es": "pausar", "ar": "إيقاف مؤقت"},
	"unpause":         {"en": "unpause", "fr": "reprendre", "zh": "恢复", "hi": "फिर से शुरू करें", "es": "reanudar", "ar": "استئناف"},
	"autostart":       {"en": "autostart", "fr": "démarrage auto", "zh": "自动启动", "hi": "ऑटोस्टार्ट", "es": "inicio auto", "ar": "تشغيل تلقائي"},
	"autostart_on":    {"en": "autostart ON", "fr": "démarrage auto ON", "zh": "自动启动 开", "hi": "ऑटोस्टार्ट चालू", "es": "inicio auto ON", "ar": "تشغيل تلقائي مفعل"},
	"autostart_off":   {"en": "autostart OFF", "fr": "démarrage auto OFF", "zh": "自动启动 关", "hi": "ऑटोस्टार्ट बंद", "es": "inicio auto OFF", "ar": "تشغيل تلقائي معطل"},
	"update":          {"en": "update", "fr": "mettre à jour", "zh": "更新", "hi": "अपडेट", "es": "actualizar", "ar": "تحديث"},
	"update_all":      {"en": "update all", "fr": "tout mettre à jour", "zh": "全部更新", "hi": "सभी अपडेट करें", "es": "actualizar todo", "ar": "تحديث الكل"},
	"up_to_date":      {"en": "already up to date", "fr": "déjà à jour", "zh": "已是最新", "hi": "पहले से अपडेट है", "es": "ya actualizado", "ar": "محدّث بالفعل"},
	"updating":        {"en": "Updating", "fr": "Mise à jour de", "zh": "正在更新", "hi": "अपडेट हो रहा है", "es": "Actualizando", "ar": "جارٍ التحديث"},
	"updating_all":    {"en": "Updating all containers", "fr": "Mise à jour de tous les containers", "zh": "正在更新所有容器", "hi": "सभी कंटेनर अपडेट हो रहे हैं", "es": "Actualizando todos los contenedores", "ar": "جارٍ تحديث جميع الحاويات"},
	"no_webui":        {"en": "No WebUI for %s", "fr": "Pas de WebUI pour %s", "zh": "%s 没有 WebUI", "hi": "%s के लिए कोई WebUI नहीं", "es": "Sin WebUI para %s", "ar": "WebUI لا يوجد لـ %s"},
	"webui_opened":    {"en": "WebUI opened for %s", "fr": "WebUI ouvert pour %s", "zh": "已打开 %s 的 WebUI", "hi": "%s का WebUI खोला गया", "es": "WebUI abierta para %s", "ar": "WebUI تم فتح لـ %s"},
	"not_running":     {"en": "%s is not running", "fr": "%s n'est pas running", "zh": "%s 未运行", "hi": "%s नहीं चल रहा है", "es": "%s no está en ejecución", "ar": "%s لا يعمل"},
	"console_done":    {"en": "Console finished", "fr": "Console terminée", "zh": "控制台已结束", "hi": "कंसोल समाप्त", "es": "Consola terminada", "ar": "انتهت وحدة التحكم"},
	"console_error":   {"en": "Console finished with error", "fr": "Console terminée avec erreur", "zh": "控制台异常结束", "hi": "कंसोल त्रुटि के साथ समाप्त", "es": "Consola terminada con error", "ar": "انتهت وحدة التحكم مع خطأ"},
	"connected_to":    {"en": "Connected to %s via SSH", "fr": "Connecté à %s via SSH", "zh": "已通过 SSH 连接到 %s", "hi": "SSH द्वारा %s से कनेक्ट हो गया", "es": "Conectado a %s via SSH", "ar": "SSH متصل بـ %s عبر"},
	"logs_error":      {"en": "Logs error: %s", "fr": "Erreur logs: %s", "zh": "日志错误: %s", "hi": "लॉग त्रुटि: %s", "es": "Error de logs: %s", "ar": "خطأ في السجلات: %s"},
	"action_ok":       {"en": "%s %s OK", "fr": "%s %s OK", "zh": "%s %s 成功", "hi": "%s %s ठीक", "es": "%s %s OK", "ar": "%s %s تم"},
	"action_error":    {"en": "Error %s %s: %s", "fr": "Erreur %s %s: %s", "zh": "%s %s 错误: %s", "hi": "त्रुटि %s %s: %s", "es": "Error %s %s: %s", "ar": "خطأ %s %s: %s"},
	"follow_on":       {"en": "FOLLOW", "fr": "SUIVI", "zh": "跟踪", "hi": "फॉलो", "es": "SEGUIR", "ar": "متابعة"},
	"follow_off":      {"en": "PAUSE", "fr": "PAUSE", "zh": "暂停", "hi": "रुकें", "es": "PAUSA", "ar": "إيقاف"},

	// === VMs ===
	"loading_vms":  {"en": "Loading VMs...", "fr": "Chargement des VMs...", "zh": "加载虚拟机中...", "hi": "VMs लोड हो रहे हैं...", "es": "Cargando VMs...", "ar": "...جارٍ تحميل الأجهزة الافتراضية"},
	"no_vms":       {"en": "No VMs configured", "fr": "Aucune VM configurée", "zh": "未配置虚拟机", "hi": "कोई VM कॉन्फ़िगर नहीं है", "es": "No hay VMs configuradas", "ar": "لا توجد أجهزة افتراضية"},
	"vms_disabled": {"en": "VMs are not enabled on this server.", "fr": "Les VMs ne sont pas activées sur ce serveur.", "zh": "此服务器未启用虚拟机。", "hi": "इस सर्वर पर VMs सक्षम नहीं हैं।", "es": "Las VMs no están habilitadas en este servidor.", "ar": ".الأجهزة الافتراضية غير مفعلة على هذا الخادم"},
	"vms_enable":   {"en": "Enable them in Settings > VM Manager.", "fr": "Activez-les dans Settings > VM Manager.", "zh": "请在 Settings > VM Manager 中启用。", "hi": "इन्हें Settings > VM Manager में सक्षम करें।", "es": "Actívelas en Settings > VM Manager.", "ar": ".Settings > VM Manager قم بتفعيلها من"},
	"reboot":       {"en": "reboot", "fr": "redémarrer", "zh": "重启", "hi": "रीबूट", "es": "reiniciar", "ar": "إعادة تشغيل"},
	"force_stop":   {"en": "force stop", "fr": "forcer l'arrêt", "zh": "强制停止", "hi": "बलपूर्वक रोकें", "es": "forzar detención", "ar": "إيقاف إجباري"},
	"resume":       {"en": "resume", "fr": "reprendre", "zh": "恢复", "hi": "फिर से शुरू करें", "es": "reanudar", "ar": "استئناف"},

	// === Notifications ===
	"loading_notifs": {"en": "Loading notifications...", "fr": "Chargement des notifications...", "zh": "加载通知中...", "hi": "सूचनाएं लोड हो रही हैं...", "es": "Cargando notificaciones...", "ar": "...جارٍ تحميل الإشعارات"},
	"no_notifs":      {"en": "No unread notifications", "fr": "Aucune notification non lue", "zh": "没有未读通知", "hi": "कोई अपठित सूचना नहीं", "es": "Sin notificaciones sin leer", "ar": "لا توجد إشعارات غير مقروءة"},
	"notifications":  {"en": "Notifications", "fr": "Notifications", "zh": "通知", "hi": "सूचनाएं", "es": "Notificaciones", "ar": "إشعارات"},

	// === Shares ===
	"loading_shares": {"en": "Loading shares...", "fr": "Chargement des shares...", "zh": "加载共享中...", "hi": "शेयर लोड हो रहे हैं...", "es": "Cargando shares...", "ar": "...جارٍ تحميل المشاركات"},
	"no_shares":      {"en": "No shares configured", "fr": "Aucun share configuré", "zh": "未配置共享", "hi": "कोई शेयर कॉन्फ़िगर नहीं है", "es": "No hay shares configurados", "ar": "لا توجد مشاركات"},
	"shares":         {"en": "Shares", "fr": "Shares", "zh": "共享", "hi": "शेयर", "es": "Shares", "ar": "مشاركات"},

	// === Onboarding ===
	"onboarding_title": {"en": "UNRAID TUI — Configuration", "fr": "UNRAID TUI — Configuration", "zh": "UNRAID TUI — 配置", "hi": "UNRAID TUI — कॉन्फ़िगरेशन", "es": "UNRAID TUI — Configuracion", "ar": "UNRAID TUI — الإعداد"},
	"welcome":          {"en": "Welcome!", "fr": "Bienvenue !", "zh": "欢迎！", "hi": "स्वागत है!", "es": "¡Bienvenido!", "ar": "!مرحباً"},
	"welcome_desc": {
		"en": "This wizard will help you configure the connection\nto your Unraid server in a few steps:",
		"fr": "Cet assistant va vous aider à configurer la connexion\nà votre serveur Unraid en quelques étapes :",
		"zh": "此向导将帮助您通过几个步骤配置\n与 Unraid 服务器的连接：",
		"hi": "यह विज़ार्ड कुछ चरणों में आपके Unraid सर्वर\nसे कनेक्शन कॉन्फ़िगर करने में मदद करेगा:",
		"es": "Este asistente le ayudará a configurar la conexión\na su servidor Unraid en unos pasos:",
		"ar": ":سيساعدك هذا المعالج على إعداد الاتصال\nبخادم Unraid الخاص بك في خطوات قليلة",
	},
	"step_enter_name":    {"en": "1. Name your server", "fr": "1. Nommer votre serveur", "zh": "1. 命名您的服务器", "hi": "1. अपने सर्वर का नाम दें", "es": "1. Nombre su servidor", "ar": "1. قم بتسمية خادمك"},
	"step_enter_url":     {"en": "2. Enter your server address", "fr": "2. Saisir l'adresse de votre serveur", "zh": "2. 输入服务器地址", "hi": "2. सर्वर का पता दर्ज करें", "es": "2. Ingrese la dirección del servidor", "ar": "2. أدخل عنوان الخادم"},
	"step_test":          {"en": "3. Test the connection", "fr": "3. Tester la connexion", "zh": "3. 测试连接", "hi": "3. कनेक्शन का परीक्षण करें", "es": "3. Probar la conexión", "ar": "3. اختبر الاتصال"},
	"step_api_key":       {"en": "4. Configure your API key", "fr": "4. Configurer votre clé API", "zh": "4. 配置 API 密钥", "hi": "4. अपनी API कुंजी कॉन्फ़िगर करें", "es": "4. Configure su clave API", "ar": "4. قم بإعداد مفتاح API"},
	"step_save":          {"en": "5. Save the configuration", "fr": "5. Sauvegarder la configuration", "zh": "5. 保存配置", "hi": "5. कॉन्फ़िगरेशन सहेजें", "es": "5. Guardar la configuración", "ar": "5. احفظ الإعدادات"},
	"config_saved_in":    {"en": "Config will be saved in ~/.unraid-tui/config.yaml", "fr": "Le fichier sera sauvegardé dans ~/.unraid-tui/config.yaml", "zh": "配置将保存到 ~/.unraid-tui/config.yaml", "hi": "कॉन्फ़िगरेशन ~/.unraid-tui/config.yaml में सहेजा जाएगा", "es": "La configuración se guardará en ~/.unraid-tui/config.yaml", "ar": "~/.unraid-tui/config.yaml سيتم حفظ الإعدادات في"},
	"server_name_title":  {"en": "Server name", "fr": "Nom du serveur", "zh": "服务器名称", "hi": "सर्वर का नाम", "es": "Nombre del servidor", "ar": "اسم الخادم"},
	"server_name_desc":   {"en": "Give your server a name (e.g. NAS, Backup, Media).", "fr": "Donnez un nom à votre serveur (ex: NAS, Backup, Media).", "zh": "为您的服务器命名（如 NAS、Backup、Media）。", "hi": "अपने सर्वर को एक नाम दें (जैसे NAS, Backup, Media)।", "es": "Dele un nombre a su servidor (ej: NAS, Backup, Media).", "ar": ".(NAS, Backup, Media :أعطِ خادمك اسماً (مثل"},
	"server_name_hint":   {"en": "This name identifies the server in the list.", "fr": "Ce nom permet d'identifier le serveur dans la liste.", "zh": "此名称用于在列表中标识服务器。", "hi": "यह नाम सूची में सर्वर की पहचान करता है।", "es": "Este nombre identifica el servidor en la lista.", "ar": ".هذا الاسم يحدد الخادم في القائمة"},
	"server_name_empty":  {"en": "Server name cannot be empty", "fr": "Le nom du serveur ne peut pas être vide", "zh": "服务器名称不能为空", "hi": "सर्वर का नाम खाली नहीं हो सकता", "es": "El nombre del servidor no puede estar vacío", "ar": "لا يمكن أن يكون اسم الخادم فارغاً"},
	"server_url_title":   {"en": "Unraid server address", "fr": "Adresse du serveur Unraid", "zh": "Unraid 服务器地址", "hi": "Unraid सर्वर का पता", "es": "Dirección del servidor Unraid", "ar": "Unraid عنوان خادم"},
	"server_url_desc":    {"en": "Enter the URL of your Unraid server (with port).", "fr": "Entrez l'URL de votre serveur Unraid (avec le port).", "zh": "输入 Unraid 服务器的 URL（含端口）。", "hi": "अपने Unraid सर्वर का URL दर्ज करें (पोर्ट के साथ)।", "es": "Ingrese la URL de su servidor Unraid (con puerto).", "ar": ".(أدخل عنوان URL لخادم Unraid (مع المنفذ"},
	"server_url_hint":    {"en": "By default, the Unraid API listens on port 3001.", "fr": "Par défaut, l'API Unraid écoute sur le port 3001.", "zh": "默认情况下，Unraid API 监听端口 3001。", "hi": "डिफ़ॉल्ट रूप से, Unraid API पोर्ट 3001 पर सुनता है।", "es": "Por defecto, la API de Unraid escucha en el puerto 3001.", "ar": ".3001 بشكل افتراضي، واجهة Unraid API تستمع على المنفذ"},
	"server_url_empty":   {"en": "Server URL cannot be empty", "fr": "L'URL du serveur ne peut pas être vide", "zh": "服务器 URL 不能为空", "hi": "सर्वर URL खाली नहीं हो सकता", "es": "La URL del servidor no puede estar vacía", "ar": "لا يمكن أن يكون عنوان URL فارغاً"},
	"testing_connection": {"en": "Testing connection to", "fr": "Test de la connexion à", "zh": "正在测试连接到", "hi": "कनेक्शन का परीक्षण कर रहे हैं", "es": "Probando conexión a", "ar": "جارٍ اختبار الاتصال بـ"},
	"testing_api_key":    {"en": "Verifying API key...", "fr": "Vérification de la clé API...", "zh": "验证 API 密钥中...", "hi": "API कुंजी सत्यापित हो रही है...", "es": "Verificando clave API...", "ar": "...جارٍ التحقق من مفتاح API"},
	"saving_config":      {"en": "Saving configuration...", "fr": "Sauvegarde de la configuration...", "zh": "保存配置中...", "hi": "कॉन्फ़िगरेशन सहेजा जा रहा है...", "es": "Guardando configuracion...", "ar": "...جارٍ حفظ الإعدادات"},
	"api_key_title":      {"en": "Enter API key", "fr": "Saisir la clé API", "zh": "输入 API 密钥", "hi": "API कुंजी दर्ज करें", "es": "Ingresar clave API", "ar": "API أدخل مفتاح"},
	"api_key_desc":       {"en": "Paste your Unraid API key below.", "fr": "Collez votre clé API Unraid ci-dessous.", "zh": "请在下方粘贴您的 Unraid API 密钥。", "hi": "अपनी Unraid API कुंजी नीचे चिपकाएं।", "es": "Pegue su clave API de Unraid a continuacion.", "ar": ".أدناه Unraid API الصق مفتاح"},
	"api_key_hint":       {"en": "The key is masked for security.", "fr": "La clé est masquée pour des raisons de sécurité.", "zh": "密钥已隐藏以确保安全。", "hi": "सुरक्षा के लिए कुंजी छिपी हुई है।", "es": "La clave está oculta por seguridad.", "ar": ".المفتاح مخفي لأسباب أمنية"},
	"api_key_empty":      {"en": "API key cannot be empty", "fr": "La clé API ne peut pas être vide", "zh": "API 密钥不能为空", "hi": "API कुंजी खाली नहीं हो सकती", "es": "La clave API no puede estar vacía", "ar": "فارغاً API لا يمكن أن يكون مفتاح"},
	"api_key_info_title": {"en": "Create an API key", "fr": "Créer une clé API", "zh": "创建 API 密钥", "hi": "API कुंजी बनाएं", "es": "Crear una clave API", "ar": "API إنشاء مفتاح"},
	"api_key_howto":      {"en": "How to get an API key:", "fr": "Comment obtenir une clé API :", "zh": "如何获取 API 密钥：", "hi": "API कुंजी कैसे प्राप्त करें:", "es": "Cómo obtener una clave API:", "ar": ":API كيفية الحصول على مفتاح"},
	"api_step_1":         {"en": "1. Open the Unraid web interface", "fr": "1. Ouvrez l'interface web de votre serveur Unraid", "zh": "1. 打开 Unraid Web 界面", "hi": "1. Unraid वेब इंटरफेस खोलें", "es": "1. Abra la interfaz web de Unraid", "ar": "Unraid 1. افتح واجهة ويب"},
	"api_step_2":         {"en": "2. Go to Settings > Management Access", "fr": "2. Allez dans Settings > Management Access", "zh": "2. 前往 Settings > Management Access", "hi": "2. Settings > Management Access पर जाएं", "es": "2. Vaya a Settings > Management Access", "ar": "Settings > Management Access 2. انتقل إلى"},
	"api_step_3":         {"en": "3. Enable Developer Options", "fr": "3. Activez Developer Options", "zh": "3. 启用 Developer Options", "hi": "3. Developer Options सक्षम करें", "es": "3. Active Developer Options", "ar": "Developer Options 3. قم بتفعيل"},
	"api_step_4":         {"en": "4. Open Apollo GraphQL Studio", "fr": "4. Ouvrez Apollo GraphQL Studio", "zh": "4. 打开 Apollo GraphQL Studio", "hi": "4. Apollo GraphQL Studio खोलें", "es": "4. Abra Apollo GraphQL Studio", "ar": "Apollo GraphQL Studio 4. افتح"},
	"api_step_5":         {"en": "5. Execute this mutation:", "fr": "5. Exécutez cette mutation :", "zh": "5. 执行此 mutation：", "hi": "5. यह mutation चलाएं:", "es": "5. Ejecute esta mutación:", "ar": ":mutation 5. نفّذ هذا الـ"},
	"api_step_6":         {"en": "6. Copy the returned key", "fr": "6. Copiez la clé retournée", "zh": "6. 复制返回的密钥", "hi": "6. लौटाई गई कुंजी कॉपी करें", "es": "6. Copie la clave devuelta", "ar": "6. انسخ المفتاح المُرجع"},
	"config_done":        {"en": "Configuration complete!", "fr": "Configuration terminée !", "zh": "配置完成！", "hi": "कॉन्फ़िगरेशन पूर्ण!", "es": "¡Configuración completa!", "ar": "!اكتمل الإعداد"},
	"config_saved_at":    {"en": "Your configuration has been saved in:", "fr": "Votre configuration a été sauvegardée dans :", "zh": "您的配置已保存到：", "hi": "आपका कॉन्फ़िगरेशन यहां सहेजा गया है:", "es": "Su configuración ha sido guardada en:", "ar": ":تم حفظ الإعدادات في"},
	"server_label":       {"en": "Server", "fr": "Serveur", "zh": "服务器", "hi": "सर्वर", "es": "Servidor", "ar": "الخادم"},
	"api_key_label":      {"en": "API key", "fr": "Clé API", "zh": "API 密钥", "hi": "API कुंजी", "es": "Clave API", "ar": "API مفتاح"},
	"api_key_saved":      {"en": "********** (saved)", "fr": "********** (sauvegardée)", "zh": "********** (已保存)", "hi": "********** (सहेजा गया)", "es": "********** (guardada)", "ar": "(تم الحفظ) **********"},
	"dash_will_launch":   {"en": "The dashboard will now launch.", "fr": "Le dashboard va maintenant se lancer.", "zh": "仪表盘即将启动。", "hi": "डैशबोर्ड अब लॉन्च होगा।", "es": "El panel se lanzará ahora.", "ar": ".سيتم فتح لوحة التحكم الآن"},
	"test_connection":    {"en": "test connection", "fr": "tester la connexion", "zh": "测试连接", "hi": "कनेक्शन परीक्षण", "es": "probar conexión", "ar": "اختبار الاتصال"},

	// === Server picker ===
	"server_picker_title": {"en": "Servers", "fr": "Serveurs", "zh": "服务器", "hi": "सर्वर", "es": "Servidores", "ar": "الخوادم"},
	"add_server":          {"en": "+ Add a server...", "fr": "+ Ajouter un serveur...", "zh": "+ 添加服务器...", "hi": "+ सर्वर जोड़ें...", "es": "+ Agregar un servidor...", "ar": "...+ إضافة خادم"},

	// === Progress ===
	"step_name":       {"en": "Name", "fr": "Nom", "zh": "名称", "hi": "नाम", "es": "Nombre", "ar": "الاسم"},
	"step_url":        {"en": "URL", "fr": "URL", "zh": "URL", "hi": "URL", "es": "URL", "ar": "URL"},
	"step_connection": {"en": "Connection", "fr": "Connexion", "zh": "连接", "hi": "कनेक्शन", "es": "Conexión", "ar": "الاتصال"},
	"step_api":        {"en": "API Key", "fr": "Clé API", "zh": "API 密钥", "hi": "API कुंजी", "es": "Clave API", "ar": "API مفتاح"},
	"step_done":       {"en": "Done", "fr": "Terminé", "zh": "完成", "hi": "पूर्ण", "es": "Listo", "ar": "تم"},
}

// SetLang sets the current language.
func SetLang(lang string) {
	lang = strings.ToLower(lang)
	for _, supported := range SupportedLanguages {
		if lang == supported || strings.HasPrefix(lang, supported+"_") || strings.HasPrefix(lang, supported+"-") {
			currentLang = supported
			return
		}
	}
	currentLang = "en"
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
