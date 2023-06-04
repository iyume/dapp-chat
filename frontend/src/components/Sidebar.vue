<template>
  <div class="drawer overflow-x-auto">
    <input :id="drawerId" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <!-- Page content here -->
      <slot name="content">
        {{ console.error("main content must be provided") }}
      </slot>
    </div>
    <!-- NOTE: DaisyUI v3 use fixed 0 0 for drawer-side which cause element overlap. It is not easy to fix. -->
    <!-- Animation bug ref: https://github.com/saadeghi/daisyui/issues/1888 -->
    <div class="drawer-side top-16 z-10">
      <label :for="drawerId" class="drawer-overlay"></label>
      <ul class="menu p-4 w-80 h-full bg-base-100 text-base">
        <!-- Sidebar content here -->
        <li>
          <label :for="drawerId" @click="currentPage = 'main'" class="py-3"
            >主页面</label
          >
        </li>
        <li>
          <label :for="drawerId" @click="currentPage = 'other'" class="py-3"
            >其他页面</label
          >
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
