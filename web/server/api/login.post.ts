import { z } from "zod"

const bodySchema = z.object({
	username: z.string(),
	password: z.string()
})

type LoginResponse = {
	session_id: string,
	user: {
		id: string,
		username: string,
		created_at: Date
	},
	access_token: string,
	refresh_token: string,
	access_token_expires_at: Date,
	refresh_token_expores_at: Date,
	error: {
		code: number,
		message: string
	}


}

export default defineEventHandler(async (event) => {
	const { username, password } = await readValidatedBody(event, bodySchema.parse)
	const response = await $fetch<LoginResponse>("http://localhost:8080/v1/login_user", {
		method: 'POST',
		body: {
			username, password
		}
	})



	if (response?.error) {
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

	setCookie(event, 'refresh_token', response.refresh_token, {
		httpOnly: true,
		secure: process.env.NODE_ENV === 'production',
		sameSite: 'strict'
	})

	await setUserSession(event, {
		session: {
			id: response.session_id,
			user: response.user,
			access: response.access_token,
			refresh: response.refresh_token
		},
	})

	return






})
