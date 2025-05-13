<script setup lang="ts">
const auth = useCookie('auth')

console.log(auth.value)


</script>

<template>
  <div class='flex w-screen h-screen items-center justify-center'>
    <WelcomePage />

  </div>
</template>
