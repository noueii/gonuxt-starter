import createClient from "openapi-fetch";
import type { paths } from "~/types/api";


export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig() as unknown as {
    public: {
      apiBase: string
    }
  }
  const apiClient = createClient<paths>({ baseUrl: config.public.apiBase || 'http://localhost:8080', credentials: 'include' })
  return {
    provide: {
      apiClient: apiClient
    }
  }
})


