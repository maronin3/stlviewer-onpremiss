package webview

import (
	"log"
	"stlviewer/utils"

	"github.com/webview/webview"
)

func Webview() {
	// Debug
	w := webview.New(utils.Env.DevMode)
	defer w.Destroy()

	// Setting
	w.SetTitle(utils.Env.Title)
	w.SetSize(utils.Env.Width, utils.Env.Height, webview.HintNone)

	// Binding
	w.Bind("quit", func() {
		w.Terminate()
	})

	w.Bind("test", func() {
		log.Println("test")
	})

	w.Navigate(utils.Env.URL)
	w.Run()
}
