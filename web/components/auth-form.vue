<script setup lang="ts">



const { $apiClient } = useNuxtApp()
const toast = useToast()

const { refresh, loggedIn, user, login, register } = useAuthStore()

const variant = ref("signin")

async function handleGoogle() {
  const { data, error } = await $apiClient.GET('/auth/google')
  if (error) {
    return
  }

  if (data.redirect_url) {
    navigateTo(data.redirect_url, { external: true })
  }


}


const authData = reactive({
  email: '',
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
  const { email, password, confirmPassword } = authData

  if (!email || email.length == 0 || !password || password.length == 0 || (variant.value === "register" && !confirmPassword)) {
    return
  }

  if (variant.value === "signin") {
    await login({
      email, password
    })

    console.log('login done')
  }

  if (variant.value === "register") {
    register({ email, password, confirmPassword })
  }
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
        <label>Email </label>
        <input v-model="authData.email" class='border-2 border-accent rounded w-full p-1' name="username">
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
    <Button @click="refresh"> Refresh </Button>
    <Button @click="handleGoogle">Google</Button>
  </form>
</template>
