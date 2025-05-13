import type { SecureSessionData, UserSession } from '#auth-utils'
import type { H3Event } from 'h3'
import { FetchError } from 'ofetch'


export default defineEventHandler(async (event) => {
	try {
		return await proxyRequest(event)
	} catch (err) {
		if (err instanceof FetchError) {
			if (err.statusCode === 401 && getCookie(event, 'access_token')) {
				await refreshToken(event)
				return await proxyRequest(event)

			}
		}
	}

})

async function proxyRequest(event: H3Event): Promise {
	const body = await readBody(event)
	const headers = new Headers(event.headers)
	const session = await getUserSession(event) as UserSession & { secure: SecureSessionData }
	const access_token = session.secure.access_token

	if (access_token) {
		headers.set('Authorization', `Bearer ${access_token}`)

	}

	return $fetch(`http://localhost:8080/${event.path.replace('/api/proxy/', '')}`, {
		method: event.method,
		headers: headers,
		body: ['GET', 'HEAD'].includes(event.method) ? undefined : body,
	})
}

type RefreshTokenResponse = {
	access_token: string,
	access_token_expires_at: Date
}

async function refreshToken(event: H3Event) {
	const session = await getUserSession(event) as UserSession & { secure: SecureSessionData }
	const refresh_token = session.secure.refresh_token

	const response = await $fetch<RefreshTokenResponse>(`http://localhost:8080/v1/refresh_token`, {
		method: "POST",
		body: {
			refresh_token: refresh_token
		}
	})

	setUserSession(event, {
		...session,
		secure: {
			...session.secure,
			access_token: response.access_token
		}
	})
}

