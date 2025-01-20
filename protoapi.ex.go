package protoapi

func (e *Error) DecodeField(r *JsonDecoder, f string) {
	switch f {
	case profile.ResultCodeField:
		DecodeUint32(r, &e.Code)
	case profile.ResultNameField:
		DecodeString(r, &e.Name)
	case profile.ResultMessageField:
		DecodeString(r, &e.Message)
	}
}

func (e *Error) EncodeField(w *JsonEncoder) {
	EncodeUint32_WithEmpty(w, profile.ResultCodeField, e.Code)
	EncodeString_OmitEmpty(w, profile.ResultNameField, e.Name)
	EncodeString_OmitEmpty(w, profile.ResultMessageField, e.Message)
}

func (e *Error) Error() string {
	w := NewJsonEncoder(nil, 1024)
	EncodeMessage(w, e)
	return UnsafeString(w.buff)
}

var _ error = (*Error)(nil)
var _ FieldCodec = (*Error)(nil)
