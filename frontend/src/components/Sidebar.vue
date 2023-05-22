<template>
  <div class="drawer">
    <input :id="drawerId" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <!-- Page content here -->
      <slot name="content">
        {{ console.error("main content must be provided") }}
      </slot>
    </div>
    <div class="drawer-side">
      <label :for="drawerId" class="drawer-overlay"></label>
      <ul class="menu p-4 w-80 bg-base-100 text-base-content">
        <!-- Sidebar content here -->
        <li>
          <!-- toggle could be done by this.$refs.input.checked = false -->
          <label :for="drawerId" @click="currentPage = 'main'"> 主页面 </label>
        </li>
        <li>
          <label :for="drawerId" @click="currentPage = 'p2p'"> P2P 页面 </label>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { currentPage } from "@/store";

defineProps({
  drawerId: { type: String, required: true },
});
defineSlots<{
  content: (props: {}) => {};
}>();
</script>
