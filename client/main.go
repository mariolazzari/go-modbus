package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://localhost:5502",
		Timeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatalf("Client init error: %v", err)
	}

	if err := client.Open(); err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer client.Close()

	client.SetUnitId(1)

	// Intercetta l'interruzione da tastiera (Ctrl+C) o segnale di terminazione
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Ticker configurato ogni 5 secondi
	interval := 5 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("Avvio monitoraggio Modbus ogni %v (Premi Ctrl+C per fermare)...\n\n", interval)

	// Esegue subito la prima lettura senza attendere il primo tick
	readTelemetry(client)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nMonitoraggio interrotto dall'utente. Chiusura client...")
			return

		case t := <-ticker.C:
			fmt.Printf("--- Lettura del %s ---\n", t.Format("15:04:05"))
			readTelemetry(client)
		}
	}
}

// Helper per eseguire la lettura della telemetria
func readTelemetry(client *modbus.ModbusClient) {
	actualTemp, errTemp := client.ReadFloat32(200, modbus.INPUT_REGISTER)
	currentSpeed, errSpeed := client.ReadUint32(202, modbus.INPUT_REGISTER)
	faultCode, errFault := client.ReadRegister(204, modbus.INPUT_REGISTER)
	uptime, errUptime := client.ReadUint32(205, modbus.INPUT_REGISTER)

	if errTemp != nil || errSpeed != nil || errFault != nil || errUptime != nil {
		log.Printf("Errore durante la lettura dei registri Modbus")
		return
	}

	fmt.Printf("  Temperatura Attuale (Reg 200): %.2f °C\n", actualTemp)
	fmt.Printf("  Velocità Attuale    (Reg 202): %d RPM\n", currentSpeed)
	fmt.Printf("  Codice Guasto       (Reg 204): %d\n", faultCode)
	fmt.Printf("  Uptime Sistema      (Reg 205): %d sec\n", uptime)
	fmt.Println()
}
