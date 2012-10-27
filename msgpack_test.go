package msgpack_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/vmihailenco/msgpack"

	msgpack2 "github.com/ugorji/go-msgpack"
	. "launchpad.net/gocheck"
)

func Test(t *testing.T) { TestingT(t) }

type MsgpackTest struct {
	buf *bytes.Buffer
	enc *msgpack.Encoder
	dec *msgpack.Decoder
}

var _ = Suite(&MsgpackTest{})

func (t *MsgpackTest) SetUpTest(c *C) {
	t.buf = &bytes.Buffer{}
	t.enc = msgpack.NewEncoder(t.buf)
	t.dec = msgpack.NewDecoder(bufio.NewReader(t.buf))
}

func (t *MsgpackTest) TestUint(c *C) {
	table := []struct {
		v uint
		b []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{2, []byte{0x02}},
		{125, []byte{0x7d}},
		{126, []byte{0x7e}},
		{127, []byte{0x7f}},
		{128, []byte{0xcc, 0x80}},
		{253, []byte{0xcc, 0xfd}},
		{254, []byte{0xcc, 0xfe}},
		{255, []byte{0xcc, 0xff}},
		{256, []byte{0xcd, 0x01, 0x00}},
		{65533, []byte{0xcd, 0xff, 0xfd}},
		{65534, []byte{0xcd, 0xff, 0xfe}},
		{65535, []byte{0xcd, 0xff, 0Xff}},
		{65536, []byte{0xce, 0x00, 0x01, 0x00, 0x00}},
		{4294967293, []byte{0xce, 0xff, 0xff, 0xff, 0xfd}},
		{4294967294, []byte{0xce, 0xff, 0xff, 0xff, 0xfe}},
		{4294967295, []byte{0xce, 0xff, 0xff, 0xff, 0xff}},
		{4294967296, []byte{0xcf, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},
		{18446744073709551613, []byte{0xcf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfd}},
		{18446744073709551614, []byte{0xcf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}},
		{18446744073709551615, []byte{0xcf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}
	for _, r := range table {
		c.Assert(t.enc.Encode(r.v), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, r.b, Commentf("err encoding %v", r.v))
		var v uint
		c.Assert(t.dec.Decode(&v), IsNil)
		c.Assert(v, Equals, r.v)
	}
}

func (t *MsgpackTest) TestInt(c *C) {
	table := []struct {
		v int
		b []byte
	}{
		{-9223372036854775808, []byte{0xd3, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{-9223372036854775807, []byte{0xd3, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{-9223372036854775806, []byte{0xd3, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}},
		{-2147483651, []byte{0xd3, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xfd}},
		{-2147483650, []byte{0xd3, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xfe}},
		{-2147483649, []byte{0xd3, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff}},
		{-2147483648, []byte{0xd2, 0x80, 0x00, 0x00, 0x00}},
		{-2147483647, []byte{0xd2, 0x80, 0x00, 0x00, 0x01}},
		{-2147483646, []byte{0xd2, 0x80, 0x00, 0x00, 0x002}},
		{-32771, []byte{0xd2, 0xff, 0xff, 0x7f, 0xfd}},
		{-32770, []byte{0xd2, 0xff, 0xff, 0x7f, 0xfe}},
		{-32769, []byte{0xd2, 0xff, 0xff, 0x7f, 0xff}},
		{-32768, []byte{0xd1, 0x80, 0x00}},
		{-32767, []byte{0xd1, 0x80, 0x01}},
		{-131, []byte{0xd1, 0xff, 0x7d}},
		{-130, []byte{0xd1, 0xff, 0x7e}},
		{-129, []byte{0xd1, 0xff, 0x7f}},
		{-128, []byte{0xd0, 0x80}},
		{-127, []byte{0xd0, 0x81}},
		{-34, []byte{0xd0, 0xde}},
		{-33, []byte{0xd0, 0xdf}},
		{-32, []byte{0xe0}},
		{-31, []byte{0xe1}},
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{126, []byte{0x7e}},
		{127, []byte{0x7f}},
		{128, []byte{0xd1, 0x00, 0x80}},
		{129, []byte{0xd1, 0x00, 0x81}},
		{130, []byte{0xd1, 0x00, 0x82}},
		{32765, []byte{0xd1, 0x7f, 0xfd}},
		{32766, []byte{0xd1, 0x7f, 0xfe}},
		{32767, []byte{0xd1, 0x7f, 0xff}},
		{32768, []byte{0xd2, 0x00, 0x00, 0x80, 0x00}},
		{32769, []byte{0xd2, 0x00, 0x00, 0x80, 0x01}},
		{32770, []byte{0xd2, 0x00, 0x00, 0x80, 0x02}},
		{2147483645, []byte{0xd2, 0x7f, 0xff, 0xff, 0xfd}},
		{2147483646, []byte{0xd2, 0x7f, 0xff, 0xff, 0xfe}},
		{2147483647, []byte{0xd2, 0x7f, 0xff, 0xff, 0xff}},
		{2147483648, []byte{0xd3, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00}},
		{2147483649, []byte{0xd3, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x01}},
		{2147483650, []byte{0xd3, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x02}},
		{4294967296, []byte{0xd3, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},
		{4294967297, []byte{0xd3, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01}},
		{4294967298, []byte{0xd3, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02}},
	}
	for _, r := range table {
		c.Assert(t.enc.Encode(r.v), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, r.b, Commentf("err encoding %v", r.v))
		var v int
		c.Assert(t.dec.Decode(&v), IsNil)
		c.Assert(v, Equals, r.v)
	}
}

func (t *MsgpackTest) TestFloat32(c *C) {
	table := []struct {
		v float32
		b []byte
	}{
		{.1, []byte{0xca, 0x3d, 0xcc, 0xcc, 0xcd}},
		{.2, []byte{0xca, 0x3e, 0x4c, 0xcc, 0xcd}},
		{-.1, []byte{0xca, 0xbd, 0xcc, 0xcc, 0xcd}},
		{-.2, []byte{0xca, 0xbe, 0x4c, 0xcc, 0xcd}},
		{float32(math.Inf(1)), []byte{0xca, 0x7f, 0x80, 0x00, 0x00}},
		{float32(math.Inf(-1)), []byte{0xca, 0xff, 0x80, 0x00, 0x00}},
		{math.MaxFloat32, []byte{0xca, 0x7f, 0x7f, 0xff, 0xff}},
		{math.SmallestNonzeroFloat32, []byte{0xca, 0x0, 0x0, 0x0, 0x1}},
	}
	for _, r := range table {
		c.Assert(t.enc.Encode(r.v), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, r.b, Commentf("err encoding %v", r.v))
		var v float32
		c.Assert(t.dec.Decode(&v), IsNil)
		c.Assert(v, Equals, r.v)
	}

	in := float32(math.NaN())
	c.Assert(t.enc.Encode(in), IsNil)
	var out float32
	c.Assert(t.dec.Decode(&out), IsNil)
	c.Assert(math.IsNaN(float64(out)), Equals, true)
}

func (t *MsgpackTest) TestFloat64(c *C) {
	table := []struct {
		v float64
		b []byte
	}{
		{.1, []byte{0xcb, 0x3f, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{.2, []byte{0xcb, 0x3f, 0xc9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{-.1, []byte{0xcb, 0xbf, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{-.2, []byte{0xcb, 0xbf, 0xc9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{math.Inf(1), []byte{0xcb, 0x7f, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		{math.Inf(-1), []byte{0xcb, 0xff, 0xf0, 0x00, 0x00, 0x0, 0x0, 0x0, 0x0}},
		{math.MaxFloat64, []byte{0xcb, 0x7f, 0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{math.SmallestNonzeroFloat64, []byte{0xcb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}},
	}
	for _, r := range table {
		c.Assert(t.enc.Encode(r.v), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, r.b, Commentf("err encoding %v", r.v))
		var v float64
		c.Assert(t.dec.Decode(&v), IsNil)
		c.Assert(v, Equals, r.v)
	}

	in := math.NaN()
	c.Assert(t.enc.Encode(in), IsNil)
	var out float64
	c.Assert(t.dec.Decode(&out), IsNil)
	c.Assert(math.IsNaN(out), Equals, true)
}

func (t *MsgpackTest) TestBool(c *C) {
	table := []struct {
		v bool
		b []byte
	}{
		{false, []byte{0xc2}},
		{true, []byte{0xc3}},
	}
	for _, r := range table {
		c.Assert(t.enc.Encode(r.v), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, r.b, Commentf("err encoding %v", r.v))
		var v bool
		c.Assert(t.dec.Decode(&v), IsNil)
		c.Assert(v, Equals, r.v)
	}
}

func (t *MsgpackTest) TestBytes(c *C) {
	for _, i := range []struct {
		s string
		b []byte
	}{
		{"", []byte{0xa0}},
		{"a", []byte{0xa1, 'a'}},
		{"hello", append([]byte{0xa5}, []byte("hello")...)},
		{
			"world world world",
			append([]byte{0xb1}, []byte("world world world")...),
		},
		{
			"world world world world world world",
			append([]byte{0xda, 0x0, 0x23}, []byte("world world world world world world")...),
		},
	} {
		c.Assert(t.enc.Encode(i.s), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, i.b)
		var s string
		c.Assert(t.dec.Decode(&s), IsNil)
		c.Assert(s, Equals, i.s)
	}
}

func (t *MsgpackTest) TestNil(c *C) {
	c.Assert(t.enc.Encode(nil), IsNil)
	c.Assert(t.buf.Bytes(), DeepEquals, []byte{0xC0})
	c.Assert(t.dec.Decode(new(string)), IsNil)
}

func (t *MsgpackTest) TestDecodingNil(c *C) {
	c.Assert(t.dec.Decode(nil), NotNil)
}

func (t *MsgpackTest) TestTime(c *C) {
	in := time.Now()
	var out time.Time
	c.Assert(t.enc.Encode(in), IsNil)
	c.Assert(t.dec.Decode(&out), IsNil)
	c.Assert(out.Equal(in), Equals, true)
}

func (t *MsgpackTest) TestIntArray(c *C) {
	for _, i := range []struct {
		s []int
		b []byte
	}{
		{[]int{}, []byte{0x90}},
		{[]int{0}, []byte{0x91, 0x0}},
	} {
		c.Assert(t.enc.Encode(i.s), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, i.b, Commentf("err encoding %v", i.s))
		var s []int
		c.Assert(t.dec.Decode(&s), IsNil)
		c.Assert(s, DeepEquals, i.s)
	}
}

func (t *MsgpackTest) TestMap(c *C) {
	for _, i := range []struct {
		m map[string]string
		b []byte
	}{
		{map[string]string{}, []byte{0x80}},
		{map[string]string{"hello": "world"}, []byte{0x81, 0xa5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0xa5, 0x77, 0x6f, 0x72, 0x6c, 0x64}},
	} {
		c.Assert(t.enc.Encode(i.m), IsNil)
		c.Assert(t.buf.Bytes(), DeepEquals, i.b, Commentf("err encoding %v", i.m))
		var m map[string]string
		c.Assert(t.dec.Decode(&m), IsNil)
		c.Assert(m, DeepEquals, i.m)
	}
}

type struct2 struct {
	Name string
}

type struct1 struct {
	Name    string
	Struct2 *struct2
}

func (t *MsgpackTest) TestNestedStructs(c *C) {
	in := struct1{Name: "hello", Struct2: &struct2{Name: "world"}}
	var out struct1
	c.Assert(t.enc.Encode(in), IsNil)
	c.Assert(t.dec.Decode(&out), IsNil)
	c.Assert(in.Name, Equals, out.Name)
	c.Assert(in.Struct2.Name, Equals, out.Struct2.Name)
}

func (t *MsgpackTest) BenchmarkBool(c *C) {
	var v bool
	for i := 0; i < c.N; i++ {
		t.enc.Encode(true)
		t.dec.Decode(&v)
	}
	c.Assert(t.buf.Len(), Equals, 0)
}

func (t *MsgpackTest) BenchmarkMsgpack2Bool(c *C) {
	buf := &bytes.Buffer{}
	dec := msgpack2.NewDecoder(buf, nil)
	enc := msgpack2.NewEncoder(buf)

	var v bool
	for i := 0; i < c.N; i++ {
		enc.Encode(true)
		dec.Decode(&v)
	}
	c.Assert(t.buf.Len(), Equals, 0)
}

func (t *MsgpackTest) BenchmarkBytes(c *C) {
	in := make([]byte, 1024)
	var out []byte
	for i := 0; i < c.N; i++ {
		t.enc.Encode(in)
		t.dec.Decode(&out)
	}
	c.Assert(t.buf.Len(), Equals, 0)
}

func (t *MsgpackTest) BenchmarkMap(c *C) {
	in := make(map[string]string)
	in["hello"] = "world"
	in["foo"] = "bar"
	var out map[string]string

	for i := 0; i < c.N; i++ {
		t.enc.Encode(in)
		t.dec.Decode(&out)
	}

	c.Assert(t.buf.Len(), Equals, 0)
}

type benchmarkStruct struct {
	Name string
	Age  int
	Tm   time.Time
}

func (t *MsgpackTest) BenchmarkStruct(c *C) {
	in := &benchmarkStruct{Name: "Hello World", Age: math.MaxInt32, Tm: time.Now()}
	out := &benchmarkStruct{}
	for i := 0; i < c.N; i++ {
		t.enc.Encode(in)
		t.dec.Decode(out)
	}
	c.Assert(t.buf.Len(), Equals, 0)
}

func (t *MsgpackTest) BenchmarkMsgpack2Struct(c *C) {
	buf := &bytes.Buffer{}
	dec := msgpack2.NewDecoder(buf, nil)
	enc := msgpack2.NewEncoder(buf)

	in := &benchmarkStruct{Name: "Hello World", Age: math.MaxInt32, Tm: time.Now()}
	out := &benchmarkStruct{}
	for i := 0; i < c.N; i++ {
		enc.Encode(in)
		dec.Decode(out)
	}
	c.Assert(t.buf.Len(), Equals, 0)
}

func (t *MsgpackTest) BenchmarkJSONStruct(c *C) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	dec := json.NewDecoder(buf)

	in := &benchmarkStruct{Name: "Hello World", Age: math.MaxInt32, Tm: time.Now()}
	out := &benchmarkStruct{}
	for i := 0; i < c.N; i++ {
		enc.Encode(in)
		dec.Decode(out)
	}
	c.Assert(buf.Len(), Equals, 0)
}