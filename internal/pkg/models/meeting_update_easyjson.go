// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjson7da70205DecodeKonamiBackendInternalPkgModels(in *jlexer.Lexer, out *MeetingUpdate) {
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
		case "meetId":
			out.MeetId = int(in.Int())
		case "fields":
			if in.IsNull() {
				in.Skip()
				out.Fields = nil
			} else {
				if out.Fields == nil {
					out.Fields = new(MeetUpdateFields)
				}
				easyjson7da70205DecodeKonamiBackendInternalPkgModels1(in, out.Fields)
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
func easyjson7da70205EncodeKonamiBackendInternalPkgModels(out *jwriter.Writer, in MeetingUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"meetId\":"
		out.RawString(prefix[1:])
		out.Int(int(in.MeetId))
	}
	{
		const prefix string = ",\"fields\":"
		out.RawString(prefix)
		if in.Fields == nil {
			out.RawString("null")
		} else {
			easyjson7da70205EncodeKonamiBackendInternalPkgModels1(out, *in.Fields)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MeetingUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7da70205EncodeKonamiBackendInternalPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MeetingUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7da70205EncodeKonamiBackendInternalPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MeetingUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7da70205DecodeKonamiBackendInternalPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MeetingUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7da70205DecodeKonamiBackendInternalPkgModels(l, v)
}
func easyjson7da70205DecodeKonamiBackendInternalPkgModels1(in *jlexer.Lexer, out *MeetUpdateFields) {
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
		case "isRegistered":
			if in.IsNull() {
				in.Skip()
				out.Reg = nil
			} else {
				if out.Reg == nil {
					out.Reg = new(bool)
				}
				*out.Reg = bool(in.Bool())
			}
		case "isLiked":
			if in.IsNull() {
				in.Skip()
				out.Like = nil
			} else {
				if out.Like == nil {
					out.Like = new(bool)
				}
				*out.Like = bool(in.Bool())
			}
		case "card":
			if in.IsNull() {
				in.Skip()
				out.Card = nil
			} else {
				if out.Card == nil {
					out.Card = new(MeetingData)
				}
				(*out.Card).UnmarshalEasyJSON(in)
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
func easyjson7da70205EncodeKonamiBackendInternalPkgModels1(out *jwriter.Writer, in MeetUpdateFields) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"isRegistered\":"
		out.RawString(prefix[1:])
		if in.Reg == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.Reg))
		}
	}
	{
		const prefix string = ",\"isLiked\":"
		out.RawString(prefix)
		if in.Like == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.Like))
		}
	}
	{
		const prefix string = ",\"card\":"
		out.RawString(prefix)
		if in.Card == nil {
			out.RawString("null")
		} else {
			(*in.Card).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}
func easyjson7da70205DecodeKonamiBackendInternalPkgModels2(in *jlexer.Lexer, out *MeetingData) {
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
		case "address":
			if in.IsNull() {
				in.Skip()
				out.Address = nil
			} else {
				if out.Address == nil {
					out.Address = new(string)
				}
				*out.Address = string(in.String())
			}
		case "city":
			if in.IsNull() {
				in.Skip()
				out.City = nil
			} else {
				if out.City == nil {
					out.City = new(string)
				}
				*out.City = string(in.String())
			}
		case "start":
			if in.IsNull() {
				in.Skip()
				out.Start = nil
			} else {
				if out.Start == nil {
					out.Start = new(string)
				}
				*out.Start = string(in.String())
			}
		case "end":
			if in.IsNull() {
				in.Skip()
				out.End = nil
			} else {
				if out.End == nil {
					out.End = new(string)
				}
				*out.End = string(in.String())
			}
		case "meet-description":
			if in.IsNull() {
				in.Skip()
				out.Text = nil
			} else {
				if out.Text == nil {
					out.Text = new(string)
				}
				*out.Text = string(in.String())
			}
		case "meetingTags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Tags = append(out.Tags, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "name":
			if in.IsNull() {
				in.Skip()
				out.Title = nil
			} else {
				if out.Title == nil {
					out.Title = new(string)
				}
				*out.Title = string(in.String())
			}
		case "photo":
			if in.IsNull() {
				in.Skip()
				out.Photo = nil
			} else {
				if out.Photo == nil {
					out.Photo = new(string)
				}
				*out.Photo = string(in.String())
			}
		case "seats":
			if in.IsNull() {
				in.Skip()
				out.Seats = nil
			} else {
				if out.Seats == nil {
					out.Seats = new(int)
				}
				*out.Seats = int(in.Int())
			}
		case "seatsLeft":
			if in.IsNull() {
				in.Skip()
				out.SeatsLeft = nil
			} else {
				if out.SeatsLeft == nil {
					out.SeatsLeft = new(int)
				}
				*out.SeatsLeft = int(in.Int())
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
func easyjson7da70205EncodeKonamiBackendInternalPkgModels2(out *jwriter.Writer, in MeetingData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix[1:])
		if in.Address == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Address))
		}
	}
	{
		const prefix string = ",\"city\":"
		out.RawString(prefix)
		if in.City == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.City))
		}
	}
	{
		const prefix string = ",\"start\":"
		out.RawString(prefix)
		if in.Start == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Start))
		}
	}
	{
		const prefix string = ",\"end\":"
		out.RawString(prefix)
		if in.End == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.End))
		}
	}
	{
		const prefix string = ",\"meet-description\":"
		out.RawString(prefix)
		if in.Text == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Text))
		}
	}
	{
		const prefix string = ",\"meetingTags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Tags {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		if in.Title == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Title))
		}
	}
	{
		const prefix string = ",\"photo\":"
		out.RawString(prefix)
		if in.Photo == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Photo))
		}
	}
	{
		const prefix string = ",\"seats\":"
		out.RawString(prefix)
		if in.Seats == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.Seats))
		}
	}
	{
		const prefix string = ",\"seatsLeft\":"
		out.RawString(prefix)
		if in.SeatsLeft == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.SeatsLeft))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MeetingData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7da70205EncodeKonamiBackendInternalPkgModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MeetingData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7da70205EncodeKonamiBackendInternalPkgModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MeetingData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7da70205DecodeKonamiBackendInternalPkgModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MeetingData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7da70205DecodeKonamiBackendInternalPkgModels2(l, v)
}
