package emulation

import (
	"errors"
	"fmt"
)

var (
	ErrNoMatchingRule        = errors.New("target emulation no matching rule")
	ErrTimingConstraint      = errors.New("target emulation timing constraint")
	ErrInvalidConfiguration  = errors.New("target emulation invalid configuration")
	ErrRequestTargetMismatch = errors.New("target emulation request target mismatch")
)

type Frame struct {
	Source    byte
	Target    byte
	Primary   byte
	Secondary byte
	Data      []byte
}

type RequestMatcher interface {
	Match(frame Frame) bool
}

type MatchFunc func(frame Frame) bool

func (fn MatchFunc) Match(frame Frame) bool {
	if fn == nil {
		return false
	}
	return fn(frame)
}

type ResponseBuilder interface {
	Build(frame Frame) (ResponsePlan, error)
}

type BuildFunc func(frame Frame) (ResponsePlan, error)

func (fn BuildFunc) Build(frame Frame) (ResponsePlan, error) {
	if fn == nil {
		return ResponsePlan{}, fmt.Errorf("missing response builder: %w", ErrInvalidConfiguration)
	}
	return fn(frame)
}

type ResponsePlan struct {
	DelayMillis uint32
	Data        []byte
}

type TimingConstraints struct {
	MinResponseDelayMillis uint32
	MaxResponseDelayMillis uint32
}

func (c TimingConstraints) validate(delayMillis uint32) error {
	if c.MaxResponseDelayMillis > 0 && c.MinResponseDelayMillis > c.MaxResponseDelayMillis {
		return fmt.Errorf("min delay exceeds max delay: %w", ErrInvalidConfiguration)
	}
	if delayMillis < c.MinResponseDelayMillis {
		return fmt.Errorf(
			"response delay %dms below min %dms: %w",
			delayMillis,
			c.MinResponseDelayMillis,
			ErrTimingConstraint,
		)
	}
	if c.MaxResponseDelayMillis > 0 && delayMillis > c.MaxResponseDelayMillis {
		return fmt.Errorf(
			"response delay %dms above max %dms: %w",
			delayMillis,
			c.MaxResponseDelayMillis,
			ErrTimingConstraint,
		)
	}
	return nil
}

func (c TimingConstraints) active() bool {
	return c.MinResponseDelayMillis != 0 || c.MaxResponseDelayMillis != 0
}

type Rule struct {
	Name    string
	Matcher RequestMatcher
	Builder ResponseBuilder
	Timing  TimingConstraints
}

type RequestEvent struct {
	AtMillis uint32
	Frame    Frame
}

type EmulatedResponse struct {
	Rule              string
	RequestedAtMillis uint32
	RespondAtMillis   uint32
	Frame             Frame
}

type Target struct {
	Name          string
	Address       byte
	DefaultTiming TimingConstraints
	Rules         []Rule
}

func (t *Target) Emulate(event RequestEvent) (EmulatedResponse, error) {
	if t == nil {
		return EmulatedResponse{}, fmt.Errorf("missing target: %w", ErrInvalidConfiguration)
	}
	if event.Frame.Target != t.Address {
		return EmulatedResponse{}, fmt.Errorf(
			"request target 0x%02x does not match emulated target 0x%02x: %w",
			event.Frame.Target,
			t.Address,
			ErrRequestTargetMismatch,
		)
	}

	for _, rule := range t.Rules {
		if rule.Matcher == nil {
			return EmulatedResponse{}, fmt.Errorf("rule %q missing matcher: %w", rule.Name, ErrInvalidConfiguration)
		}
		if !rule.Matcher.Match(event.Frame) {
			continue
		}
		if rule.Builder == nil {
			return EmulatedResponse{}, fmt.Errorf("rule %q missing builder: %w", rule.Name, ErrInvalidConfiguration)
		}
		plan, err := rule.Builder.Build(event.Frame)
		if err != nil {
			return EmulatedResponse{}, err
		}

		timing := rule.Timing
		if !timing.active() {
			timing = t.DefaultTiming
		}
		if err := timing.validate(plan.DelayMillis); err != nil {
			return EmulatedResponse{}, err
		}

		return EmulatedResponse{
			Rule:              rule.Name,
			RequestedAtMillis: event.AtMillis,
			RespondAtMillis:   event.AtMillis + plan.DelayMillis,
			Frame: Frame{
				Source:    t.Address,
				Target:    event.Frame.Source,
				Primary:   event.Frame.Primary,
				Secondary: event.Frame.Secondary,
				Data:      append([]byte(nil), plan.Data...),
			},
		}, nil
	}

	return EmulatedResponse{}, fmt.Errorf(
		"request pb=0x%02x sb=0x%02x: %w",
		event.Frame.Primary,
		event.Frame.Secondary,
		ErrNoMatchingRule,
	)
}

func MatchPrimarySecondary(primary, secondary byte) MatchFunc {
	return func(frame Frame) bool {
		return frame.Primary == primary && frame.Secondary == secondary
	}
}

func MatchPrimarySecondaryWithPrefix(primary, secondary byte, prefix []byte) MatchFunc {
	copiedPrefix := append([]byte(nil), prefix...)
	return func(frame Frame) bool {
		if frame.Primary != primary || frame.Secondary != secondary {
			return false
		}
		if len(copiedPrefix) == 0 {
			return true
		}
		if len(frame.Data) < len(copiedPrefix) {
			return false
		}
		for idx := range copiedPrefix {
			if frame.Data[idx] != copiedPrefix[idx] {
				return false
			}
		}
		return true
	}
}
