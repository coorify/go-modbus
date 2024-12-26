package modbus

import (
	"context"
	"io"
	"time"
)

var TIMEOUT time.Duration = 5 * time.Second

func doRequest(stream io.ReadWriter, encoder ModbusEncoder, req *ModbusPDU) (rep *ModbusPDU, err error) {
	reqStream, err := encoder.Endode(req)
	if err != nil {
		return
	}

	_, err = reqStream.WriteTo(stream)
	if err != nil {
		return
	}

	repStream := NewStream(nil)
	timeout, cannel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cannel()

	for {
		select {
		case <-timeout.Done():
			err = timeout.Err()
			return
		default:
			_, err = repStream.ReadFrom(stream)
			if err != nil {
				return
			}

			repStream.Seek(0, io.SeekStart)
			rep, err = encoder.Decode(repStream)
			if err == io.ErrShortBuffer {
				time.Sleep(5 * time.Millisecond)
			} else if err != nil || rep != nil {
				return
			}
		}
	}
}

func uint16ToBytes(values ...uint16) []byte {
	raw := make([]byte, len(values)*2)
	for i, v := range values {
		raw[i*2+0] = byte(v >> 8)
		raw[i*2+1] = byte(v)
	}
	return raw
}

func bytesToUint16(values ...byte) []uint16 {
	raw := make([]uint16, len(values)/2)
	for i := 0; i < len(raw); i++ {
		raw[i] = uint16(values[i*2+0]) << 8
		raw[i] |= uint16(values[i*2+1]) << 0
	}

	return raw
}

func bytesToBool(values ...byte) []bool {

	raw := make([]bool, len(values)*8)
	for i, v := range values {
		for j := 0; j < 8; j++ {
			mask := byte(1) << j

			raw[i*8+j] = (v & mask) == mask
		}
	}

	return raw
}

func boolToBytes(values ...bool) []byte {
	ln := len(values)
	bln := ln / 8
	if ln%8 != 0 {
		bln++
	}

	raw := make([]byte, bln)
	for i, v := range values {
		if v {
			shift := byte(i % 8)
			raw[i/8] |= byte(1) << shift
		}
	}

	return raw
}
