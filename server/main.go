package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/simonvetter/modbus"
)

type DeviceState struct {
	lock sync.RWMutex

	// --- HOLDING REGISTERS (Read / Write) ---
	// Impostazioni e controlli inviati dal client
	TargetTemp float32 // Reg 100-101: Setpoint Temperatura (°C)
	MaxSpeed   uint32  // Reg 102-103: Velocità massima consentita (RPM)

	// --- INPUT REGISTERS (Read-Only) ---
	// Telemetrie e misure lette dai sensori hardware
	ActualTemp   float32 // Reg 200-201: Temperatura attuale misurata (°C)
	CurrentSpeed uint32  // Reg 202-203: Velocità attuale (RPM)
	FaultCode    uint16  // Reg 204:     Codice di errore (0 = Nessun errore)
	Uptime       uint32  // Reg 205-206: Uptime del sistema in secondi
}

func NewDeviceState() *DeviceState {
	return &DeviceState{
		TargetTemp:   22.5,
		MaxSpeed:     3000,
		ActualTemp:   21.8,
		CurrentSpeed: 1450,
		FaultCode:    0,
		Uptime:       0,
	}
}

// HandleHoldingRegisters: Gestisce esclusivamente le configurazioni modificabili (Read/Write)
func (ds *DeviceState) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	if req.UnitId != 1 {
		return nil, modbus.ErrIllegalFunction
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	for i := 0; i < int(req.Quantity); i++ {
		regAddr := req.Addr + uint16(i)

		switch regAddr {
		// TargetTemp (float32 -> Reg 100-101)
		case 100:
			bits := math.Float32bits(ds.TargetTemp)
			if req.IsWrite {
				newBits := (uint32(req.Args[i]) << 16) | (bits & 0x0000ffff)
				ds.TargetTemp = math.Float32frombits(newBits)
			}
			res = append(res, uint16((bits>>16)&0xffff))
		case 101:
			bits := math.Float32bits(ds.TargetTemp)
			if req.IsWrite {
				newBits := (bits & 0xffff0000) | uint32(req.Args[i])
				ds.TargetTemp = math.Float32frombits(newBits)
			}
			res = append(res, uint16(bits&0xffff))

		// MaxSpeed (uint32 -> Reg 102-103)
		case 102:
			if req.IsWrite {
				ds.MaxSpeed = (uint32(req.Args[i]) << 16) | (ds.MaxSpeed & 0x0000ffff)
			}
			res = append(res, uint16((ds.MaxSpeed>>16)&0xffff))
		case 103:
			if req.IsWrite {
				ds.MaxSpeed = (ds.MaxSpeed & 0xffff0000) | uint32(req.Args[i])
			}
			res = append(res, uint16(ds.MaxSpeed&0xffff))

		default:
			return nil, modbus.ErrIllegalDataAddress
		}
	}

	return res, nil
}

// HandleInputRegisters: Gestisce le misurazioni e la telemetria (Sola Lettura)
func (ds *DeviceState) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	if req.UnitId != 1 {
		return nil, modbus.ErrIllegalFunction
	}

	ds.lock.RLock()
	defer ds.lock.RUnlock()

	for i := 0; i < int(req.Quantity); i++ {
		regAddr := req.Addr + uint16(i)

		switch regAddr {
		// ActualTemp (float32 -> Reg 200-201)
		case 200:
			bits := math.Float32bits(ds.ActualTemp)
			res = append(res, uint16((bits>>16)&0xffff))
		case 201:
			bits := math.Float32bits(ds.ActualTemp)
			res = append(res, uint16(bits&0xffff))

		// CurrentSpeed (uint32 -> Reg 202-203)
		case 202:
			res = append(res, uint16((ds.CurrentSpeed>>16)&0xffff))
		case 203:
			res = append(res, uint16(ds.CurrentSpeed&0xffff))

		// FaultCode (uint16 -> Reg 204)
		case 204:
			res = append(res, ds.FaultCode)

		// Uptime (uint32 -> Reg 205-206)
		case 205:
			res = append(res, uint16((ds.Uptime>>16)&0xffff))
		case 206:
			res = append(res, uint16(ds.Uptime&0xffff))

		default:
			return nil, modbus.ErrIllegalDataAddress
		}
	}

	return res, nil
}

func (ds *DeviceState) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	return nil, modbus.ErrIllegalFunction
}

func (ds *DeviceState) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	return nil, modbus.ErrIllegalFunction
}

func main() {
	state := NewDeviceState()

	// Goroutine per simulare l'aggiornamento dei sensori in background
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			state.lock.Lock()
			state.Uptime++
			// Simula una piccola fluttuazione della temperatura misurata
			state.ActualTemp += 0.05
			state.lock.Unlock()
		}
	}()

	server, err := modbus.NewServer(&modbus.ServerConfiguration{
		URL:        "tcp://localhost:5502",
		Timeout:    30 * time.Second,
		MaxClients: 5,
	}, state)
	if err != nil {
		log.Fatalf("Server creation failed: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
	fmt.Println("Modbus Server running on tcp://localhost:5502...")
	defer server.Stop()

	select {}
}
