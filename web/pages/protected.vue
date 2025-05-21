<script setup lang="ts">
definePageMeta({
  middleware: ['auth']
})

const { user } = useAuthStore()

</script>

<template>
  <div class="flex flex-col p-2">
    <div class="flex flex-col">
      <span>This is a protected route. Only authenticated users can access this page.</span>
      <span>Welcome, {{ user?.username }}</span>
    </div>

    <div class="flex flex-col">
      <h1 class="text-xl font-bold">User settings</h1>
      <div class="flex flex-col">
        <div class="flex gap-2">

          <label>Username</label>

        </div>
        <label>Password</label>

      </div>
    </div>
  </div>
</template>
