import createClient, { type Middleware } from "openapi-fetch"
import type { paths } from '../types/api'


// TODO: CHANGE TO MATCH ENV VARIABLE





let initialRequest: Request | null = null

const UNPROTECTED_ROUTES = ["/create_user"];

type RefreshTokenResponse = {
  access_token: string,
  access_token_expires_at: string
}


const authMiddleware: Middleware = {
  async onRequest({ schemaPath, request }) {
    if (UNPROTECTED_ROUTES.some((pathname) => schemaPath.startsWith(pathname))) return request

    console.log("Hello from server")

    const accessToken: string | undefined | null = useCookie('access_token')?.value
    console.log("access_token: " + accessToken)
    if (accessToken) {
      request.headers.set("Authorization", `Bearer ${accessToken}`)
    }

    initialRequest = request.clone()

    return request
  },
  async onResponse({ response }) {
    if (response.status === 401) { // add message check 
      const refreshToken = useCookie('refresh_token')?.value


      if (!refreshToken) {
        return response
      }

      await $fetch("/api/token/refresh")


      const accessToken = useCookie('access_token')

      initialRequest?.headers.set("Authorization", `Bearer ${accessToken}`)
      return initialRequest ? fetch(initialRequest) : response
    }

  },
}
const apiClient = createClient<paths>({ baseUrl: "http://localhost:8080" })
apiClient.use(authMiddleware)

export default apiClient






