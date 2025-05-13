type Response = {
	access_token: string,
	access_token_expires_at: Date,
	error: {
		code: number,
		message: string
	}
}

export default defineEventHandler(async (event) => {
	const refreshToken = getCookie(event, 'refresh_token')

	if (!refreshToken) {
		throw createError({
			statusCode: 400,
			message: "Bad request. Missing refresh token"
		})
	}

	const response = await $fetch<Response>("http://localhost:8080/v1/refresh_token", {
		method: "POST",
		body: {
			refresh_token: refreshToken
		}
	})

	if (response.error) {
		throw createError({
			statusCode: response.error.code,
			message: response.error.message
		})
	}

	setCookie(event, 'access_token', response.access_token, {
		httpOnly: true,
		secure: process.env.NODE_ENV === 'production',
		sameSite: 'strict'
	})



})
