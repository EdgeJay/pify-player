export const spotifyTokenKey = 'spotify_token';
export const spotifyTokenExpiresAtKey = 'spotify_token_expires_at';

export interface LoginResponse {
	logged_in: boolean;
	redirect_url: string;
	user: {
		display_name: string;
		profile_image_url: string;
	};
}

interface AccessTokenInfo {
	accessToken: string;
	expiresAt: Date | null;
}

export const checkSession = async (): Promise<LoginResponse> => {
	// check login status first
	const DOMAIN = window.location.hostname;
	const response = await fetch(`https://${DOMAIN}:8080/api/auth/login`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});

	if (!response.ok) {
		throw new Error('Login failed');
	}

	const res = (await response.json()) as LoginResponse;
	return res;
};

export const getSpotifyTokenFromStorage = (): AccessTokenInfo => {
	const expiresAtStr = localStorage.getItem(spotifyTokenExpiresAtKey);
	const expiresAt = expiresAtStr ? new Date(parseInt(expiresAtStr, 10)) : null;
	return {
		accessToken: localStorage.getItem(spotifyTokenKey) || '',
		expiresAt
	};
};

export const saveSpotifyTokenToStorage = (accessToken: string, expiresAt: Date) => {
	localStorage.setItem(spotifyTokenKey, accessToken);
	localStorage.setItem(spotifyTokenExpiresAtKey, expiresAt.getTime().toString());
};

export const clearSpotifyTokenFromStorage = () => {
	localStorage.removeItem(spotifyTokenKey);
	localStorage.removeItem(spotifyTokenExpiresAtKey);
};
