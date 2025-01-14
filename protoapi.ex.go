package protoapi

func (e *Error) DecodeJSON(r *JsonDecoder) {
	DecodeMessage(r, &e, func(r *JsonDecoder, sr *Error, k string) {
		switch k {
		case profile.ResultCodeField:
			DecodeUint32(r, &sr.Code)
		case profile.ResultNameField:
			DecodeString(r, &sr.Name)
		case profile.ResultMessageField:
			DecodeString(r, &sr.Message)
		}
	})
}

func (e *Error) EncodeJSON(w *JsonEncoder) {
	EncodeMessage(w, e, func(w *JsonEncoder, m *Error) {
		EncodeUint32_WithEmpty(w, profile.ResultCodeField, e.Code)
		EncodeString_OmitEmpty(w, profile.ResultNameField, e.Name)
		EncodeString_OmitEmpty(w, profile.ResultMessageField, e.Message)
	})
}

func (e *Error) Error() string {
	w := NewJsonEncoder(nil, 1024)
	e.EncodeJSON(w)
	_ = w.Close()
	return UnsafeString(w.buff)
}

var _ error = (*Error)(nil)
var _ JsonCodec = (*Error)(nil)
