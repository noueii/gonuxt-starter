import { defineStore } from "pinia"

type State = {
  session: Session | undefined,
  user: User | undefined
}

type Session = {
  id: string | undefined,
  expiresAt: Date
}

type User = {
  id: string | undefined,
  username: string | undefined,
  role: string | undefined
}


export const useAuthStore = defineStore('auth', {
  state: (): State => ({
    session: undefined,
    user: undefined
  }),



  getters: {
    loggedIn(): boolean {
      return !!this.session
    },

  },

  actions: {
    async login(input: { username: string, password: string }): Promise<void> {
      const { username, password } = input
      const { $apiClient } = useNuxtApp()
      const toast = useToast()

      console.log(username, password)
      if (!username || !password || username.length === 0 || password.length === 0) return

      const { response, error, data } = await $apiClient.POST('/v1/login_user', {
        body: {
          username, password
        }
      })



      if (error) {
        console.log(error)
        toast.add({
          title: "Login error",
          description: error.message,
          color: 'error'
        })
        return
      }

      if (response.ok) {
        console.log(response)
        toast.add({
          title: "Logged in",
          color: "success"
        })
      }

      this.user = {
        id: data.user?.id,
        username: data.user?.username,
        role: data.user?.role
      }

      this.session = {
        id: data?.session_id,
        expiresAt: new Date()
      }

      reloadNuxtApp({
        path: '/'
      })

    },
    async logout(): Promise<void> {
      const refreshToken = useCookie('refresh_token')
      const accessToken = useCookie('access_token')
      refreshToken.value = undefined
      accessToken.value = undefined
      this.session = undefined
      this.user = undefined

    },
    async register(input: { username: string, password: string, confirmPassword: string }): Promise<void> {
      const { username, password, confirmPassword } = input
      const { $apiClient } = useNuxtApp()
      const toast = useToast()
      if (!username || !password || username.length === 0 || password.length === 0 || password !== confirmPassword) return

      const { response, error } = await $apiClient.POST('/v1/create_user', {
        body: {
          username, password
        }
      })

      if (error) {
        toast.add({
          title: "Registration error",
          description: error.message,
          color: 'error'
        })
        return
      }

      if (response.ok) {
        toast.add({
          title: "Registered successfully",
          color: "success"
        })
      }
      //return await this.login({ username, password })
    },

    async refresh(): Promise<void> {
      const { $apiClient } = useNuxtApp()
      const response = await $apiClient.GET("/v1/refresh_token")
      const toast = useToast()
      if (response.error) {
        toast.add({
          title: "Error refreshing session",
          description: response.error.message,
          color: 'error'
        })
      }

      if (response.response.ok) {
        toast.add({
          title: "Session refreshed",
          color: 'success',
        })
      }
    },
    async verify(): Promise<void> {
      const { $apiClient } = useNuxtApp()
      const { response, data, error } = await $apiClient.GET('/v1/verify_token')

      if (error) {
        this.session = undefined
        this.user = undefined
        return
      }

      if (response.ok) {
        if (!this.user) {
          return
        }
        this.user.id = data.user_id
        this.user.role = data.role
      }
    }
  },
  persist: true
})

