module testClient

go 1.14

replace git.vnnox.net/ncp/xframe/functional/unet => ../../xframe/functional/unet

replace git.vnnox.net/ncp/xframe/functional/upackage => ../../xframe/functional/upackage

replace tools => ../common/tools

require (
	git.vnnox.net/ncp/xframe/functional/unet v0.0.0-00010101000000-000000000000
	tools v0.0.0-00010101000000-000000000000
)
