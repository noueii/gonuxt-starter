<script setup lang="ts">



const { session, user, loggedIn } = useUserSession()
const { $apiClient } = useNuxtApp()
const toast = useToast()

console.log(user.value)

const variant = ref("signin")




const authData = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

function handleVariantChange() {
  if (variant.value == "signin") {
    variant.value = "register"
  } else {
    variant.value = "signin"
  }
}

async function handleSubmit() {
  const { username, password, confirmPassword } = authData

  if (!username || !password || (variant.value === "register" && !confirmPassword)) {
    return
  }

  if (variant.value === "signin") {

    const { response, error } = await $apiClient.POST('/v1/login_user', {
      body: {
        username: username,
        password: password
      }
    })



    if (error) {
      toast.add({
        title: 'Server Error:',
        description: error.message,
        color: 'error'
      })
      return
    }

    if (response.ok) {
      //reloadNuxtApp({ path: '/' })
    }








  }




}

async function handleRefresh() {
  $apiClient.GET('/v1/refresh_token')
}

async function handleUpdate() {

  console.log('HELLO')
  const { response, error } = await $apiClient.PATCH("/v1/update_user", {
    body: {
      balance: 40,
      username: 'nxshappy',
      password: 'razielsvenom'
    }
  })

  console.log(error)
  console.log(response)

  if (error) {
    toast.add({
      title: 'Server Error:',
      description: error.message,
      color: 'error'
    })

  }

  if (response.ok) {
    console.log('OK RESPONSE')
    toast.add({
      id: 'user-data',
      title: 'User updated',
      description: 'Successfully updated user data',
      color: 'success'
    })
  }





}
</script>

<template>
  <form class='border-2 flex flex-col rounded-2xl p-4 w-full h-full gap-4' @submit.prevent="handleSubmit">
    <h1 v-if='loggedIn'>{{ user?.username }}</h1>
    <h2 v-if='variant === "signin"'> Sign in </h2>
    <h2 v-if='variant === "register"'> Register</h2>
    <div class='flex flex-col gap-4'>
      <div class='flex flex-col'>
        <label>Username </label>
        <input v-model="authData.username" class='border-2 border-accent rounded w-full p-1' name="username">
      </div>
      <div class='flex flex-col'>
        <label class> Password </label>
        <input v-model="authData.password" class='border-2 border-accent rounded w-full p-1' type="password"
          name="password">
      </div>
      <div v-if="variant == 'register'" class='flex flex-col'>
        <label class> Confirm Password </label>
        <input v-model="authData.confirmPassword" class='border-2 border-accent rounded w-full p-1' type="password"
          name="password-confirm">
      </div>
    </div>
    <span v-if="variant === 'signin'" class="cursor-pointer" @click="handleVariantChange">
      Don't have an account?
    </span>
    <span v-if="variant === 'register'" class="cursor-pointer" @click="handleVariantChange">
      Already have an account ?
    </span>
    <Button type="submit"> Sign in </Button>
    <Button @click="handleUpdate"> Update </Button>
    <Button @click="handleRefresh"> Refresh </Button>
  </form>
</template>
