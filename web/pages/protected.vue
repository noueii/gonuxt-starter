<script setup lang="ts">
definePageMeta({
  middleware: ['auth']
})
const isDisabled = computed((): boolean => {
  const { username, password, passwordConfirm } = newUser
  let disabled = true
  if (password.length > 0) {
    disabled = false
  }
  if (username.length > 0 && username !== user.value?.username) {
    disabled = false
  }
  if (passwordConfirm !== password) disabled = true

  return disabled
})

const newUser = reactive({
  username: '',
  password: '',
  passwordConfirm: '',

})

const authStore = useAuthStore()

const { user } = storeToRefs(authStore)


async function handleUpdate() {
  if (isDisabled.value) return

  const { $apiClient } = useNuxtApp()
  const { password, username } = newUser
  const { refresh, user } = useAuthStore()
  const toast = useToast()

  const { response, error } = await $apiClient.PATCH('/v1/user', {
    body: {
      id: user?.id,
      username: username.length > 0 ? username : undefined,
      password: password.length > 0 ? password : undefined
    }
  })

  if (error) {
    toast.add({
      title: "Error updating user",
      description: error.message,
      color: 'error'
    })
    return
  }

  if (response.ok) {
    await refresh()

  }
}



</script>

<template>
  <div class="flex flex-col p-2">
    <div class="flex flex-col">
      <span>This is a protected route. Only authenticated users can access this page.</span>
      <span>Welcome, {{ user?.username }}</span>
    </div>

    <div class="flex flex-col gap-2 w-96">
      <h1 class="text-xl font-bold">User settings</h1>
      <div class="flex flex-col gap-2">
        <div class="flex gap-2 flex-col">
          <label>Username</label>
          <input v-model="newUser.username" class="border-2 rounded p-2" :placeholder="user?.username">

        </div>
        <div class="flex gap-2 flex-col">
          <label>New password</label>
          <input v-model="newUser.password" class="border-2 rounded p-2" type="password">

        </div>
        <div class="flex gap-2 flex-col">
          <label>Confirm password</label>
          <input v-model="newUser.passwordConfirm" class="border-2 rounded p-2" type="password">

        </div>
        <div class="flex gap-2 flex-col w-full">
          <label>Email</label>
          <input class="border-2 rounded p-2" :placeholder="user?.email" disabled>

        </div>

        <Button :disabled="isDisabled" @click="handleUpdate">Update</Button>
      </div>
    </div>
  </div>
</template>
