interface TrackResponse {
	external_urls: {
		[key: string]: string;
	};
}

interface YoutubeVideoResponse {
	data: {
		video_id: string;
	};
}

export const getTrack = async (basicAuthToken: string, trackId: string) => {
	const DOMAIN = window.location.hostname;
	const response = await fetch(`https://${DOMAIN}:8080/api/player/track/${trackId}`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Basic ${basicAuthToken}`
		}
	});

	if (!response.ok) {
		throw new Error('Get track failed');
	}

	const res = (await response.json()) as TrackResponse;
	return res;
};

export const getAndSaveYoutubeVideo = async (
	basicAuthToken: string,
	query: string,
	trackId: string
) => {
	const DOMAIN = window.location.hostname;
	const response = await fetch(`https://${DOMAIN}:8080/api/player/youtube`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Basic ${basicAuthToken}`
		},
		body: JSON.stringify({
			query,
			spotify_track_id: trackId,
			cache_results: true
		})
	});

	if (!response.ok) {
		throw new Error('Get track failed');
	}

	const res = (await response.json()) as YoutubeVideoResponse;
	return res;
};
