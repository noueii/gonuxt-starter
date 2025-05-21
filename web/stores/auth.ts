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
  email: string | undefined,
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
    async login(input: { email: string, password: string }): Promise<void> {
      const { email, password } = input
      const { $apiClient } = useNuxtApp()
      const toast = useToast()

      if (!email || !password || email.length === 0 || password.length === 0) return

      const { response, error, data } = await $apiClient.POST('/v1/login_user', {
        body: {
          email, password
        }
      })



      if (error) {
        toast.add({
          title: "Login error",
          description: error.message,
          color: 'error'
        })
        return
      }

      if (response.ok) {
        toast.add({
          title: "Logged in",
          color: "success"
        })
      }

      this.user = {
        id: data.user?.id,
        username: data.user?.username,
        role: data.user?.role,
        email: data.user?.email
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
    async register(input: { email: string, password: string, confirmPassword: string }): Promise<void> {
      const { email, password, confirmPassword } = input
      const { $apiClient } = useNuxtApp()
      const toast = useToast()
      if (!email || !password || email.length === 0 || password.length === 0 || password !== confirmPassword) return

      const { response, error } = await $apiClient.POST('/v1/create_user', {
        body: {
          email, password
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
      const { response, error, data } = await $apiClient.GET("/v1/refresh_token")
      //const toast = useToast()
      if (error) {
        return
      }

      if (response.ok) {

        if (!data || !data.user || !data.session) return

        this.user = {
          id: data.user.id,
          email: data.user.email,
          role: data.user.role,
          username: data.user.username
        }

        this.session = {
          id: data.session.id,
          expiresAt: new Date(data.session.expires_at ?? 0)
        }

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

