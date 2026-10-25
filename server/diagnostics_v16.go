package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var engineVersionRE = regexp.MustCompile(`CVPN webOS (?:Engine )?([0-9]+\.[0-9]+\.[0-9]+)`)

func engineVersion() string {
	b, err := os.ReadFile(path("bin", "cvpn.sh"))
	if err != nil {
		return "CVPN Engine (версия неизвестна)"
	}
	if len(b) > 600 {
		b = b[:600]
	}
	if match := engineVersionRE.FindSubmatch(b); len(match) == 2 {
		return "CVPN Engine " + string(match[1])
	}
	return "CVPN Engine (версия неизвестна)"
}

// Only fixed, non-secret strings go to the UI. In particular, Xray log lines
// may include credentials or full subscription URLs; NEVER expose raw logs.
func summarizeLogFile(name string) string {
	b, err := os.ReadFile(path("state", name))
	if os.IsNotExist(err) {
		return "файл отсутствует"
	}
	if err != nil {
		return "не удалось прочитать"
	}
	if len(b) > 100000 {
		b = b[len(b)-100000:]
	}
	content := strings.ToLower(string(b))
	if len(content) == 0 {
		return "записей пока нет"
	}
	lines := strings.Count(content, "\n")
	if !strings.HasSuffix(content, "\n") {
		lines++
	}
	errors := 0
	warnings := 0
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "error") || strings.Contains(line, "failed") || strings.Contains(line, "ошиб") {
			errors++
		}
		if strings.Contains(line, "warn") || strings.Contains(line, "restor") || strings.Contains(line, "откат") {
			warnings++
		}
	}
	return fmt.Sprintf("%d строк, %d сообщений с ошибками, %d предупреждений/восстановлений", lines, errors, warnings)
}

func safeEngineLogSummary() string {
	return "СВОДКА ЖУРНАЛОВ (без секретов)\n" +
		"Xray: " + summarizeLogFile("xray.log") + "\n" +
		"Watchdog: " + summarizeLogFile("watchdog.log") + "\n" +
		"Восстановление после запуска: " + safeBootRecoverySummary() + "\n" +
		"Таймер безопасности: " + summarizeLogFile("timer.log") + "\n" +
		"Панель: " + summarizePanelLog() + "\n" +
		"Для полной локальной диагностики администратор может изучить журналы по SSH. Не публикуйте их без удаления токенов и UUID."
}

func summarizePanelLog() string {
	b, e := os.ReadFile(path("webui", "state", "panel.log"))
	if e != nil {
		return "файл отсутствует или недоступен"
	}
	if len(b) > 100000 {
		b = b[len(b)-100000:]
	}
	count := strings.Count(string(b), "\n")
	if len(b) > 0 && b[len(b)-1] != '\n' {
		count++
	}
	return fmt.Sprintf("%d строк (содержимое скрыто)", count)
}

// On-demand network check. This is not run in the status polling loop and
// never starts/stops VPN or changes routing. It is advisory, not a proof of
// privacy for every application on the television.
func httpsProbe() (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	args := []string{"-4", "-sS", "--noproxy", "*", "--connect-timeout", "4", "--max-time", "7", "-o", "/dev/null", "-w", "%{http_code}", "https://www.google.com/generate_204"}
	out, err := exec.CommandContext(ctx, "curl", args...).Output()
	if ctx.Err() != nil {
		return "timeout", "HTTPS-проверка превысила 8 секунд"
	}
	if err != nil {
		return "failed", "HTTPS-проверка не прошла (DNS, TLS или маршрут)"
	}
	if strings.TrimSpace(string(out)) == "204" {
		return "ok", "Google HTTPS: HTTP 204"
	}
	return "unexpected", "Сервер ответил, но код HTTP отличается от 204"
}
