package main

import (
	"fmt"
	"log"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://localhost:5502",
		Timeout: 1 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	if err := client.Open(); err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	defer client.Close()

	client.SetUnitId(1)

	// 1. Scrittura di un int16 sul registro 103 (valido nel server)
	var s int16 = -200
	err = client.WriteRegister(103, uint16(s))
	if err != nil {
		log.Fatalf("Failed to write register 103: %v", err)
	}
	fmt.Println("Successfully wrote int16 to register 103")

	// 2. Lettura int16 dal registro 103
	reg16, err := client.ReadRegister(103, modbus.HOLDING_REGISTER)
	if err != nil {
		log.Fatalf("Failed to read register 103: %v", err)
	}
	fmt.Printf("Read register 103 as int16: %d\n", int16(reg16))

	// 3. Scrittura di un uint32 sui registri 200-201 (valido nel server)
	err = client.WriteUint32(200, 123456)
	if err != nil {
		log.Fatalf("Failed to write uint32 to registers 200-201: %v", err)
	}
	fmt.Println("Successfully wrote uint32 (123456) to registers 200-201")

	// 4. Lettura del float32 hardcodato nel server agli Input Registers 300-301 (Read-Only 3.1415)
	valFloat, err := client.ReadFloat32(300, modbus.INPUT_REGISTER)
	if err != nil {
		log.Fatalf("Failed to read input register float32 at 300: %v", err)
	}
	fmt.Println("Read static float32 from input register 300:", valFloat)
}
