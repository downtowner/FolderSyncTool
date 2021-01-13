module FolderSyncTool

go 1.14

replace git.vnnox.net/ncp/xframe/functional/unet => ../xframe/functional/unet

replace git.vnnox.net/ncp/xframe/functional/upackage => ../xframe/functional/upackage

require (
	git.vnnox.net/ncp/xframe/functional/unet v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.4.9
)
