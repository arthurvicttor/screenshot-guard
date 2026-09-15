package capture

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sys/windows"
)

// Constantes de formatos de clipboard do Windows.
const (
	cfBitmap = 2  // CF_BITMAP
	cfDIB    = 8  // CF_DIB
	cfDIBV5  = 17 // CF_DIBV5
)

var (
	user32                          = windows.NewLazySystemDLL("user32.dll")
	procGetClipboardSequenceNumber  = user32.NewProc("GetClipboardSequenceNumber")
	procIsClipboardFormatAvailable  = user32.NewProc("IsClipboardFormatAvailable")
)

// Event representa um evento de mudança no clipboard, indicando se o conteúdo é uma imagem ou não.
type Event struct {
	Time    time.Time
	IsImage bool
}

// ClipboardWatcher é responsável por monitorar mudanças no clipboard do Windows.
type ClipboardWatcher struct {
	Interval time.Duration
	Events   chan Event
}

// NewClipboardWatcher cria um novo ClipboardWatcher com o intervalo de polling especificado.
func NewClipboardWatcher(interval time.Duration) *ClipboardWatcher {
	return &ClipboardWatcher{
		Interval: interval,
		Events:   make(chan Event),
	}
}

// Run inicia o monitoramento do clipboard, enviando eventos para o canal Events quando mudanças são detectadas.
func (w *ClipboardWatcher) Run(ctx context.Context) error {
	defer close(w.Events)

	lastSeq, err := getClipboardSequenceNumber()
	if err != nil {
		return fmt.Errorf("capture: leitura inicial do clipboard falhou: %w", err)
	}

	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			seq, err := getClipboardSequenceNumber()
			if err != nil {
				continue // falha pontual: ignora este tick e tenta de novo no próximo
			}
			if seq == lastSeq {
				continue
			}
			lastSeq = seq

			ev := Event{
				Time:    time.Now(),
				IsImage: isImageFormatAvailable(),
			}
			select {
			case w.Events <- ev:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func getClipboardSequenceNumber() (uint32, error) {
	r, _, callErr := procGetClipboardSequenceNumber.Call()
	if r == 0 {
		// Se o retorno for 0, pode ser que o clipboard não esteja disponível ou que tenha ocorrido algum erro.
		return 0, fmt.Errorf("GetClipboardSequenceNumber retornou 0 (%v)", callErr)
	}
	return uint32(r), nil
}

func isImageFormatAvailable() bool {
	for _, format := range []uintptr{cfDIB, cfDIBV5, cfBitmap} {
		r, _, _ := procIsClipboardFormatAvailable.Call(format)
		if r != 0 {
			return true
		}
	}
	return false
}