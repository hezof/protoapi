package protojson

import "encoding/json"

func EncodeAny_OmitEmpty(w *JsonEncoder, name string, val any) {
	if val != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		if jc, ok := val.(JsonCodec); ok {
			jc.EncodeJSON(w)
		} else {
			bs, err := json.Marshal(val)
			if err != nil {
				if w.firstError == nil {
					w.firstError = err
				}
				return
			}
			_, _ = w.Write(bs)
		}
		w.buff = append(w.buff, comma)
	}
}

func EncodeAny_WithEmpty(w *JsonEncoder, name string, val any) {
	if val != nil {
		w.ensure(5 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon)
		if jc, ok := val.(JsonCodec); ok {
			jc.EncodeJSON(w)
		} else {
			enc := json.NewEncoder(w)
			enc.SetEscapeHTML(false)
			err := enc.Encode(val)
			if err != nil {
				if w.firstError == nil {
					w.firstError = err
				}
				return
			}
		}
		w.buff = append(w.buff, comma)
	} else {
		w.ensure(8 + len(name))
		w.buff = append(w.buff, quotes)
		w.buff = append(w.buff, name...)
		w.buff = append(w.buff, quotes, colon, 'n', 'u', 'l', 'l', comma)
	}
}
