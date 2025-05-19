export default defineNuxtRouteMiddleware(async (to, from) => {
  const { verify, user } = useAuthStore()
  console.log(user)

  await verify()

  console.log(user)


  if (!user || user.role !== "admin") {
    return navigateTo('/forbidden')
  }

})
