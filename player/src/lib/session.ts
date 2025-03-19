export interface LoginResponse {
	logged_in: boolean;
	redirect_url: string;
	user: {
		display_name: string;
		profile_image_url: string;
	};
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
