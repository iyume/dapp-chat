<template>
  <div class="dropdown">
    <label tabindex="0" class="btn btn-sm btn-primary whitespace-nowrap"
      >后端: 00.00.00.00:00000</label
    >
    <ul
      tabindex="0"
      class="dropdown-content menu p-2 shadow bg-base-100 rounded-box backdrop-blur bg-opacity-50 z-20"
    >
      <li class="rounded">
        <button class="flex items-center gap-x-1 text-sm px-2">
          <MiniCheckIcon />00.00.00.00:00
        </button>
      </li>
      <li class="rounded">
        <button class="flex items-center gap-x-1 text-sm px-2">
          <MiniPlusIcon />添加后端
        </button>
      </li>
    </ul>
  </div>
</template>
<script setup lang="ts">
import MiniCheckIcon from "@/components/icons/MiniCheckIcon.vue";
import MiniPlusIcon from "@/components/icons/MiniPlusIcon.vue";
</script>
