export default defineNuxtRouteMiddleware((to, from) => {
  const { loggedIn } = useUserSession()
  console.log(`logged in: ${loggedIn.value}`)
  if (!loggedIn.value) {
    return navigateTo('/auth')
  }
})
