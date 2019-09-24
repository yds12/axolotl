package main

import (
	"bufio"
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"sync"

	astilectron "github.com/asticode/go-astilectron"
	log "github.com/sirupsen/logrus"

	"github.com/nanu-c/textsecure-qml/app/config"
	"github.com/nanu-c/textsecure-qml/app/helpers"
	"github.com/nanu-c/textsecure-qml/app/push"
	"github.com/nanu-c/textsecure-qml/app/ui"
	"github.com/nanu-c/textsecure-qml/app/webserver"
	"github.com/nanu-c/textsecure-qml/app/worker"
)

var e string

func init() {
	flag.StringVar(&config.MainQml, "qml", "qml/phoneui/main.qml", "The qml file to load.")
	flag.StringVar(&e, "e", "", "Usage")
}
func print(stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func setup() {
	config.SetupConfig()
	helpers.SetupLogging()
	log.SetLevel(log.DebugLevel)
	log.Infof("Starting Signal for Ubuntu version %s", config.AppVersion)
}
func runBackend() {
	ui.SetEngine()
	//
	// ui.Engine.AddImageProvider("avatar", store.AvatarImageProvider)
	ui.InitModels()
	//
	// ui.Engine.Context().SetVar("textsecure", worker.Api)
	// ui.Engine.Context().SetVar("appVersion", config.AppVersion)
	// ui.Engine.Context().SetVar("cacheDir", config.CacheDir)
	ui.SetComponent()
	//
	// ui.Win.Show()
	go worker.RunBackend()
	if config.IsPushHelper {
		push.PushHelperProcess()
	}
}
func runUI() error {
	defer wg.Done()
	if e != "ut" && e != "me" && e != "lorca" {
		ui.RunUi(e)
		runElectron()
	} else {
		ui.RunUi(e)
	}
	os.Exit(0)
	return nil
}
func runElectron() {
	var a, _ = astilectron.New(astilectron.Options{
		AppName:            "axolotl",
		AppIconDefaultPath: "axolotl-web/public/axolotl.png", // If path is relative, it must be relative to the data directory
		AppIconDarwinPath:  "axolotl-web/public/axolotl.png", // Same here
		BaseDirectoryPath:  "dist",
	})
	defer a.Close()

	// Start astilectron
	a.Start()
	var w, _ = a.NewWindow("http://localhost:9080", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	})
	w.Create()
	// Blocking pattern
	a.Wait()
}
func runWebserver() {
	// Decrement the counter when the goroutine completes.
	defer wg.Done()
	log.Printf("Axolotl server started")

	// Fetch the URL.
	webserver.Run()
}

var wg sync.WaitGroup

func main() {

	setup()
	runBackend()
	log.Println("Setup completed")
	wg.Add(1)
	go runWebserver()
	wg.Add(1)
	go runUI()
	wg.Wait()
}
