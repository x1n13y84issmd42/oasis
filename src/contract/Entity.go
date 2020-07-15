package contract

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
