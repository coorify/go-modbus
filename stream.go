package modbus

import (
	"bytes"
	"errors"
	"io"
)

type modbusStream struct {
	buff []byte
	rpos int64
	wpos int64
}

func NewStream(bs []byte) ModbusStream {
	if bs == nil {
		bs = make([]byte, 0)
	}

	return &modbusStream{buff: bs, rpos: 0, wpos: int64(len(bs))}
}

func (s *modbusStream) Bytes() []byte {
	if s.wpos <= s.rpos {
		return make([]byte, 0)
	}

	return s.buff[s.rpos:s.wpos]
}

func (s *modbusStream) Len() int64 {
	return s.wpos - s.rpos
}

func (s *modbusStream) grow(n int) {
	rml := len(s.buff) - int(s.wpos)
	if rml < n {
		chunk := make([]byte, n-rml)
		s.buff = append(s.buff, chunk...)
	}
}

func (s *modbusStream) WriteTo(w io.Writer) (n int64, err error) {
	r := bytes.NewReader(s.Bytes())

	n, err = io.Copy(w, r)
	if n > 0 {
		s.rpos += n
	}

	return
}

func (s *modbusStream) ReadFrom(r io.Reader) (n int64, err error) {
	raw := make([]byte, 512)

	i32, err := r.Read(raw)
	n = int64(i32)
	if err != nil {
		return
	}

	if i32 > 0 {
		i32, err = s.Write(raw[:i32])
		n = int64(i32)
	}

	return
}

func (s *modbusStream) Seek(offset int64, whence int) (int64, error) {
	var abs int64

	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = s.rpos + offset
	case io.SeekEnd:
		abs = s.wpos + offset
	default:
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}

	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}

	s.rpos = abs
	return abs, nil
}

func (s *modbusStream) ReadByte() (byte, error) {
	raw := make([]byte, 1)
	_, err := s.Read(raw)
	if err != nil {
		return 0, err
	}
	return raw[0], nil
}

func (s *modbusStream) Read(p []byte) (n int, err error) {
	bl := s.Len()
	el := int64(len(p))

	if bl <= 0 {
		return 0, io.EOF
	} else if el > bl {
		return 0, io.ErrShortBuffer
	}

	n = copy(p, s.buff[s.rpos:])
	s.rpos += int64(n)

	return
}

func (s *modbusStream) WriteByte(v byte) error {
	_, err := s.Write([]byte{v})
	return err
}

func (s *modbusStream) Write(p []byte) (n int, err error) {
	s.grow(len(p))
	n = copy(s.buff[s.wpos:], p)
	s.wpos += int64(n)
	return
}
