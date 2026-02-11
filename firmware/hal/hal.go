package hal

// UART is the minimal serial interface required by the firmware bus layer.
// Concrete TinyGo/board bindings are introduced in later milestones.
type UART interface {
	ReadByte() (byte, error)
	WriteByte(value byte) error
}

// Clock abstracts timing access used by polling loops and timeouts.
type Clock interface {
	NowMillis() uint32
	SleepMillis(duration uint32)
}

// Board groups hardware capabilities needed by the firmware bootstrap path.
type Board interface {
	UART() UART
	Clock() Clock
}
