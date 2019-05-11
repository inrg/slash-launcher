package main

import (
	"os"
	"time"

	"github.com/nokka/slash-launcher/bridge"
	"github.com/nokka/slash-launcher/d2"
	"github.com/nokka/slash-launcher/window"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/widgets"
)

func main() {
	// Enable high dpi scaling, useful for devices with high pixel density displays.
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	// Create the new QApplication.
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create the view and configure it.
	view := quick.NewQQuickView(nil)
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.SetFlags(core.Qt__FramelessWindowHint)

	// Set our main.qml to the view.
	view.SetSource(core.NewQUrl3("qml/main.qml", 0))

	// Create a new QML bridge that will bridge the client to Go.
	var qmlBridge = bridge.NewQmlBridge(nil)

	// Setup the bridge dependencies.
	qmlBridge.D2Launcher = d2.NewLauncher("/")
	qmlBridge.View = view

	// Connect the QML signals on the bridge to Go.
	qmlBridge.Connect()

	// Allows for windows to minimize on Darwin.
	window.AllowMinimize(view.WinId())

	// Set the bridge on the view.
	view.RootContext().SetContextProperty("QmlBridge", qmlBridge)

	// Center the view.
	view.SetPosition2(
		(widgets.QApplication_Desktop().Screen(0).Width()-view.Width())/2,
		(widgets.QApplication_Desktop().Screen(0).Height()-view.Height())/2,
	)

	// Make the view visible.
	view.Show()

	// Finally, execute the application.
	app.Exec()

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			qmlBridge.SetPatchProgress(0.1 * float32(i))
		}

	}()
}
