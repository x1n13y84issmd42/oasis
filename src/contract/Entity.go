package contract

// IEntityTrait is an interface of entity traits.
type IEntityTrait interface {
	SetLogger(log Logger)
}

// EntityTrait contains basic functions shared by multiple different app components,
// such as logging.
type EntityTrait struct {
	Log Logger
}

// Entity creates a new EntityTrait instance.
func Entity(log Logger) EntityTrait {
	return EntityTrait{
		Log: log,
	}
}

// SetLogger sets a logger instance to use.
func (e *EntityTrait) SetLogger(log Logger) {
	e.Log = log
}
