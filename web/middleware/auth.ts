export default defineNuxtRouteMiddleware((to, from) => {
  const pinia = usePinia()
  const { loggedIn } = useAuthStore(pinia)
  console.log(`logged in: ${loggedIn}`)
  if (!loggedIn) {
    return navigateTo('/auth')
  }
})
