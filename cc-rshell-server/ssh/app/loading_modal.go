package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type LoadingModal struct {
	*tview.TextView

	spinnerCharSet []string
	spinnerIndex   int
}

func (m *LoadingModal) SetText(text string) *LoadingModal {
	m.TextView.SetText(text + "\n\n" + m.spinnerCharSet[m.spinnerIndex])
	return m
}

func (m *LoadingModal) Draw(screen tcell.Screen) {
	txt := m.TextView.GetText(false)
	spinnerLength := len(m.spinnerCharSet[m.spinnerIndex])

	b := strings.Builder{}
	b.WriteString(txt[:len(txt)-1-spinnerLength])
	m.spinnerIndex = (m.spinnerIndex + 1) % len(m.spinnerCharSet)
	b.WriteString(m.spinnerCharSet[m.spinnerIndex])

	m.TextView.SetText(b.String())
	m.TextView.Draw(screen)
}
