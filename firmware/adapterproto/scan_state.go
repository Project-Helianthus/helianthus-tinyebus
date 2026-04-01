package adapterproto

import "fmt"

type ProtocolStateSnapshot struct {
	State            byte   `json:"state"`
	StateName        string `json:"state_name"`
	Flags            byte   `json:"flags"`
	FlagsName        string `json:"flags_name"`
	Tick             uint32 `json:"tick"`
	Deadline         uint32 `json:"deadline"`
	WindowDelay      uint32 `json:"window_delay"`
	WindowLimit      uint32 `json:"window_limit"`
	ScanSeed         uint32 `json:"scan_seed"`
	MergedWindow     uint32 `json:"merged_window"`
	ActiveScanSlot   byte   `json:"active_scan_slot"`
	DescriptorCursor uint16 `json:"descriptor_cursor"`
}

type ScanDeadlineExample struct {
	Now            uint32 `json:"now"`
	RequestedDelay uint32 `json:"requested_delay"`
	EffectiveDelay uint32 `json:"effective_delay"`
	Deadline       uint32 `json:"deadline"`
}

type ProtocolStateDispatchSample struct {
	Trigger   string                `json:"trigger"`
	InputCode byte                  `json:"input_code,omitempty"`
	Return    byte                  `json:"return"`
	Before    ProtocolStateSnapshot `json:"before"`
	After     ProtocolStateSnapshot `json:"after"`
	Notes     string                `json:"notes"`
}

type StatusFrameSample struct {
	Builder    string   `json:"builder"`
	PrefixHex  []string `json:"prefix_hex"`
	SummaryHex []string `json:"summary_hex"`
	Length     int      `json:"length"`
	Notes      string   `json:"notes"`
}

type OracleScanReport struct {
	Initial          ProtocolStateSnapshot         `json:"initial"`
	DeadlineExamples []ScanDeadlineExample         `json:"deadline_examples"`
	DispatchSamples  []ProtocolStateDispatchSample `json:"dispatch_samples"`
	StatusSnapshot   StatusFrameSample             `json:"status_snapshot"`
	StatusVariant    StatusFrameSample             `json:"status_variant"`
}

type ProtocolScanModel struct {
	state            byte
	flags            byte
	tick             uint32
	deadline         uint32
	windowDelay      uint32
	windowLimit      uint32
	scanSeed         uint32
	mergedWindow     uint32
	activeScanSlot   byte
	descriptorCursor uint16
	statusSeedLatch  byte
	dispatchCursor   byte
}

const (
	scanMinimumDelay  = uint32(0x3c)
	scanMinimumLimit  = uint32(0x000000f0)
	statusFrameMax    = 24
	statusFramePrefix = 6
)

func NewProtocolScanModel() *ProtocolScanModel {
	return &ProtocolScanModel{
		state:            3,
		flags:            0,
		tick:             0x00000140,
		deadline:         0x000001a4,
		windowDelay:      0x0000003c,
		windowLimit:      0x00000156,
		scanSeed:         0x000002a4,
		mergedWindow:     0x000001a8,
		activeScanSlot:   0x06,
		descriptorCursor: 0x00e9,
		statusSeedLatch:  0xa4,
	}
}

func (model *ProtocolScanModel) Clone() *ProtocolScanModel {
	if model == nil {
		return NewProtocolScanModel()
	}
	clone := *model
	return &clone
}

func (model *ProtocolScanModel) Snapshot() ProtocolStateSnapshot {
	if model == nil {
		return ProtocolStateSnapshot{}
	}
	return ProtocolStateSnapshot{
		State:            model.state,
		StateName:        protocolStateName(model.state),
		Flags:            model.flags,
		FlagsName:        protocolStateFlagsName(model.flags),
		Tick:             model.tick,
		Deadline:         model.deadline,
		WindowDelay:      model.windowDelay,
		WindowLimit:      model.windowLimit,
		ScanSeed:         model.scanSeed,
		MergedWindow:     model.mergedWindow,
		ActiveScanSlot:   model.activeScanSlot,
		DescriptorCursor: model.descriptorCursor,
	}
}

