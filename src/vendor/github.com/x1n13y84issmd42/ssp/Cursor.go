package ssp

// Cursor is a pair of argument list and an index of a current argument.
type Cursor struct {
	Args []string
	I    int
}

// Clone creates a new Cursor instance which contains only unparsed items from args.
func (args *Cursor) Clone() *Cursor {
	return &Cursor{
		Args: args.Tail(),
	}
}

// Tail return a slice of unparsed items from args.
func (args *Cursor) Tail() []string {
	if args.I >= len(args.Args) {
		return []string{}
	}

	return args.Args[args.I:]
}

// Get returns the current argument.
func (args *Cursor) Get() string {
	if args.I >= len(args.Args) {
		return ""
	}

	return args.Args[args.I]
}

// Len returns length of the internal slice.
func (args *Cursor) Len() int {
	return len(args.Args)
}

// Left returns the number of items left unparsed yet.
func (args *Cursor) Left() int {
	return len(args.Args) - args.I
}

// Shift shifts the internal pointer so it points to the next item to parse.
func (args *Cursor) Shift(v int) {
	args.I += v
}
