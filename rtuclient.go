package modbus

import (
	"fmt"
	"io"
)

type rtuClient struct {
	stream  io.ReadWriter
	encoder ModbusEncoder
}

func NewRTUClient(stream io.ReadWriter) ModbusClient {
	return &rtuClient{stream: stream, encoder: NewRTUEncoder()}
}

func (c *rtuClient) ReadCoils(address, quantity uint16, slave byte) (rep []bool, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 2000)
		return
	}

	pduReq := &ModbusPDU{
		OpCode:  OPCODE_READ_COILS,
		Slave:   slave,
		Payload: uint16ToBytes(address, quantity),
	}

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return
	}

	return bytesToBool(pduRep.Payload...)[:quantity], nil
}

func (c *rtuClient) ReadDiscreteInputs(address, quantity uint16, slave byte) (rep []bool, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 2000)
		return
	}

	pduReq := &ModbusPDU{
		OpCode:  OPCODE_READ_DISCRETE_INPUTS,
		Slave:   slave,
		Payload: uint16ToBytes(address, quantity),
	}

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return
	}

	return bytesToBool(pduRep.Payload...)[:quantity], nil
}

func (c *rtuClient) WriteCoil(address uint16, value bool, slave byte) (bool, error) {
	u16True := uint16(0xFF00)

	u16val := uint16(0)
	if value {
		u16val = u16True
	}

	pduReq := &ModbusPDU{
		OpCode:  OPCODE_WRITE_COIL,
		Slave:   slave,
		Payload: uint16ToBytes(address, u16val),
	}

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return false, err
	}

	rep := bytesToUint16(pduRep.Payload...)
	return rep[1] == u16True, nil
}

func (c *rtuClient) WriteCoils(address uint16, values []bool, slave byte) (rep []uint16, err error) {
	u8vals := boolToBytes(values...)

	pduReq := &ModbusPDU{
		OpCode:  OPCODE_WRITE_COILS,
		Slave:   slave,
		Payload: uint16ToBytes(address, uint16(len(values))),
	}

	pduReq.Payload = append(pduReq.Payload, byte(len(u8vals)))
	pduReq.Payload = append(pduReq.Payload, u8vals...)

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return
	}

	return bytesToUint16(pduRep.Payload...), nil
}

func (c *rtuClient) ReadInputRegisters(address, quantity uint16, slave byte) (rep []uint16, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return
	}

	pduReq := &ModbusPDU{
		OpCode:  OPCODE_READ_INPUT_REGISTERS,
		Slave:   slave,
		Payload: uint16ToBytes(address, quantity),
	}

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return
	}

	return bytesToUint16(pduRep.Payload...), nil
}

func (c *rtuClient) ReadHoldingRegisters(address, quantity uint16, slave byte) (rep []uint16, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return
	}

	pduReq := &ModbusPDU{
		OpCode:  OPCODE_READ_HOLDING_REGISTERS,
		Slave:   slave,
		Payload: uint16ToBytes(address, quantity),
	}

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return
	}

	return bytesToUint16(pduRep.Payload...), nil
}

func (c *rtuClient) WriteRegister(address, value uint16, slave byte) (rep []uint16, err error) {
	pduReq := &ModbusPDU{
		OpCode:  OPCODE_WRITE_REGISTER,
		Slave:   slave,
		Payload: uint16ToBytes(address, value),
	}

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return
	}

	return bytesToUint16(pduRep.Payload...), nil
}

func (c *rtuClient) WriteRegisters(address uint16, values []uint16, slave byte) (rep []uint16, err error) {
	pduReq := &ModbusPDU{
		OpCode:  OPCODE_WRITE_REGISTER,
		Slave:   slave,
		Payload: uint16ToBytes(address, uint16(len(values))),
	}

	pduReq.Payload = append(pduReq.Payload, byte(len(values)*2))
	pduReq.Payload = append(pduReq.Payload, uint16ToBytes(values...)...)

	pduRep, err := doRequest(c.stream, c.encoder, pduReq)
	if err != nil {
		return
	}

	return bytesToUint16(pduRep.Payload...), nil
}
