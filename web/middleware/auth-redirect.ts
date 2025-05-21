export default defineNuxtRouteMiddleware(async (to, from) => {
  const pinia = usePinia()
  const auth = useAuthStore(pinia)

  const { loggedIn } = storeToRefs(auth)
  await auth.refresh()

  if (loggedIn.value) {
    return navigateTo('/')
  }
})
