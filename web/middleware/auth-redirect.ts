export default defineNuxtRouteMiddleware(async (to, from) => {
  const pinia = usePinia()
  const { refresh } = useAuthStore(pinia)
  await refresh()

  const { loggedIn } = useAuthStore(pinia)
  if (loggedIn) {
    return navigateTo('/')
  }
})
