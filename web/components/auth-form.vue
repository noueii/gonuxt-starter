<script setup lang="ts">
import { client } from "../src/api/v1/client.ts"
const variant = ref("signin")


function handleVariantChange() {
  console.log("hello")
  if (variant.value == "signin") {
    variant.value = "register"
  } else {
    variant.value = "signin"
  }
}

async function handleSubmit(event) {
  const formData = new FormData(event.target)
  const username = formData.get("username")
  const password = formData.get("password")
  const passwordConfirm = formData.get("password-confirm")

  if (!username || !password || (variant.value === "register" && !passwordConfirm)) {
    return
  }

  if (variant.value === "signin") {
    const res = await client.POST("/login_user", {
      body: {
        username: username,
        password: password
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

    const res = await client.POST("/create_user", {
      body: {
        username, password
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
  <form className='border-2 flex flex-col rounded-2xl p-4 w-full h-full gap-4' @submit.prevent="handleSubmit">
    <h2 v-if='variant === "signin"'>Sign in</h2>
    <h2 v-if='variant === "register"'>Register</h2>
    <div className='flex flex-col'>
      <div className='flex flex-col'>
        <label>Username</label>
        <input className='border-2 border-accent rounded w-full p-1' name="username">
      </div>
      <div className='flex flex-col'>
        <label class>Password</label>
        <input className='border-2 border-accent rounded w-full p-1' name="password">
      </div>
      <div v-if="variant == 'register'" className='flex flex-col'>
        <label class>Confirm Password</label>
        <input className='border-2 border-accent rounded w-full p-1' type="password" name="password-confirm">
      </div>
    </div>
    <span v-if="variant === 'signin'" className="cursor-pointer" @click="handleVariantChange">Don't have an
      account?</span>
    <span v-if="variant === 'register'" className="cursor-pointer" @click="handleVariantChange">Already have an
      account?</span>
    <Button type="submit">Sign in</Button>
  </form>
</template>
