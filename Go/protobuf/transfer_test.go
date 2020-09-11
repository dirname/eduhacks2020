package protobuf

import (
	"bytes"
	"testing"
)

type RequestTest struct {
	Req      *Request
	Type     int32
	Path     string
	Data     []byte
	Sign     []byte
	Location string
}

type ResponseTest struct {
	Res  *Response
	Code int32
	Msg  string
	Type int32
	Data []byte
	Html *Render
}

var testRes = &Response{
	Code:   0,
	Msg:    "OK",
	Type:   0,
	Data:   []byte("test"),
	Render: true,
	Html: &Render{
		Code:   "<script>alert(\"hello world!\");</script>",
		Type:   0,
		Id:     "msgBox",
		Iframe: false,
	},
}

var testReq1 = &Request{
	Type:     1,
	Path:     "Hello",
	Data:     []byte("Hello world !"),
	Sign:     nil,
	Location: "/login",
}

var testReq2 = &Request{
	Type:     2,
	Path:     "World",
	Data:     nil,
	Sign:     []byte("Hello world !"),
	Location: "/logout",
}

var testReq3 = &Request{
	Type:     3,
	Path:     "!",
	Data:     nil,
	Sign:     nil,
	Location: "",
}

var requestTests = []RequestTest{
	{Type: 1, Path: "Hello", Data: []byte("Hello world !"), Sign: nil, Location: "/login", Req: testReq1},
	{Type: 2, Path: "World", Data: nil, Sign: []byte("Hello world !"), Location: "/logout", Req: testReq2},
	{Type: 3, Path: "!", Data: nil, Sign: nil, Location: "", Req: testReq3},
}

var responseTests = []ResponseTest{
	{Code: 0, Msg: "OK", Type: 0, Data: []byte("test"), Html: &Render{
		Code:   "<script>alert(\"hello world!\");</script>",
		Type:   0,
		Id:     "msgBox",
		Iframe: false,
	}, Res: testRes},
}

func TestRequest(t *testing.T) {
	for _, test := range requestTests {
		gotType := test.Req.GetType()
		if gotType != test.Type {
			t.Errorf("Type = %v, want %v", gotType, test.Type)
		}
		gotPath := test.Req.GetPath()
		if gotPath != test.Path {
			t.Errorf("Path = %v, want %v", gotPath, test.Path)
		}
		gotData := test.Req.GetData()
		if !bytes.Equal(gotData, test.Data) {
			t.Errorf("Data = %v, want %v", gotData, test.Data)
		}
		gotSign := test.Req.GetSign()
		if !bytes.Equal(gotSign, test.Sign) {
			t.Errorf("Sign = %v, want %v", gotSign, test.Sign)
		}
		gotLocation := test.Req.GetLocation()
		if gotLocation != test.Location {
			t.Errorf("Location = %v, want %v", gotLocation, test.Location)
		}
	}
}

func TestResponse(t *testing.T) {
	for _, test := range responseTests {
		gotType := test.Res.GetType()
		if gotType != test.Type {
			t.Errorf("Type = %v, want %v", gotType, test.Type)
		}
		gotData := test.Res.GetData()
		if !bytes.Equal(gotData, test.Data) {
			t.Errorf("Data = %v, want %v", gotData, test.Data)
		}
		gotMsg := test.Res.GetMsg()
		if gotMsg != test.Msg {
			t.Errorf("Msg = %v, want %v", gotMsg, test.Msg)
		}
		gotCode := test.Res.GetCode()
		if gotCode != test.Code {
			t.Errorf("Code = %v, want %v", gotCode, test.Code)
		}
		gotRenderCode := test.Res.Html.GetCode()
		if gotRenderCode != test.Html.Code {
			t.Errorf("Render Code = %v, want %v", gotRenderCode, test.Html.Code)
		}
		gotRenderType := test.Res.Html.GetType()
		if gotRenderType != test.Html.Type {
			t.Errorf("Render Type = %v, want %v", gotRenderType, test.Html.Type)
		}
		gotRenderId := test.Res.Html.GetId()
		if gotRenderId != test.Html.Id {
			t.Errorf("Render Id = %v, want %v", gotRenderId, test.Html.Id)
		}
		gotRenderIframe := test.Res.Html.GetIframe()
		if gotRenderIframe != test.Html.GetIframe() {
			t.Errorf("Render Iframe = %v, want %v", gotRenderIframe, test.Html.Iframe)
		}
	}
}
