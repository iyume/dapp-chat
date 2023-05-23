<template>
  <div :class="active ? 'active' : ''">
    <div class="avatar">
      <div class="w-8 rounded">
        <img src="avatar" alt="empty" />
      </div>
    </div>
    <span class="ml15">{{ name }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  active: boolean;
  name: string;
  avatar?: string;
}>();
</script>
