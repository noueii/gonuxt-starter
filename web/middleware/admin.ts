export default defineNuxtRouteMiddleware(async (to, from) => {
  const { verify, user } = useAuthStore()
  await verify()
  if (!user || user.role !== "admin") {
    return navigateTo('/forbidden')
  }

})
