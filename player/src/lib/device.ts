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

export const getAllDevices = async (): Promise<AllDevicesResponse> => {
	const DOMAIN = window.location.hostname;
	const response = await fetch(`https://${DOMAIN}:8080/api/device/all`, {
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
