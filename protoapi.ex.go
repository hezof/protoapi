package protoapi

func (e *Error) EncodeJSON(w *JsonEncoder) {
	EncodeMessage(w, e, func(w *JsonEncoder, m *Error) {
		EncodeUint32_WithEmpty(w, profile.ResultCodeField, e.Code)
		EncodeString_OmitEmpty(w, profile.ResultNameField, e.Name)
		EncodeString_OmitEmpty(w, profile.ResultMessageField, e.Message)
	})
}

func (e *Error) Error() string {
	out := NewJsonEncoder(nil, 1024)
	e.EncodeJSON(out)
	bs, _ := out.Close()
	return UnsafeString(bs)
}

var _ error = (*Error)(nil)
var _ MessageEncoder = (*Error)(nil)
