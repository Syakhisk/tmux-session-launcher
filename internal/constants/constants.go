package constants

const SockAddress = "/tmp/tmux-session-launcher.sock"

// JSON-RPC method names following RPC conventions
const (
	MethodModeNext       = "mode.next"
	MethodModePrev       = "mode.previous"
	MethodModeGet        = "mode.get"
	MethodContentGet     = "content.get"
	MethodLauncherOpenIn = "launcher.openIn"
)
