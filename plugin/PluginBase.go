package plugin

type PluginBase interface {
	// identification
	Name() string
	GetPrefix() string
	Get() interface{}

	// lifecycle
	InitFlags()
	Configure() error
	Run() error
	Stop() <-chan bool
}
