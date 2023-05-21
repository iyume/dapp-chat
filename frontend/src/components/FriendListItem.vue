<template>
  <div class="user">
    <div class="avatar">
      <div class="w-8 rounded">
        <img :src="avatar" />
      </div>
    </div>
    <span class="ml15">{{ remark }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{ avatar?: string; remark: string }>()
</script>
