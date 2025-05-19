<script setup lang="ts">
definePageMeta({
  middleware: ['auth']
})

const { user } = useAuthStore()

</script>

<template>
  <div class="flex flex-col">
    <span>Welcome, {{ user?.username }}</span>
    <span>This is a protected route. Only authenticated users can access this page.</span>
  </div>
</template>
