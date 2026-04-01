package adapterproto

type OracleReport struct {
	ENH     OracleENHReport       `json:"enh"`
	ENS     OracleENSReport       `json:"ens"`
	INFO    OracleInfoReport      `json:"info"`
	Runtime RuntimeContractReport `json:"runtime"`
	Scan    OracleScanReport      `json:"scan"`
}

type OracleENHReport struct {
	ShortFormReceived OracleENHSample `json:"short_form_received"`
	EncodedInfo       OracleENHSample `json:"encoded_info"`
}

type OracleENHSample struct {
	InputHex     []string        `json:"input_hex"`
	Decoded      ENHDecodedFrame `json:"decoded"`
	ReencodedHex []string        `json:"reencoded_hex"`
}

type OracleENSReport struct {
	RawHex     []string `json:"raw_hex"`
	EncodedHex []string `json:"encoded_hex"`
	DecodedHex []string `json:"decoded_hex"`
}

type OracleInfoReport struct {
	Version   OracleVersionSample     `json:"version"`
	ResetInfo OracleResetSample       `json:"reset_info"`
	Queries   []OracleInfoQuerySample `json:"queries"`
}

type OracleVersionSample struct {
	InputHex           []string       `json:"input_hex"`
	Parsed             AdapterVersion `json:"parsed"`
	SupportsInfoQuery0 bool           `json:"supports_info_query_0"`
	SupportsInfoQuery6 bool           `json:"supports_info_query_6"`
}

type OracleResetSample struct {
	InputHex []string         `json:"input_hex"`
	Parsed   AdapterResetInfo `json:"parsed"`
}

type OracleInfoQuerySample struct {
	ID          AdapterInfoID     `json:"id"`
	Name        string            `json:"name"`
	Supported   bool              `json:"supported"`
	Modeled     bool              `json:"modeled"`
	RequestHex  []string          `json:"request_hex"`
	ResponseHex []string          `json:"response_hex,omitempty"`
	Version     *AdapterVersion   `json:"version,omitempty"`
	ResetInfo   *AdapterResetInfo `json:"reset_info,omitempty"`
	Notes       string            `json:"notes,omitempty"`
}

func DefaultOracleReport() OracleReport {
	shortInput := []byte{0x5A}
	shortDecoded, _ := DecodeENHStream(shortInput)
	shortReencoded := EncodeENHStream(shortDecoded.Command, shortDecoded.Data)

	encodedInput := EncodeENHStream(ENHReqInfo, 0x06)
	encodedDecoded, _ := DecodeENHStream(encodedInput)
	encodedReencoded := EncodeENHStream(encodedDecoded.Command, encodedDecoded.Data)

	ensRaw := []byte{0x11, ENSByteEscape, ENSByteSync, 0x7F}
	ensEncoded := EncodeENS(ensRaw)
	ensDecoded, _ := DecodeENS(ensEncoded)

	versionInput := []byte{0x12, 0x01, 0xAB, 0xCD, 0x1E, 0x34, 0x56, 0x78}
	versionParsed, _ := ParseAdapterVersion(versionInput)

	resetInput := []byte{0x03, 0x07}
	resetParsed, _ := ParseAdapterResetInfo(resetInput)

	infoQueries := buildInfoQuerySamples(versionInput, resetInput, versionParsed, resetParsed)

	return OracleReport{
		ENH: OracleENHReport{
			ShortFormReceived: OracleENHSample{
				InputHex:     hexStrings(shortInput),
				Decoded:      shortDecoded,
				ReencodedHex: hexStrings(shortReencoded),
			},
			EncodedInfo: OracleENHSample{
				InputHex:     hexStrings(encodedInput),
				Decoded:      encodedDecoded,
				ReencodedHex: hexStrings(encodedReencoded),
			},
		},
		ENS: OracleENSReport{
			RawHex:     hexStrings(ensRaw),
			EncodedHex: hexStrings(ensEncoded),
			DecodedHex: hexStrings(ensDecoded),
		},
		INFO: OracleInfoReport{
			Version: OracleVersionSample{
				InputHex:           hexStrings(versionInput),
				Parsed:             versionParsed,
				SupportsInfoQuery0: versionParsed.SupportsInfoID(AdapterInfoVersion),
				SupportsInfoQuery6: versionParsed.SupportsInfoID(AdapterInfoResetInfo),
			},
			ResetInfo: OracleResetSample{
				InputHex: hexStrings(resetInput),
				Parsed:   resetParsed,
			},
			Queries: infoQueries,
		},
		Runtime: DefaultRuntimeContractReport(),
		Scan:    DefaultOracleScanReport(),
	}
}

func buildInfoQuerySamples(versionPayload, resetPayload []byte, version AdapterVersion, reset AdapterResetInfo) []OracleInfoQuerySample {
	queries := make([]OracleInfoQuerySample, 0, 8)
	for id := AdapterInfoVersion; id <= AdapterInfoWiFiRSSI; id++ {
		query := OracleInfoQuerySample{
			ID:         id,
			Name:       id.String(),
			Supported:  version.SupportsInfoID(id),
			Modeled:    false,
			RequestHex: hexStrings(EncodeENHStream(ENHReqInfo, byte(id))),
			Notes:      "supported by adapter version, but no response parser is modeled here yet",
		}
		switch id {
		case AdapterInfoVersion:
			query.Modeled = true
			query.ResponseHex = hexStrings(versionPayload)
			query.Version = &version
			query.Notes = "version response is modeled as the canonical INFO 0x00 payload"
		case AdapterInfoResetInfo:
			query.Modeled = true
			query.ResponseHex = hexStrings(resetPayload)
			query.ResetInfo = &reset
			query.Notes = "reset-info response is modeled as the canonical INFO 0x06 payload"
		}
		queries = append(queries, query)
	}
	return queries
}

func hexStrings(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	out := make([]string, len(data))
	for idx, value := range data {
		out[idx] = hexByte(value)
	}
	return out
}

func hexByte(value byte) string {
	const digits = "0123456789abcdef"
	return string([]byte{digits[value>>4], digits[value&0x0F]})
}
