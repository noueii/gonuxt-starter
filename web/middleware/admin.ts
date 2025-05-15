export default defineNuxtRouteMiddleware((to, from) => {
  const { user } = useUserSession()
  if (!user.value?.role || user.value?.role !== 'admin') {
    return navigateTo('/forbidden')
  }
})
