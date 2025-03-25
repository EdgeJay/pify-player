import { PLAYER_NAME, BASIC_AUTH_USERNAME, BASIC_AUTH_PASSWORD } from '$env/static/private';

export function load() {
	return {
		playerName: PLAYER_NAME,
		basicAuthToken: Buffer.from(`${BASIC_AUTH_USERNAME}:${BASIC_AUTH_PASSWORD}`).toString('base64')
	};
}
