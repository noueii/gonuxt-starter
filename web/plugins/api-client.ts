import createClient from "openapi-fetch";
import type { paths } from "~/types/api";


export default defineNuxtPlugin(() => {
  const apiClient = createClient<paths>({ baseUrl: "http://localhost:8080", credentials: 'include' })
  return {
    provide: {
      apiClient: apiClient
    }
  }
})


