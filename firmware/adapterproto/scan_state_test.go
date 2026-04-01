package adapterproto

import "testing"

func TestNormalizeScanDelay(t *testing.T) {
	t.Parallel()

	if got := normalizeScanDelay(0); got != scanMinimumDelay {
		t.Fatalf("normalizeScanDelay(0) = 0x%02x; want 0x%02x", got, scanMinimumDelay)
	}
	if got := normalizeScanDelay(0x10); got != scanMinimumDelay {
		t.Fatalf("normalizeScanDelay(0x10) = 0x%02x; want 0x%02x", got, scanMinimumDelay)
	}
	if got := normalizeScanDelay(0x3c); got != scanMinimumDelay {
		t.Fatalf("normalizeScanDelay(0x3c) = 0x%02x; want 0x%02x", got, scanMinimumDelay)
	}
	if got := normalizeScanDelay(0x80); got != 0x80 {
		t.Fatalf("normalizeScanDelay(0x80) = 0x%02x; want 0x80", got)
	}
}

func TestDeadlineReachedIsWrapSafe(t *testing.T) {
	t.Parallel()

	if !deadlineReached(0x00000010, 0x00000010) {
		t.Fatalf("deadlineReached(equal) = false; want true")
	}
	if deadlineReached(0x0000000f, 0x00000010) {
		t.Fatalf("deadlineReached(before) = true; want false")
	}
	if !deadlineReached(0x00000020, 0xfffffff0) {
		t.Fatalf("deadlineReached(wrapped) = false; want true")
	}
	if deadlineReached(0xffffffe0, 0x00000010) {
		t.Fatalf("deadlineReached(future wrapped) = true; want false")
	}
}

func TestDefaultOracleScanReportBuildsDeterministicFrames(t *testing.T) {
	t.Parallel()

	report := DefaultOracleScanReport()

	if report.StatusSnapshot.Length != len(report.StatusSnapshot.PrefixHex)+len(report.StatusSnapshot.SummaryHex) {
		t.Fatalf("snapshot length = %d; want prefix+summary", report.StatusSnapshot.Length)
	}
	if report.StatusVariant.Length != len(report.StatusVariant.PrefixHex)+len(report.StatusVariant.SummaryHex) {
		t.Fatalf("variant length = %d; want prefix+summary", report.StatusVariant.Length)
	}
	if report.StatusSnapshot.Length != statusFrameMax || report.StatusVariant.Length != statusFrameMax {
		t.Fatalf("status frame lengths = %d/%d; want %d", report.StatusSnapshot.Length, report.StatusVariant.Length, statusFrameMax)
	}
	if report.StatusSnapshot.PrefixHex[4] != "35" {
		t.Fatalf("snapshot prefix kind = %#v; want 0x35", report.StatusSnapshot.PrefixHex)
	}
	if report.StatusVariant.PrefixHex[4] != "37" {
		t.Fatalf("variant prefix kind = %#v; want 0x37", report.StatusVariant.PrefixHex)
	}
	if got := report.StatusSnapshot.SummaryHex[len(report.StatusSnapshot.SummaryHex)-6:]; got[0] != "32" || got[1] != "34" {
		t.Fatalf("snapshot trailer = %#v; want 24....", got)
	}
	if got := report.StatusVariant.SummaryHex[len(report.StatusVariant.SummaryHex)-6:]; got[0] != "32" || got[1] != "34" {
		t.Fatalf("variant trailer = %#v; want 24....", got)
	}
}

func TestScanWindowBuildersMirrorRuntimeOrdering(t *testing.T) {
	t.Parallel()

	model := NewProtocolScanModel()
	model.StartScanWindow()
	if snapshot := model.Snapshot(); snapshot.State != 7 || snapshot.ActiveScanSlot != 0x06 || snapshot.DescriptorCursor != 0x24F4 {
		t.Fatalf("StartScanWindow() snapshot = %+v; want state=7 slot=0x06 cursor=0x24F4", snapshot)
	}

	model.ContinueScanWindow()
	if snapshot := model.Snapshot(); snapshot.State != 5 || snapshot.ActiveScanSlot != 0x01 || snapshot.DescriptorCursor != 0x2454 {
		t.Fatalf("ContinueScanWindow() snapshot = %+v; want state=5 slot=0x01 cursor=0x2454", snapshot)
	}
}

func TestDispatchScanCodeSamples(t *testing.T) {
	t.Parallel()

	report := DefaultOracleScanReport()

	if len(report.DispatchSamples) != 7 {
		t.Fatalf("dispatch sample len = %d; want 7", len(report.DispatchSamples))
	}
	if report.DispatchSamples[0].After.DescriptorCursor != 0x0264 {
		t.Fatalf("dispatch[0].after.descriptor_cursor = 0x%04x; want 0x0264", report.DispatchSamples[0].After.DescriptorCursor)
	}
	if report.DispatchSamples[2].After.WindowLimit != 0xa402a402 {
		t.Fatalf("dispatch[2].after.window_limit = 0x%08x; want 0xa402a402", report.DispatchSamples[2].After.WindowLimit)
	}
	if report.DispatchSamples[5].After.WindowDelay != 0xa402a402 {
		t.Fatalf("dispatch[5].after.window_delay = 0x%08x; want 0xa402a402", report.DispatchSamples[5].After.WindowDelay)
	}
	if report.DispatchSamples[6].After.MergedWindow != 0xa402a402 {
		t.Fatalf("dispatch[6].after.merged_window = 0x%08x; want 0xa402a402", report.DispatchSamples[6].After.MergedWindow)
	}
}
