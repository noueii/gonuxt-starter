<script setup lang="ts">

import { client } from '~/src/api/v1/client'
const variant = ref("signin")

const authData = ref({
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
  const username = authData.value.username
  const password = authData.value.password
  const passwordConfirm = authData.value.confirmPassword

  console.log(`username: ${username}, password: ${password}, passwordConfirm: ${passwordConfirm}`)

  if (!username || !password || (variant.value === "register" && !passwordConfirm)) {
    return
  }

  if (variant.value === "signin") {
    const res = await client.POST("/v1/login_user", {
      body: {
        username: username.toString(),
        password: password.toString()
      }
    })

    console.log(res)

    if (res?.error) {
      alert(res.error.message)
      return
    }


    const { data } = res

    console.log(data)
  }

  if (variant.value === "register") {
    if (password !== passwordConfirm) {
      alert("Passwords do not match")
      return
    }

    const res = await client.POST("/v1/create_user", {
      body: {
        username: username.toString(),
        password: password.toString()
      }
    })

    if (res?.error) {
      alert(res.error.message)
      return
    }

    const { data } = res
    console.log(data)
  }



}
</script>

<template>
  <form class='border-2 flex flex-col rounded-2xl p-4 w-full h-full gap-4' @submit.prevent="handleSubmit">
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
  </form>
</template>
