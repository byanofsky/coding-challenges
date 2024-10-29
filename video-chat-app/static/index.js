// Claude created this client.
class WebSocketClient {
  constructor(url) {
    this.url = url;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 3000; // 3 seconds
    this.messageHandlers = [];
  }

  connect(clientId) {
    try {
      this.ws = new WebSocket(this.url + "?clientId=" + clientId);

      this.ws.onopen = () => {
        console.log("Connected to WebSocket server");
        this.reconnectAttempts = 0;
      };

      this.ws.onmessage = (event) => {
        // const message = JSON.parse(event.data);
        const message = event.data;
        console.log("Received message:", message);
        // Handle incoming message
        const splitIdx = message.indexOf(":");
        const type = message.slice(0, splitIdx);
        const rest = message.slice(splitIdx + 1);

        this.messageHandlers.forEach(([t, fn]) => {
          if (t === type) {
            fn(rest);
          }
        });
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

  addMessageHandler(message, fn) {
    this.messageHandlers.push([message, fn]);
  }
}

const client = new WebSocketClient("ws://localhost:8000/signal");

let stream;

async function startVideo() {
  // TODO: Erro rhandling
  stream = await navigator.mediaDevices.getUserMedia({
    video: true,
    audio: false,
  });
  const video = document.getElementById("stream");
  video.srcObject = stream;
}

async function call() {
  const remoteClientId = document.getElementById("remoteClientId").value;
  client.sendMessage({
    to: remoteClientId,
    from: "local client id",
    message: "hello world",
  });
  // const videoTracks = stream.getVideoTracks();
  // if (videoTracks.length > 0) {
  //   console.log(`Using video device: ${videoTracks[0].label}`);
  // }

  // const pc = new RTCPeerConnection();
  // pc.addEventListener("icecandidate", (e) => console.log(e));
  // pc.addEventListener("iceconnectionstatechange", (e) => console.log(e));

  // stream.getTracks().forEach((track) => pc.addTrack(track, stream));

  // try {
  //   const offer = await pc.createOffer({ offerToReceiveVideo: 1 });
  //   await pc.setLocalDescription(offer);
  //   console.log("ofer", offer);
  //   const input = document.getElementById("clientId");
  //   // Send a message
  //   client.sendMessage(`findClient:${input.value}:${offer.sdp}`);
  // } catch (e) {
  //   console.error(e);
  // }
  // client.addMessageHandler("answer", async (message) => {
  //   console.log("answer", message);
  //   const colonIdx = message.indexOf(":");
  //   await pc.setRemoteDescription({
  //     type: "answer",
  //     sdp: message.slice(colonIdx + 1),
  //   });
  // });
}

async function acceptCalls() {
  const clientId = document.getElementById("clientId").value;
  client.addMessageHandler("remote", async (message) => {
    const pc = new RTCPeerConnection();
    pc.addEventListener("track", gotRemoteStream);
    const colonIdx = message.indexOf(":");
    await pc.setRemoteDescription({
      type: "offer",
      sdp: message.slice(colonIdx + 1),
    });

    const answer = await pc.createAnswer();
    pc.setLocalDescription(answer);
    client.sendMessage(`answer:${clientId}:${answer.sdp}`);
  });
}

document
  .getElementById("stream")
  .addEventListener("loadedmetadata", function () {
    console.log(
      `Remote video videoWidth: ${this.videoWidth}px,  videoHeight: ${this.videoHeight}px`
    );
  });

function gotRemoteStream(e) {
  console.log(e);
  const video = document.getElementById("stream");
  if (video.srcObject !== e.streams[0]) {
    video.srcObject = e.streams[0];
    console.log("pc2 received remote stream");
  }
}

async function connect() {
  const localClientId = document.getElementById("localClientId").value;
  client.connect(localClientId);
}

document.getElementById("connectForm").addEventListener("submit", (e) => {
  e.preventDefault();
  connect();
});

document.getElementById("callForm").addEventListener("submit", (e) => {
  e.preventDefault();
  call();
});