func (model *ProtocolScanModel) DispatchScanCode(code byte) bool {
	if model == nil {
		return false
	}

	transformed := seedWindowTransform(model.scanSeed)
	limit := normalizeScanLimit(transformed)
	delay := normalizeScanDelay(transformed)
	if delay > limit {
		delay = limit
	}

	switch code {
	case 0x01:
		model.descriptorCursor = 0x0264
		model.activeScanSlot = 0x01
	case 0x03:
		model.descriptorCursor = 0x0260
		model.activeScanSlot = 0x03
	case 0x33:
		model.windowLimit = limit
	case 0x35:
		model.statusSeedLatch = byte(model.scanSeed)
	case 0x36:
		model.descriptorCursor = 0x0268
		model.activeScanSlot = 0x06
	case 0x3a:
		model.windowDelay = delay
		model.deadline = model.tick + delay
	case 0x3b:
		merged := transformed
		if merged < delay {
			merged = delay
		}
		if merged > limit {
			merged = limit
		}
		model.mergedWindow = merged
	default:
		return false
	}

	model.flags = 3
	return true
}

func (model *ProtocolScanModel) StartScanWindow() {
	if model == nil {
		return
	}

	model.dispatchCursor = 0
	model.state = 7
	model.flags = 3
	model.activeScanSlot = 0x06
	model.initializeScanSlotFull(0x06)
	model.statusSeedLatch = byte(model.scanSeed)
}

func (model *ProtocolScanModel) ContinueScanWindow() {
	dispatchCodes := []byte{0x01, 0x33, 0x35, 0x36, 0x3a, 0x3b, 0x03}

	if model == nil {
		return
	}

	if model.state == 0 || model.state == 3 {
		model.StartScanWindow()
		return
	}

	code := dispatchCodes[int(model.dispatchCursor)%len(dispatchCodes)]
	if !model.DispatchScanCode(code) {
		return
	}

	model.initializeScanSlotFull(model.activeScanSlot)

	model.state = 5
	model.dispatchCursor = (model.dispatchCursor + 1) % byte(len(dispatchCodes))
}

func (model *ProtocolScanModel) initializeScanSlotFull(slotID byte) {
	addrLo := byte(slotID*0x20 + 8)
	addrHi := byte((slotID >> 3) & 0x1F)
	if 0xF7 < byte(slotID*0x20) {
		addrHi++
	}
	addrHi |= 0x24
	model.descriptorCursor = uint16(addrHi)<<8 | uint16(addrLo)

	// recompute_scan_masks_tail equivalent: advance cursor by 0x2C
	model.descriptorCursor += 0x2C
	model.statusSeedLatch = byte(model.scanSeed)
}

func DefaultOracleScanReport() OracleScanReport {
	dispatchCodes := []byte{0x01, 0x03, 0x33, 0x35, 0x36, 0x3a, 0x3b}
	dispatchNotes := []string{
		"select descriptor window 0x0264",
		"select descriptor window 0x0260",
		"derive scan window limit from scan seed",
		"latch scan seed low byte into status path",
		"select descriptor window 0x0268",
		"derive scan window delay from scan seed",
		"merge scan window candidate from scan seed",
	}

	model := NewProtocolScanModel()
	initial := model.Snapshot()
	deadlineExamples := []ScanDeadlineExample{
		deadlineExample(0x00000140, 0x00),
		deadlineExample(0x00000140, 0x10),
		deadlineExample(0xfffffff0, 0x20),
		deadlineExample(0x00000140, 0x80),
	}

	dispatchSamples := make([]ProtocolStateDispatchSample, 0, len(dispatchCodes))
	for idx, code := range dispatchCodes {
		dispatchModel := NewProtocolScanModel()
		before := dispatchModel.Snapshot()
		ok := dispatchModel.DispatchScanCode(code)
		dispatchSamples = append(dispatchSamples, ProtocolStateDispatchSample{
			Trigger:   "command_id_dispatch",
			InputCode: code,
			Return:    boolToByte(ok),
			Before:    before,
			After:     dispatchModel.Snapshot(),
			Notes:     dispatchNotes[idx],
		})
	}

	snapshotModel := NewProtocolScanModel()
	snapshotModel.StartScanWindow()
	snapshotFrame := buildStatusFrameSample("build_status_snapshot_frame", 0x35, snapshotModel, "runtime wire-level status snapshot sample")

	variantModel := NewProtocolScanModel()
	variantModel.StartScanWindow()
	variantModel.ContinueScanWindow()
	variantFrame := buildStatusFrameSample("build_status_variant_frame", 0x37, variantModel, "runtime wire-level status variant sample after first scan dispatch")

	return OracleScanReport{
		Initial:          initial,
		DeadlineExamples: deadlineExamples,
		DispatchSamples:  dispatchSamples,
		StatusSnapshot:   snapshotFrame,
		StatusVariant:    variantFrame,
	}
}

