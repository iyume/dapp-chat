<template>
  <div class="navbar bg-base-300">
    <div class="flex-none">
      <label :for="drawerId" class="btn btn-square btn-ghost">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
          />
        </svg>
      </label>
    </div>
    <div class="flex-1">
      <a class="btn btn-ghost normal-case text-xl" href="/">Awesome Chat App</a>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    drawerId: { type: String, required: true },
  },
};
</script>
