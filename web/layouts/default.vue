<script setup>
const { loggedIn, logout } = useAuthStore()

async function handleLogOut() {
  await logout()
  reloadNuxtApp()
}
</script>

<template>
  <div class="w-screen flex">
    <main class="flex flex-col w-full gap-2">
      <div class="flex gap-4 w-full items-center justify-center">
        <NuxtLink class="p-2 bg-secondary rounded" to="/"> Home </NuxtLink>
        <NuxtLink class="p-2 bg-secondary rounded" to="/protected"> Protected </NuxtLink>
        <NuxtLink class="p-2 bg-secondary rounded" to="/admin"> Admin </NuxtLink>
        <NuxtLink v-if="!loggedIn" class="p-2 bg-secondary rounded" to="/auth"> Log in </NuxtLink>
        <div v-if="loggedIn" class="p-2 bg-secondary rounded cursor-pointer" @click="handleLogOut">Log out</div>
      </div>
      <slot />
    </main>
  </div>




</template>
