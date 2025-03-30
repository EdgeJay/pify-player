import { google } from 'googleapis';
import { json } from '@sveltejs/kit';
import { YOUTUBE_API_KEY, VITE_DOMAIN } from '$env/static/private';

const youtube = google.youtube({
	version: 'v3',
	auth: YOUTUBE_API_KEY
});

async function searchVideos(query: string, maxResults = 1) {
	try {
		const response = await youtube.search.list(
			{
				part: ['snippet'],
				q: query,
				type: ['video'],
				maxResults: maxResults
			},
			{
				headers: {
					Referer: `https://${VITE_DOMAIN}`
				}
			}
		);

		/*
		response.data.items?.forEach((item) => {
			const { title, description, publishedAt } = item.snippet || {};
			const videoId = item?.id?.videoId;
			console.log(`Title: ${title}`);
			console.log(`Video ID: ${videoId}`);
			console.log(`Published At: ${publishedAt}`);
			console.log(`Description: ${description?.substring(0, 100)}...`);
			console.log(`URL: https://www.youtube.com/watch?v=${videoId}`);
			console.log('---------------------------------------------');
		});
		*/
		return response.data;
	} catch (error) {
		console.error('Error fetching videos:', error);
	}
}

export async function GET({ url }) {
	const query = url.searchParams.get('query');
	const maxResults = url.searchParams.get('maxResults') || 1;

	console.log('Query:', query);

	if (!query) {
		return json({ error: 'Query parameter is required' }, { status: 400 });
	}

	const response = await searchVideos(query, Number(maxResults));

	return json(response);
}
