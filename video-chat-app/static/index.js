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
        const message = JSON.parse(event.data);
        console.log("Received message:", message);
        // Handle incoming message
        this.messageHandlers.forEach(([t, fn]) => {
          const { type, from, data } = message;
          if (t === type) {
            fn({ from, data });
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

  // TODO: save clientid for reconnect attempts
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

let localClientId;

let pc;

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

  const videoTracks = stream.getVideoTracks();
  if (videoTracks.length > 0) {
    console.log(`Using video device: ${videoTracks[0].label}`);
  }

  pc.addEventListener("icecandidate", (e) => {
    console.log("ice candidate", e);
    if (e.candidate) {
      client.sendMessage({
        to: remoteClientId,
        from: localClientId,
        message: JSON.stringify({
          type: "new-ice-candidate",
          from: localClientId,
          data: {
            candidate: e.candidate,
          },
        }),
      });
    }
  });
  // pc.addEventListener("iceconnectionstatechange", (e) => console.log(e));

  stream.getTracks().forEach((track) => pc.addTrack(track, stream));

  try {
    const offer = await pc.createOffer({ offerToReceiveVideo: 1 });
    await pc.setLocalDescription(offer);
    // Send a message
    client.sendMessage({
      to: remoteClientId,
      from: localClientId,
      message: JSON.stringify({
        type: "offer",
        from: localClientId,
        data: {
          sdp: offer.sdp,
        },
      }),
    });
  } catch (e) {
    console.error(e);
  }
  client.addMessageHandler("answer", async ({ from, data }) => {
    console.log("answer", from, data);
    await pc.setRemoteDescription({
      type: "answer",
      sdp: data.sdp,
    });
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
  pc = new RTCPeerConnection();
  localClientId = document.getElementById("localClientId").value;
  client.connect(localClientId);
  client.addMessageHandler("offer", async ({ from, data }) => {
    console.log(
      `received message. From: ${from}. Data: ${JSON.stringify(data)}`
    );
    pc.addEventListener("track", gotRemoteStream);
    await pc.setRemoteDescription({
      type: "offer",
      sdp: data.sdp,
    });

    const answer = await pc.createAnswer();
    pc.setLocalDescription(answer);
    console.log("send answer");
    client.sendMessage({
      to: from,
      from: localClientId,
      message: JSON.stringify({
        type: "answer",
        from: localClientId,
        data: {
          sdp: answer.sdp,
        },
      }),
    });
  });

  client.addMessageHandler("new-ice-candidate", async ({ from, data }) => {
    console.log("receive ice candidate", from, data);
    const candidate = new RTCIceCandidate(data.candidate);
    pc.addIceCandidate(candidate);
  });
}

document.getElementById("connectForm").addEventListener("submit", (e) => {
  e.preventDefault();
  connect();
});

document.getElementById("callForm").addEventListener("submit", (e) => {
  e.preventDefault();
  call();
});

document
  .getElementById("startVideoBtn")
  .addEventListener("click", () => startVideo());
