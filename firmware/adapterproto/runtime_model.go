package adapterproto

type RuntimeSessionSnapshot struct {
	Active    bool `json:"active"`
	Initiator byte `json:"initiator,omitempty"`
}

type RuntimeExchangeSample struct {
	RequestHex  []string               `json:"request_hex"`
	Request     ENHDecodedFrame        `json:"request"`
	ResponseHex []string               `json:"response_hex"`
	Response    *ENHDecodedFrame       `json:"response"`
	Outcome     string                 `json:"outcome"`
	StateBefore RuntimeSessionSnapshot `json:"state_before"`
	StateAfter  RuntimeSessionSnapshot `json:"state_after"`
}

type RuntimeContractReport struct {
	Init             RuntimeExchangeSample `json:"init"`
	StartSuccess     RuntimeExchangeSample `json:"start_success"`
	SendAfterStart   RuntimeExchangeSample `json:"send_after_start"`
	SendCancel       RuntimeExchangeSample `json:"send_cancel"`
	SendWithoutStart RuntimeExchangeSample `json:"send_without_start"`
	StartCancel      RuntimeExchangeSample `json:"start_cancel"`
}

type RuntimeContractModel struct {
	session RuntimeSessionSnapshot
}

const (
	runtimeOutcomeResetted       = "resetted"
	runtimeOutcomeStarted        = "started"
	runtimeOutcomeReceived       = "received"
	runtimeOutcomeErrorHost      = "error_host"
	runtimeOutcomeStartCancelled = "start_cancelled"
	runtimeErrorNoActiveStart    = 0x05
)

func NewRuntimeContractModel() *RuntimeContractModel {
	return &RuntimeContractModel{}
}

func (model *RuntimeContractModel) Snapshot() RuntimeSessionSnapshot {
	if model == nil {
		return RuntimeSessionSnapshot{}
	}
	return model.session
}

func (model *RuntimeContractModel) Init(features byte) RuntimeExchangeSample {
	before := model.Snapshot()
	requestHex := encodeRuntimeRequest(ENHReqInit, features)
	responseHex := EncodeENHStream(ENHResResetted, features)
	request := runtimeENHFrameSample(ENHReqInit, "init", features, len(requestHex))
	response := runtimeENHFrameSample(ENHResResetted, "resetted", features, len(responseHex))

	model.session = RuntimeSessionSnapshot{}

	return RuntimeExchangeSample{
		RequestHex:  hexStrings(requestHex),
		Request:     request,
		ResponseHex: hexStrings(responseHex),
		Response:    &response,
		Outcome:     runtimeOutcomeResetted,
		StateBefore: before,
		StateAfter:  model.Snapshot(),
	}
}

func (model *RuntimeContractModel) Start(initiator byte) RuntimeExchangeSample {
	before := model.Snapshot()
	requestHex := encodeRuntimeRequest(ENHReqStart, initiator)
	request := runtimeENHFrameSample(ENHReqStart, "start", initiator, len(requestHex))

	if initiator == 0xAA {
		model.session = RuntimeSessionSnapshot{}
		return RuntimeExchangeSample{
			RequestHex:  hexStrings(requestHex),
			Request:     request,
			ResponseHex: nil,
			Response:    nil,
			Outcome:     runtimeOutcomeStartCancelled,
			StateBefore: before,
			StateAfter:  model.Snapshot(),
		}
	}

	model.session = RuntimeSessionSnapshot{
		Active:    true,
		Initiator: initiator,
	}
	responseHex := EncodeENHStream(ENHResStarted, initiator)
	response := runtimeENHFrameSample(ENHResStarted, "started", initiator, len(responseHex))

	return RuntimeExchangeSample{
		RequestHex:  hexStrings(requestHex),
		Request:     request,
		ResponseHex: hexStrings(responseHex),
		Response:    &response,
		Outcome:     runtimeOutcomeStarted,
		StateBefore: before,
		StateAfter:  model.Snapshot(),
	}
}

func (model *RuntimeContractModel) Send(data byte) RuntimeExchangeSample {
	before := model.Snapshot()
	requestHex := encodeRuntimeRequest(ENHReqSend, data)
	request := runtimeENHFrameSample(ENHReqSend, "send", data, len(requestHex))

	if !model.session.Active {
		responseHex := EncodeENHStream(ENHResErrorHost, runtimeErrorNoActiveStart)
		response := runtimeENHFrameSample(ENHResErrorHost, "error_host", runtimeErrorNoActiveStart, len(responseHex))
		return RuntimeExchangeSample{
			RequestHex:  hexStrings(requestHex),
			Request:     request,
			ResponseHex: hexStrings(responseHex),
			Response:    &response,
			Outcome:     runtimeOutcomeErrorHost,
			StateBefore: before,
			StateAfter:  model.Snapshot(),
		}
	}

	responseHex := EncodeENHStream(ENHResReceived, data)
	response := runtimeENHFrameSample(ENHResReceived, "received", data, len(responseHex))
	if data == 0xAA {
		model.session = RuntimeSessionSnapshot{}
	}
	return RuntimeExchangeSample{
		RequestHex:  hexStrings(requestHex),
		Request:     request,
		ResponseHex: hexStrings(responseHex),
		Response:    &response,
		Outcome:     runtimeOutcomeReceived,
		StateBefore: before,
		StateAfter:  model.Snapshot(),
	}
}

func DefaultRuntimeContractReport() RuntimeContractReport {
	model := NewRuntimeContractModel()
	init := model.Init(0x01)
	startSuccess := model.Start(0x31)
	sendAfterStart := model.Send(0x55)
	sendCancel := model.Send(0xAA)
	sendWithoutStart := NewRuntimeContractModel().Send(0x55)
	startCancel := model.Start(0xAA)

	return RuntimeContractReport{
		Init:             init,
		StartSuccess:     startSuccess,
		SendAfterStart:   sendAfterStart,
		SendCancel:       sendCancel,
		SendWithoutStart: sendWithoutStart,
		StartCancel:      startCancel,
	}
}

func runtimeENHFrameSample(command ENHCommand, name string, data byte, consumed int) ENHDecodedFrame {
	return ENHDecodedFrame{
		Command:     command,
		CommandName: name,
		Data:        data,
		Consumed:    consumed,
	}
}

func encodeRuntimeRequest(command ENHCommand, data byte) []byte {
	encoded := EncodeENH(command, data)
	return []byte{encoded[0], encoded[1]}
}
