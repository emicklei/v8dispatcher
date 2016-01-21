package v8dispatcher

// Module represents a Javascript object with functions that call into its Go counterpart.
type Module interface {
	// ModuleDefinition returns the name of the module as it will be known in Javascript
	// and Javascript source to create this module (global variable).
	// It returns an error if loading the source failed.
	Definition() (name string, source string, err error)

	// Perform will call the function associated with the Method of the message.
	// The returnValue must be JSON marshallable
	// It returns an error if performing the function failed.
	Perform(msg MessageSend) (returnValue interface{}, err error)
}
