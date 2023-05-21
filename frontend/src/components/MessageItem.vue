<template>
  <div :class="`chat ${item.direction == 1 ? 'chat-start' : 'chat-end'}`">
    <div class="chat-image avatar">
      <div class="w-10 rounded-full">
        <img :src="item.avatar" />
      </div>
    </div>
    <div class="chat-header">
      {{ item.name }}
      <time class="text-xs opacity-50">{{ item.sendTime }}</time>
    </div>
    <div class="chat-bubble">{{ item.content }}</div>
    <div class="opacity-50 chat-footer" v-if="item.direction == 2">
      {{ item.seenTime }}
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ item: MessageItemType }>()

export type MessageItemType = {
  message_id: number
  content: string
  sendTime: string
  seenTime?: string
  direction: 1 | 2 // 1: received, 2: sent
  name?: string
  avatar?: string
}
</script>
