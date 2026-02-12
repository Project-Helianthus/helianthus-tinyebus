package emulation

import "fmt"

type ResponseEnvelope struct {
	MinDelayMillis uint32
	MaxDelayMillis uint32
}

func (e ResponseEnvelope) validate() error {
	if e.MaxDelayMillis > 0 && e.MinDelayMillis > e.MaxDelayMillis {
		return fmt.Errorf("envelope min delay exceeds max delay: %w", ErrInvalidConfiguration)
	}
	return nil
}

func ValidateResponseEnvelope(responses []EmulatedResponse, envelope ResponseEnvelope) error {
	if err := envelope.validate(); err != nil {
		return err
	}
	for idx, response := range responses {
		delayMillis := response.RespondAtMillis - response.RequestedAtMillis
		if delayMillis < envelope.MinDelayMillis {
			return fmt.Errorf(
				"response[%d] delay %dms below envelope min %dms: %w",
				idx,
				delayMillis,
				envelope.MinDelayMillis,
				ErrTimingConstraint,
			)
		}
		if envelope.MaxDelayMillis > 0 && delayMillis > envelope.MaxDelayMillis {
			return fmt.Errorf(
				"response[%d] delay %dms above envelope max %dms: %w",
				idx,
				delayMillis,
				envelope.MaxDelayMillis,
				ErrTimingConstraint,
			)
		}
	}
	return nil
}

type QueryStep struct {
	AdvanceMillis uint32
	Frame         Frame
}

type Harness struct {
	target  *Target
	now     uint32
	history []EmulatedResponse
}

func NewHarness(target *Target) *Harness {
	return &Harness{
		target: target,
	}
}

func (h *Harness) NowMillis() uint32 {
	if h == nil {
		return 0
	}
	return h.now
}

func (h *Harness) AdvanceMillis(delta uint32) {
	if h == nil || delta == 0 {
		return
	}
	h.now += delta
}

func (h *Harness) Query(frame Frame) (EmulatedResponse, error) {
	if h == nil || h.target == nil {
		return EmulatedResponse{}, fmt.Errorf("missing harness target: %w", ErrInvalidConfiguration)
	}
	response, err := h.target.Emulate(RequestEvent{
		AtMillis: h.now,
		Frame:    frame,
	})
	if err != nil {
		return EmulatedResponse{}, err
	}
	h.history = append(h.history, response)
	return response, nil
}

func (h *Harness) RunSequence(steps []QueryStep) ([]EmulatedResponse, error) {
	if h == nil || h.target == nil {
		return nil, fmt.Errorf("missing harness target: %w", ErrInvalidConfiguration)
	}
	responses := make([]EmulatedResponse, 0, len(steps))
	for idx, step := range steps {
		h.AdvanceMillis(step.AdvanceMillis)
		response, err := h.Query(step.Frame)
		if err != nil {
			return nil, fmt.Errorf("step[%d]: %w", idx, err)
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (h *Harness) History() []EmulatedResponse {
	if h == nil || len(h.history) == 0 {
		return nil
	}
	out := make([]EmulatedResponse, len(h.history))
	copy(out, h.history)
	return out
}
