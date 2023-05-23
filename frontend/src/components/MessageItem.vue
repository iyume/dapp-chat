<template>
  <div :class="`chat ${direction == 1 ? 'chat-start' : 'chat-end'}`">
    <div class="chat-image avatar">
      <div class="w-10 rounded-full">
        <img :src="avatar" alt="empty" />
      </div>
    </div>
    <div class="chat-header">
      {{ name }}
      <time class="text-xs opacity-50">{{ sendTime }}</time>
    </div>
    <div class="chat-bubble">{{ content }}</div>
    <div class="opacity-50 chat-footer" v-if="direction == 2">
      {{ seenTime }}
    </div>
  </div>
</template>

<script setup lang="ts">
export type MessageItemType = {
  message_id: number;
  name: string;
  avatar?: string;
  direction: 1 | 2; // 1: received, 2: sent
  sendTime: string;
  seenTime?: string;
  content: string;
};

defineProps<MessageItemType>();
</script>
