<template>
  <div class="container h-full">
    <div class="w-1/4 h-full float-left p-2 drawer-side bg-gray-100 overflow-y-auto">
      <ul class="w-full menu menu-vertical">
        <li v-for="fr in friends" :key="fr.id">
          <friend-list-item
            @click="chat_with = fr.id"
            :class="chat_with == fr.id ? 'active' : ''"
            :remark="fr.remark"
            :avatar="fr.avatar"
          />
        </li>
      </ul>
    </div>
    <div class="w-3/4 h-full float-right p-2 drawer-content overflow-y-auto">
      <div v-if="chat_with != null" class="h-full w-full">
        <div class="h-3/4 w-full overflow-y-auto">
          <div v-for="msg in message_list" :key="msg.message_id">
            <message-item :item="msg" />
          </div>
        </div>
        <div class="container h-1/4 w-full form-control">
          <textarea
            class="textarea textarea-bordered h-full w-full"
            placeholder="Type your messaeg here"
            style="resize: none"
            v-model="text"
          ></textarea>
          <label class="label">
            <span></span>
            <span>
              <button @click="send_message" class="btn btn-sm btn-primary mx-1">Send</button>
            </span>
          </label>
        </div>
      </div>
      <div v-else class="container h-full w-full">Select a Chat</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDataStore } from '@/stores/data'
import { ref } from 'vue'
import FriendListItem from '@/components/FriendListItem.vue'
import { computed } from 'vue'
import type { MessageItemType } from '@/components/MessageItem.vue'
import MessageItem from '@/components/MessageItem.vue'

const chat_with = ref<number | undefined>()
const text = ref<string>('')

function send_message() {
  console.log(text.value)
  text.value = ''
}

const message_list = computed<MessageItemType[]>(() => {
  return messages
    .filter((val) => val.friend_id == chat_with.value)
    .map((val) => {
      let name = undefined
      let avatar = undefined
      if (val.direction == 1) {
        let fr = friends.find((f) => f.id == val.friend_id)
        name = fr?.remark || val.friend_id.toString()
        avatar = fr?.avatar
      } else {
        ;(name = my_info?.name), (avatar = my_info?.avatar)
      }
      return {
        message_id: val.message_id,
        content: val.message,
        sendTime: val.time,
        seenTime: 'not implemented',
        direction: val.direction,
        name,
        avatar
      }
    }) as MessageItemType[]
})

const { messages, friends, my_info } = useDataStore()
</script>
