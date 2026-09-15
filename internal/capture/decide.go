package capture

// DecideLogMessage determina a mensagem de log apropriada com base no evento do clipboard e na configuração de verbosidade.
func DecideLogMessage(ev Event, verbose bool) (message string, shouldLog bool) {
	if ev.IsImage {
		return "Screenshot detected", true
	}
	if verbose {
		return "Clipboard changed (non-image content)", true
	}
	return "", false
}