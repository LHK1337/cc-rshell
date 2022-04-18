package app

import (
	"cc-rshell-server/model"
	"cc-rshell-server/sockets/types"
	"context"
	"github.com/briandowns/spinner"
	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
	"log"
	"math/rand"
)

var spinnerIDPool = []int{
	1, 14, 24, 27, 32, 43, 47, 50, 54, 57, 69, 78, 79, 80,
}

func RunApp(screen tcell.Screen, d types.ComputerDescriptor, registry types.ClientRegistry, name string, pubKey ssh.PublicKey) error {
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[!] SSH session connected to %d panicked. Reason: %v\n", d.ComputerID(), r)
		}
	}()

	isAnon := name == ""
	if isAnon {
		name = "anonymous"
	}

	app := tview.NewApplication().SetScreen(screen)

	spinnerSetID := spinnerIDPool[rand.Intn(len(spinnerIDPool))]
	loading := &LoadingModal{TextView: tview.NewTextView(), spinnerCharSet: spinner.CharSets[spinnerSetID]}
	loading.SetDynamicColors(true).SetRegions(true)
	loading.SetTitle(" Connecting... ")
	loading.SetText("Forwarding you to your remote machine.").
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderPadding(2, 2, 2, 2)

	procID := 0

	fbChannel := make(chan *model.FrameBuffer)
	registry[d.ComputerID()].RegisterFramebufferChannel(procID, fbChannel)

	animateContext, cancelAnimation := context.WithCancel(appContext)
	go Animate(animateContext, app)

	pages := tview.NewPages()

	fbv := NewFramebufferView(appContext, fbChannel, d.Colors, func() {
		pages.HidePage("loading")
		pages.ShowPage("term")
		cancelAnimation()
	}, func() {
		app.Draw()
	},
		func() {
			app.Stop()
		},
	)

	pages.
		AddPage("loading", generateModal(loading, 40, 10), true, true).
		AddPage("term", fbv, true, false)
	app.SetRoot(pages, true)

	go fbv.Worker()

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
