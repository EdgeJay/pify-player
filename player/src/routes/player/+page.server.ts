import {
	PLAYER_NAME,
	BASIC_AUTH_USERNAME,
	BASIC_AUTH_PASSWORD,
	ENABLE_YOUTUBE
} from '$env/static/private';

export function load() {
	return {
		playerName: PLAYER_NAME,
		basicAuthToken: Buffer.from(`${BASIC_AUTH_USERNAME}:${BASIC_AUTH_PASSWORD}`).toString('base64'),
		enableYoutube: ENABLE_YOUTUBE === '1'
	};
}
