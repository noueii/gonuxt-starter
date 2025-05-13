export default defineNuxtRouteMiddleware((to, from) => {
  const { user } = useUserSession()
  if (!user.value?.role) {
    return navigateTo('/forbidden')
  }
})
