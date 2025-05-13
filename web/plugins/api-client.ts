import createClient from "openapi-fetch";
import type { paths } from "~/types/api";


export default defineNuxtPlugin(() => {
  const apiClient = createClient<paths>({ baseUrl: "/api/proxy" })

  return {
    provide: {
      apiClient: apiClient
    }
  }
})


