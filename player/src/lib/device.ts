export interface Device {
	id: string;
	is_active: boolean;
	is_private_session: boolean;
	is_restricted: boolean;
	name: string;
	type: string;
	volume_percent: number;
	supports_volume: boolean;
}

export interface AllDevicesResponse {
	data: {
		devices: Device[];
	};
	error_code: string;
}

export interface ControlPlaybackErrorResponse {
	error_code: string;
}

export const getAllDevices = async (): Promise<AllDevicesResponse> => {
	const domain = window.location.hostname;
	const response = await fetch(`https://${domain}:8080/api/device/all`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});

	if (!response.ok) {
		throw new Error('Get all devices failed');
	}

	const res = (await response.json()) as AllDevicesResponse;
	return res;
};

/**
 * Takes over control of playback for user account.
 * As this endpoint is designed to be called from player page, user auth is not required.
 * @returns
 */
export const controlPlayback = async (accessToken: string, deviceId: string): Promise<boolean> => {
	const domain = window.location.hostname;
	const response = await fetch(`https://${domain}:8080/api/device/control-playback`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			access_token: accessToken,
			device_id: deviceId
		})
	});

	if (!response.ok) {
		throw new Error('control playback failed');
	}

	const success = response.status === 204;
	if (!success) {
		const { error_code } = (await response.json()) as ControlPlaybackErrorResponse;
		console.log(error_code);
	}

	return success;
};
