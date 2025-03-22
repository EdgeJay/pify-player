export const spotifyTokenKey = 'spotify_token';

interface WSCommand {
	command: string;
	payload?: {
		[key: string]: string;
	};
}

interface WSResponse {
	command: string;
	body: unknown;
}

interface WSConnectResponse {
	access_token: string;
}

interface ConnectWSOptions {
	onConnect: (token: string) => void;
	onError: (message: string) => void;
}

const onConnectResponse = (payload: WSConnectResponse, onConnect: (token: string) => void) => {
	// Connect to Spotify after receiving token
	console.log(`received access token: ${payload.access_token}`);
	localStorage.setItem(spotifyTokenKey, payload.access_token);

	/**
	 * Establishes a connection to the player service.
	 * This method should be called to initialize and connect the player to the backend.
	 */
	onConnect(payload.access_token);
};

export const getSpotifyTokenFromStorage = (): string => {
	return localStorage.getItem(spotifyTokenKey) || '';
};

export const clearSpotifyTokenFromStorage = () => {
	localStorage.removeItem(spotifyTokenKey);
};

export const sendConnectCommand = (ws: WebSocket) => {
	const command: WSCommand = {
		command: 'connect'
	};
	ws.send(JSON.stringify(command));
};

export const getApiConnectWS = ({ onConnect, onError }: ConnectWSOptions): WebSocket => {
	const domain = window.location.hostname;
	const ws = new WebSocket(`wss://${domain}:8080/api/player/ws`);

	/* When WebSocket connection is successfully established, sets up the onopen event handler
	 * This callback function is executed once the WebSocket connection is opened
	 */
	ws.onopen = () => {
		const token = localStorage.getItem(spotifyTokenKey) || '';
		if (!token) {
			sendConnectCommand(ws);
		} else {
			console.log('access token found in localStorage', token);
			onConnect(token);
		}
	};

	ws.onmessage = async (event) => {
		const response = JSON.parse(event.data) as WSResponse;
		switch (response.command) {
			case 'connect':
				onConnectResponse(response.body as WSConnectResponse, onConnect);
				break;
		}
	};

	ws.onerror = (error) => {
		console.error('WebSocket error:', error);
		onError(`WebSocket error occurred: ${error}`);
	};

	ws.onclose = () => {
		console.log('WebSocket connection closed');
		onError('WebSocket connection closed');
	};

	return ws;
};
