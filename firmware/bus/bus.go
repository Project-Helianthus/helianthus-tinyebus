package bus

// Transport defines the byte-level bus link expected by higher-level bus logic.
// Implementations are provided in later milestones.
type Transport interface {
	ReadByte() (byte, error)
	WriteByte(value byte) error
}

// Engine describes the minimal bus lifecycle hooks expected by firmware wiring.
// This is a contract placeholder and intentionally has no implementation here.
type Engine interface {
	Init(link Transport) error
	Poll() error
}