func deadlineExample(now, requestedDelay uint32) ScanDeadlineExample {
	effectiveDelay := normalizeScanDelay(requestedDelay)
	return ScanDeadlineExample{
		Now:            now,
		RequestedDelay: requestedDelay,
		EffectiveDelay: effectiveDelay,
		Deadline:       deadlineAfter(now, requestedDelay),
	}
}

func normalizeScanDelay(requestedDelay uint32) uint32 {
	if requestedDelay == 0 || requestedDelay < scanMinimumDelay {
		return scanMinimumDelay
	}
	return requestedDelay
}

func normalizeScanLimit(requestedLimit uint32) uint32 {
	if requestedLimit < scanMinimumLimit {
		return scanMinimumLimit
	}
	return requestedLimit
}

func deadlineAfter(now, requestedDelay uint32) uint32 {
	return now + normalizeScanDelay(requestedDelay)
}

func deadlineReached(now, deadline uint32) bool {
	return int32(now-deadline) >= 0
}

func protocolStateName(state byte) string {
	switch state {
	case 0:
		return "idle"
	case 1:
		return "pending"
	case 3:
		return "ready"
	case 5:
		return "variant"
	case 7:
		return "scan"
	case 8:
		return "scan_holding"
	default:
		return fmt.Sprintf("unknown(0x%02x)", state)
	}
}

func protocolStateFlagsName(flags byte) string {
	switch flags {
	case 0:
		return "idle"
	case 1:
		return "pending_transition"
	case 3:
		return "ready_and_pending"
	default:
		return fmt.Sprintf("unknown(0x%02x)", flags)
	}
}

func buildStatusFrameSample(builder string, kind byte, model *ProtocolScanModel, notes string) StatusFrameSample {
	frame := buildStatusFrame(kind, model)
	return StatusFrameSample{
		Builder:    builder,
		PrefixHex:  hexStrings(frame[:statusFramePrefix]),
		SummaryHex: hexStrings(frame[statusFramePrefix:]),
		Length:     len(frame),
		Notes:      notes,
	}
}

func buildStatusFrame(kind byte, model *ProtocolScanModel) []byte {
	frame := make([]byte, 0, statusFrameMax)
	frame = append(frame, 0x63, 0x82, 0x53, 0x63, kind, 0x01)
	frame = append(frame, model.state, 0x3d, 0x07, 0x01)

	cached3 := byte(model.descriptorCursor >> 8)
	cached4 := byte(model.windowDelay)
	cached5 := byte(model.windowLimit)
	frame = append(frame,
		model.flags,
		model.activeScanSlot,
		byte(model.descriptorCursor),
		cached3,
		cached4,
		cached5,
	)
	frame = append(frame, 0x0c, model.statusSeedLatch+0x06)
	frame = append(frame,
		hexDigit(cached3>>4),
		hexDigit(cached3),
		hexDigit(cached4>>4),
		hexDigit(cached4),
		hexDigit(cached5>>4),
		hexDigit(cached5),
	)
	return frame
}

func seedWindowTransform(seed uint32) uint32 {
	r8 := seed >> 8
	l8 := seed << 8
	l24 := seed << 24
	return r8 | l8 | l24
}

func boolToByte(value bool) byte {
	if value {
		return 1
	}
	return 0
}

func hexDigit(value byte) byte {
	value &= 0x0f
	if value < 10 {
		return '0' + value
	}
	return 'a' + (value - 10)
}
