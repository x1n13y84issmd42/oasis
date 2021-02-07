package contract

// IEntityTrait is an interface of entity traits.
type IEntityTrait interface {
	GetLogger() Logger
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

// GetLogger returns logger.
func (entity EntityTrait) GetLogger() Logger {
	return entity.Log
}
