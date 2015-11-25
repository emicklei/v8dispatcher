package v8dispatcher

// Module represents a Javascript object with functions that call into its Go counterpart.
type Module interface {
	// ModuleDefinition returns the name of the module as it will be known in Javascript
	// and Javascript source to create this module (global variable).
	// It returns an error if loading the source failed.
	ModuleDefinition() (string, string)

	// Perform will call the function associated to the Method of the message.
	Perform(msg MessageSend) (interface{}, error)
}
