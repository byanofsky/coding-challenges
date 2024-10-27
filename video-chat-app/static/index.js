// Claude created this client.
class WebSocketClient {
  constructor(url) {
    this.url = url;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 3000; // 3 seconds
  }

  connect() {
    try {
      this.ws = new WebSocket(this.url);

      this.ws.onopen = () => {
        console.log("Connected to WebSocket server");
        this.reconnectAttempts = 0;
        const clientId = Math.floor(Math.random() * 100_000_000) + 1;
        this.sendMessage(`clientId:${clientId}`);
      };

      this.ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("Received message:", message);
        // Handle incoming message
      };

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      this.ws.onclose = () => {
        console.log("WebSocket connection closed");
        this.attemptReconnect();
      };
    } catch (error) {
      console.error("Failed to create WebSocket connection:", error);
      this.attemptReconnect();
    }
  }

  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(
        `Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`
      );

      setTimeout(() => {
        this.connect();
      }, this.reconnectDelay);
    } else {
      console.error("Max reconnection attempts reached");
    }
  }

  sendMessage(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not connected");
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}

const client = new WebSocketClient("ws://localhost:8000/signal");
client.connect();

function findClient() {
  const input = document.getElementById("clientId");
  // Send a message
  client.sendMessage(`findClient:${input.value}`);
}
