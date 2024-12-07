export const URI =
  process.env.REACT_APP_SOCKET_API_URI ||
  `${window.location.protocol.startsWith('https') ? 'wss' : 'ws'}://${
    window.location.host
  }`;

const debug = process.env.REACT_APP_DEBUG || null;

interface Opts {
  helloMessage?: string;
  processMessageFn: (data) => Promise<void>;
  onConnectCallback: () => void;
}

export class WebSocketReconnector {
  url: any;
  reconnectInterval: number;
  websocket: WebSocket | null;

  constructor(quizId: string, username: string) {
    this.url = `${URI}/api/quiz/${quizId}?username=${username}`;
    this.reconnectInterval = 3000; // Reconnect interval in milliseconds
    this.websocket = null;
  }

  isReady(): boolean {
    return this.websocket?.readyState === WebSocket.OPEN;
  }

  isNotReady(): boolean {
    return this.websocket?.readyState !== WebSocket.OPEN;
  }

  connect(opts?: Opts) {
    this.websocket = new WebSocket(this.url);

    this.websocket.onopen = () => {
      debug && console.log('WebSocket connection opened');
      if (opts?.helloMessage) {
        this.websocket?.send(opts?.helloMessage);
      }

      opts?.onConnectCallback && opts.onConnectCallback();
    };

    this.websocket.onmessage = (event) => {
      debug && console.log('Received message:', event.data);
      opts?.processMessageFn && opts?.processMessageFn(event.data);
    };

    this.websocket.onclose = (event) => {
      console.error(
        `WebSocket connection closed: ${event.code} - ${event.reason}`,
      );

      this.reconnect();
    };

    this.websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  reconnect() {
    console.log(`Reconnecting in ${this.reconnectInterval / 1000} seconds...`);
    setTimeout(() => {
      this.connect();
    }, this.reconnectInterval);
  }

  send(message) {
    if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
      this.websocket.send(message);
    } else {
      console.warn('WebSocket not open. Message not sent.');
    }
  }

  sendJSON(payload) {
    if (this.websocket?.readyState === WebSocket.OPEN) {
      this.websocket.send(JSON.stringify(payload));
    } else {
      console.warn('WebSocket not open. Message not sent.');
    }
  }

  close() {
    console.log('Sending close message ...');
    this.websocket?.close();
  }
}
