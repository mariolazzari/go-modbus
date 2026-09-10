# Go-Modbus

A clean, idiomatic, and high-performance Modbus library written in Go (Golang), implementing RTU, ASCII, and TCP/IP communication protocols. Designed following modern Go standards, utilizing the standard library, `pgxpool` principles for robust resource pooling where applicable, and clean concurrency patterns.

## Features

- **Protocols Supported:**
  - Modbus TCP
  - Modbus RTU (Serial)
  - Modbus ASCII (Serial)
- **Roles:**
  - Client (Master)
  - Server (Slave)
- **Design & Performance:**
  - Clean API adhering to Go best practices
  - Thread-safe connection handling and pooling
  - Comprehensive timeout and error management
  - Zero heavy third-party dependencies (pure Go implementation)

---

## Installation

Ensure you have Go installed (Go 1.22+ recommended). Install the package using `go get`:

```bash
go get github.com/mariolazzari/go-modbus
```

---

## Quick Start

### 1. Modbus TCP Client Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mariolazzari/go-modbus/modbus"
)

func main() {
	// Configure TCP client configuration
	clientCfg := &modbus.ClientConfig{
		Address: "192.168.1.100:502",
		Timeout: 5 * time.Second,
	}

	// Create a new TCP client
	client := modbus.NewTCPClient(clientCfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the Modbus device
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to modbus server: %v", err)
	}
	defer client.Close()

	// Read Holding Registers (Function Code 03)
	// Starting address: 0, Quantity: 10
	results, err := client.ReadHoldingRegisters(ctx, 1, 0, 10)
	if err != nil {
		log.Fatalf("failed to read holding registers: %v", err)
	}

	// Output the read register values (strings and comments in English as standard)
	fmt.Printf("Holding Registers Data: %v\n", results)
}
```

### 2. Modbus RTU (Serial) Client Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mariolazzari/go-modbus/modbus"
)

func main() {
	// Configure Serial RTU client configuration
	serialCfg := &modbus.SerialConfig{
		Address:  "/dev/ttyUSB0",
		BaudRate: 9600,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  2 * time.Second,
	}

	client := modbus.NewRTUClient(serialCfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("failed to open serial port: %v", err)
	}
	defer client.Close()

	// Write Single Coil (Function Code 05)
	// Address: 10, Value: true (ON)
	err := client.WriteSingleCoil(ctx, 1, 10, true)
	if err != nil {
		log.Fatalf("failed to write single coil: %v", err)
	}

	fmt.Println("Coil successfully written to TRUE")
}
```

---

## Supported Function Codes

| Function Code | Name                     | Description                               |
| :-----------: | :----------------------- | :---------------------------------------- |
|    **01**     | Read Coils               | Read discrete output status               |
|    **02**     | Read Discrete Inputs     | Read discrete input status                |
|    **03**     | Read Holding Registers   | Read 16-bit holding registers             |
|    **04**     | Read Input Registers     | Read 16-bit input registers               |
|    **05**     | Write Single Coil        | Write a single discrete output            |
|    **06**     | Write Single Register    | Write a single 16-bit register            |
|    **15**     | Write Multiple Coils     | Write a block of discrete outputs         |
|    **16**     | Write Multiple Registers | Write a block of 16-bit holding registers |

---

## Configuration & Options

The library allows fine-grained customization through configuration structs:

- **ClientConfig / SerialConfig**: Set custom timeouts, retry attempts, frame delays, and connection parameters.
- **Context Awareness**: All network and serial operations support standard `context.Context` for proper cancellation and timeout handling.

---

## Testing

To run the unit tests and integration tests:

```bash
go test -v ./...
```

To include race detection and coverage reports:

```bash
go test -race -coverprofile=coverage.out ./...
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository (`https://github.com/mariolazzari/go-modbus/fork`)
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git origin push feature/amazing-feature`)
5. Open a Pull Request

Please ensure all comments, documentation, and code commits follow standard English conventions.

---

## License

Distributed under the MIT License. See `LICENSE` for more information.
