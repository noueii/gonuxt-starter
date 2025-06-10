export default defineNuxtPlugin(async () => {
  const pinia = usePinia()
  const auth = useAuthStore(pinia)
  await auth.fetchSession()
})
