package adapterproto

import "testing"

func TestRuntimeContractModelTransitions(t *testing.T) {
	t.Parallel()

	model := NewRuntimeContractModel()

	init := model.Init(0x01)
	if init.Outcome != runtimeOutcomeResetted {
		t.Fatalf("Init outcome = %q; want %q", init.Outcome, runtimeOutcomeResetted)
	}
	if init.Response == nil || init.Response.CommandName != "resetted" {
		t.Fatalf("Init response = %+v; want resetted", init.Response)
	}
	if init.StateAfter.Active {
		t.Fatalf("Init state after = %+v; want idle", init.StateAfter)
	}

	start := model.Start(0x31)
	if start.Outcome != runtimeOutcomeStarted {
		t.Fatalf("Start outcome = %q; want %q", start.Outcome, runtimeOutcomeStarted)
	}
	if start.Response == nil || start.Response.CommandName != "started" {
		t.Fatalf("Start response = %+v; want started", start.Response)
	}
	if !start.StateAfter.Active || start.StateAfter.Initiator != 0x31 {
		t.Fatalf("Start state after = %+v; want active initiator 0x31", start.StateAfter)
	}

	sendAfterStart := model.Send(0x55)
	if sendAfterStart.Outcome != runtimeOutcomeReceived {
		t.Fatalf("Send outcome = %q; want %q", sendAfterStart.Outcome, runtimeOutcomeReceived)
	}
	if sendAfterStart.Response == nil || sendAfterStart.Response.CommandName != "received" {
		t.Fatalf("Send response = %+v; want received", sendAfterStart.Response)
	}
	if len(sendAfterStart.RequestHex) != 2 || sendAfterStart.RequestHex[0] != "c5" || sendAfterStart.RequestHex[1] != "95" {
		t.Fatalf("Send request hex = %#v; want [c5 95]", sendAfterStart.RequestHex)
	}
	if len(sendAfterStart.ResponseHex) != 1 || sendAfterStart.ResponseHex[0] != "55" {
		t.Fatalf("Send response hex = %#v; want [55]", sendAfterStart.ResponseHex)
	}

	sendCancel := model.Send(0xAA)
	if sendCancel.Outcome != runtimeOutcomeReceived {
		t.Fatalf("Send cancel outcome = %q; want %q", sendCancel.Outcome, runtimeOutcomeReceived)
	}
	if sendCancel.Response == nil || sendCancel.Response.CommandName != "received" {
		t.Fatalf("Send cancel response = %+v; want received", sendCancel.Response)
	}
	if sendCancel.StateAfter.Active {
		t.Fatalf("Send cancel state after = %+v; want idle", sendCancel.StateAfter)
	}

	cancel := model.Start(0xAA)
	if cancel.Outcome != runtimeOutcomeStartCancelled {
		t.Fatalf("Cancel outcome = %q; want %q", cancel.Outcome, runtimeOutcomeStartCancelled)
	}
	if cancel.Response != nil || cancel.ResponseHex != nil {
		t.Fatalf("Cancel response = %+v/%#v; want none", cancel.Response, cancel.ResponseHex)
	}
	if cancel.StateAfter.Active {
		t.Fatalf("Cancel state after = %+v; want idle", cancel.StateAfter)
	}

	sendWithoutStart := NewRuntimeContractModel().Send(0x55)
	if sendWithoutStart.Outcome != runtimeOutcomeErrorHost {
		t.Fatalf("Send without start outcome = %q; want %q", sendWithoutStart.Outcome, runtimeOutcomeErrorHost)
	}
	if sendWithoutStart.Response == nil || sendWithoutStart.Response.CommandName != "error_host" {
		t.Fatalf("Send without start response = %+v; want error_host", sendWithoutStart.Response)
	}
	if len(sendWithoutStart.RequestHex) != 2 || sendWithoutStart.RequestHex[0] != "c5" || sendWithoutStart.RequestHex[1] != "95" {
		t.Fatalf("Send without start request hex = %#v; want [c5 95]", sendWithoutStart.RequestHex)
	}
	if len(sendWithoutStart.ResponseHex) != 2 {
		t.Fatalf("Send without start response hex = %#v; want two bytes", sendWithoutStart.ResponseHex)
	}
	if sendWithoutStart.Response.Data != runtimeErrorNoActiveStart {
		t.Fatalf("Send without start response data = 0x%02x; want 0x%02x", sendWithoutStart.Response.Data, runtimeErrorNoActiveStart)
	}
}

func TestDefaultOracleReportIncludesRuntimeSamples(t *testing.T) {
	t.Parallel()

	report := DefaultOracleReport()

	if report.Runtime.Init.Outcome != runtimeOutcomeResetted {
		t.Fatalf("runtime.init.outcome = %q; want %q", report.Runtime.Init.Outcome, runtimeOutcomeResetted)
	}
	if report.Runtime.StartSuccess.Outcome != runtimeOutcomeStarted {
		t.Fatalf("runtime.start_success.outcome = %q; want %q", report.Runtime.StartSuccess.Outcome, runtimeOutcomeStarted)
	}
	if report.Runtime.SendAfterStart.Response == nil || report.Runtime.SendAfterStart.Response.CommandName != "received" {
		t.Fatalf("runtime.send_after_start.response = %+v; want received", report.Runtime.SendAfterStart.Response)
	}
	if report.Runtime.SendCancel.Response == nil || report.Runtime.SendCancel.Response.CommandName != "received" {
		t.Fatalf("runtime.send_cancel.response = %+v; want received", report.Runtime.SendCancel.Response)
	}
	if report.Runtime.SendWithoutStart.Response == nil || report.Runtime.SendWithoutStart.Response.CommandName != "error_host" {
		t.Fatalf("runtime.send_without_start.response = %+v; want error_host", report.Runtime.SendWithoutStart.Response)
	}
	if report.Runtime.StartCancel.Response != nil {
		t.Fatalf("runtime.start_cancel.response = %+v; want nil", report.Runtime.StartCancel.Response)
	}
}
