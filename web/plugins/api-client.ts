import createClient, { type Middleware } from "openapi-fetch";
import type { paths } from "~/types/api";


export default defineNuxtPlugin(() => {
  const apiClient = createClient<paths>({ baseUrl: "http://192.168.73.124:8080", credentials: 'include' })




  return {
    provide: {
      apiClient: apiClient
    }
  }
})


