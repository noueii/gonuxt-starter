export default defineEventHandler(async (event) => {
	return await clearUserSession(event)

})
