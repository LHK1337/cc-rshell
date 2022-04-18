package app

import (
	"cc-rshell-server/sockets/types"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
	"math/rand"
)

var spinnerIDPool = []int{
	1, 14, 24, 27, 32, 43, 47, 50, 54, 57, 69, 78, 79, 80,
}

func RunApp(screen tcell.Screen, registry types.ClientRegistry, name string, pubKey ssh.PublicKey) error {
	isAnon := name == ""
	if isAnon {
		name = "anonymous"
	}

	app := tview.NewApplication().SetScreen(screen)

	txt := tview.NewTextView().SetDynamicColors(true)
	txt.SetTitle(" Welcome stranger! ")
	txt.SetText(fmt.Sprintf("\nIt appears you are [lime]%s[white].\nNice to meet you 👋", name)).
		SetTextAlign(tview.AlignCenter).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyESC {
				app.Stop()
			}
		}).
		SetBorder(true)

	spinnerSetID := spinnerIDPool[rand.Intn(len(spinnerIDPool))]
	loading := &LoadingModal{TextView: tview.NewTextView(), spinnerCharSet: spinner.CharSets[spinnerSetID]}
	loading.SetDynamicColors(true).SetRegions(true)
	loading.SetTitle(" Connecting... ")
	loading.SetText("Forwarding you to your remote machine.").
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderPadding(2, 2, 2, 2)

	app.SetRoot(tview.NewPages().
		AddPage("loading", generateModal(loading, 40, 10), true, true).
		AddPage("term", generateModal(txt, 50, 10), true, false),
		true)

	stopAnimation := make(chan struct{})
	defer close(stopAnimation)
	go Animate(stopAnimation, app)

	return app.Run()
}

// generateModal makes a centered object
func generateModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
