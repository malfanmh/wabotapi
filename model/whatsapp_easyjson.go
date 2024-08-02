// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel(in *jlexer.Lexer, out *WATemplate) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "language":
			out.Language = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "category":
			out.Category = string(in.String())
		case "components":
			if in.IsNull() {
				in.Skip()
				out.Components = nil
			} else {
				in.Delim('[')
				if out.Components == nil {
					if !in.IsDelim(']') {
						out.Components = make([]map[string]interface{}, 0, 8)
					} else {
						out.Components = []map[string]interface{}{}
					}
				} else {
					out.Components = (out.Components)[:0]
				}
				for !in.IsDelim(']') {
					var v1 map[string]interface{}
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('{')
						v1 = make(map[string]interface{})
						for !in.IsDelim('}') {
							key := string(in.String())
							in.WantColon()
							var v2 interface{}
							if m, ok := v2.(easyjson.Unmarshaler); ok {
								m.UnmarshalEasyJSON(in)
							} else if m, ok := v2.(json.Unmarshaler); ok {
								_ = m.UnmarshalJSON(in.Raw())
							} else {
								v2 = in.Interface()
							}
							(v1)[key] = v2
							in.WantComma()
						}
						in.Delim('}')
					}
					out.Components = append(out.Components, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel(out *jwriter.Writer, in WATemplate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"language\":"
		out.RawString(prefix)
		out.String(string(in.Language))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"category\":"
		out.RawString(prefix)
		out.String(string(in.Category))
	}
	{
		const prefix string = ",\"components\":"
		out.RawString(prefix)
		if in.Components == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Components {
				if v3 > 0 {
					out.RawByte(',')
				}
				if v4 == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
					out.RawString(`null`)
				} else {
					out.RawByte('{')
					v5First := true
					for v5Name, v5Value := range v4 {
						if v5First {
							v5First = false
						} else {
							out.RawByte(',')
						}
						out.String(string(v5Name))
						out.RawByte(':')
						if m, ok := v5Value.(easyjson.Marshaler); ok {
							m.MarshalEasyJSON(out)
						} else if m, ok := v5Value.(json.Marshaler); ok {
							out.Raw(m.MarshalJSON())
						} else {
							out.Raw(json.Marshal(v5Value))
						}
					}
					out.RawByte('}')
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WATemplate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WATemplate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WATemplate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WATemplate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel1(in *jlexer.Lexer, out *WAMessageText) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "body":
			out.Body = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel1(out *jwriter.Writer, in WAMessageText) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"body\":"
		out.RawString(prefix[1:])
		out.String(string(in.Body))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAMessageText) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAMessageText) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAMessageText) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAMessageText) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel1(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel2(in *jlexer.Lexer, out *WAMessageContext) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "from":
			out.From = string(in.String())
		case "to":
			out.To = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel2(out *jwriter.Writer, in WAMessageContext) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"from\":"
		out.RawString(prefix[1:])
		out.String(string(in.From))
	}
	{
		const prefix string = ",\"to\":"
		out.RawString(prefix)
		out.String(string(in.To))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAMessageContext) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAMessageContext) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAMessageContext) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAMessageContext) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel2(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel3(in *jlexer.Lexer, out *WAMessageButton) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "text":
			out.Text = string(in.String())
		case "payload":
			out.Payload = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel3(out *jwriter.Writer, in WAMessageButton) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"payload\":"
		out.RawString(prefix)
		out.String(string(in.Payload))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAMessageButton) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAMessageButton) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAMessageButton) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAMessageButton) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel3(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel4(in *jlexer.Lexer, out *WAMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "context":
			(out.Context).UnmarshalEasyJSON(in)
		case "from":
			out.From = string(in.String())
		case "id":
			out.ID = string(in.String())
		case "timestamp":
			out.Timestamp = string(in.String())
		case "text":
			(out.Text).UnmarshalEasyJSON(in)
		case "button":
			(out.Button).UnmarshalEasyJSON(in)
		case "interactive":
			(out.Interactive).UnmarshalEasyJSON(in)
		case "type":
			out.Type = WAMessageType(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel4(out *jwriter.Writer, in WAMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"context\":"
		out.RawString(prefix[1:])
		(in.Context).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"from\":"
		out.RawString(prefix)
		out.String(string(in.From))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"timestamp\":"
		out.RawString(prefix)
		out.String(string(in.Timestamp))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		(in.Text).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"button\":"
		out.RawString(prefix)
		(in.Button).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"interactive\":"
		out.RawString(prefix)
		(in.Interactive).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel4(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel5(in *jlexer.Lexer, out *WAInteractiveListReplay) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel5(out *jwriter.Writer, in WAInteractiveListReplay) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAInteractiveListReplay) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAInteractiveListReplay) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAInteractiveListReplay) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAInteractiveListReplay) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel5(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel6(in *jlexer.Lexer, out *WAInteractive) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			out.Type = string(in.String())
		case "list_reply":
			(out.ListReplay).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel6(out *jwriter.Writer, in WAInteractive) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix[1:])
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"list_reply\":"
		out.RawString(prefix)
		(in.ListReplay).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAInteractive) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAInteractive) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAInteractive) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAInteractive) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel6(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel7(in *jlexer.Lexer, out *WAContactProfile) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = WAContactProfileName(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel7(out *jwriter.Writer, in WAContactProfile) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAContactProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAContactProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAContactProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAContactProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel7(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel8(in *jlexer.Lexer, out *WAContact) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "wa_id":
			out.WaID = string(in.String())
		case "profile":
			(out.Profile).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel8(out *jwriter.Writer, in WAContact) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"wa_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.WaID))
	}
	{
		const prefix string = ",\"profile\":"
		out.RawString(prefix)
		(in.Profile).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WAContact) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WAContact) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WAContact) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WAContact) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel8(l, v)
}
func easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel9(in *jlexer.Lexer, out *Member) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "contact":
			(out.Contact).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel9(out *jwriter.Writer, in Member) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"contact\":"
		out.RawString(prefix)
		(in.Contact).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Member) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Member) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson25b52c02EncodeGithubComMalfanmhWabotapiModel9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Member) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Member) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson25b52c02DecodeGithubComMalfanmhWabotapiModel9(l, v)
}
