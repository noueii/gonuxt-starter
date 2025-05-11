import createClient, { type Middleware } from "openapi-fetch"
import type { paths } from "./api"

// TODO: CHANGE TO MATCH ENV VARIABLE

let accessToken: string | undefined = undefined

const UNPROTECTED_ROUTES = ["/create_user"];


const authMiddleware: Middleware = {
    async onRequest({ schemaPath, request }) {
        if (UNPROTECTED_ROUTES.some((pathname) => schemaPath.startsWith(pathname))) return

        if (!accessToken) {
            const authRes = { accessToken: "" }; // TODO: CHANGE TO SOME FUNCTION TO GET ACCESSTOKEN
            if (authRes.accessToken) {
                accessToken = authRes.accessToken;
            } else {
                // TODO: HANDLE AUTH ERROR
            }
        }

        request.headers.set("Authorization", `Bearer ${accessToken}`)
        return request
    }
}




const client = createClient<paths>({ baseUrl: "http://localhost:8080" })
client.use(authMiddleware)

export { client }
