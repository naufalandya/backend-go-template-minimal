package websocket

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebSocketHandler(c *websocket.Conn) {
	defer func() {
		fmt.Println("Closing connection~ 🌸")
		c.Close()
	}()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		fmt.Printf("received: %s\n", msg)

		reply := fmt.Sprintf("Echo: %s | time: %s", msg, time.Now().Format("15:04:05"))
		if err = c.WriteMessage(mt, []byte(reply)); err != nil {
			fmt.Println("write error:", err)
			break
		}
	}
}

func UpgradeMiddleware() fiber.Handler {
	return websocket.New(WebSocketHandler)
}
