

// server/api/proxy.ts
export default defineEventHandler(async (event) => {
	const { url, method, headers, body } = await readBody(event)

	console.log(url, method, headers, body)
	// Add auth (from secure server context — cookies, session, etc)
	/*
	      const token = getCookie(event, 'access_token')
      
	const res = await $fetch(url, {
	  method,
	  headers: {
	    ...headers,
	    Authorization: `Bearer ${token}`
	  },
	  body
	})
      
	return res
	      */
})

