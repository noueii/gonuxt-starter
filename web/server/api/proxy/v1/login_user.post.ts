import type { components } from "~/types/api"



export default defineEventHandler(async (event) => {
	const body = await readBody(event)
	const headers = new Headers(event.headers)


	const response = await $fetch<components["schemas"]["pbLoginUserResponse"]>(`http://localhost:8080/${event.path.replace('/api/proxy/', '')}`, {
		method: event.method,
		headers: headers,
		body: ['GET', 'HEAD'].includes(event.method) ? undefined : body,
	})

	await setUserSession(event, {
		user: response.user,
		secure: {
			access_token: response.access_token,
			refresh_token: response.refresh_token
		},
		id: response.session_id
	})

})
