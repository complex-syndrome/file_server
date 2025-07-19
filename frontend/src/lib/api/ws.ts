export function ConnectSocket(
	selfCloseRef: () => boolean,
	refreshFunction: () => void,
	setTimeoutRef: (v: ReturnType<typeof setTimeout> | null) => void
): WebSocket {
	const socket = new WebSocket(import.meta.env.VITE_WS_URL);

	socket.onopen = () => {
		console.log('WebSocket connection established.');
		refreshFunction();
	};

	socket.onmessage = (_) => {
		const t = setTimeout(() => {
			refreshFunction();
			setTimeoutRef(null);
		}, 300);

		setTimeoutRef(t);
	};

	socket.onclose = (_) => {
		if (!selfCloseRef()) {
			console.log('Websocket conn disconnected. Attempting to reconnect...');
			setTimeout(() => {
				ConnectSocket(selfCloseRef, refreshFunction, setTimeoutRef);
			}, 3000);
		}
	};

	socket.onerror = (error) => {
		if (socket.readyState !== WebSocket.CLOSED && socket.readyState !== WebSocket.CLOSING) {
			console.error('WebSocket error:', error);
			socket.close();
		}
	};
	return socket;
}
