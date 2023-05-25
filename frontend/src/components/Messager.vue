<template>
  <div>
    <p>Messaging</p>
    <div v-if="chat_with != null" class="h-full w-full">
      <div class="h-3/4 w-full overflow-y-auto">
        <div v-for="item in message_list" :key="item.message_id">
          <div
            :class="`chat ${item.direction == 1 ? 'chat-start' : 'chat-end'}`"
          >
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
            <button @click="send_message" class="btn btn-sm btn-primary mx-1">
              Send
            </button>
          </span>
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { api } from "@/api";
import type { IFriend } from "@/interfaces";

export type MessageItemType = {
  message_id: number;
  content: string;
  sendTime: string;
  seenTime?: string;
  direction: 1 | 2; // 1: received, 2: sent
  name?: string;
  avatar?: string;
};

const friends = ref<IFriend[] | undefined>();
const message_list = ref<MessageItemType[] | undefined>();

const chat_with = ref<string | undefined>();
const text = ref<string>("");

function select_chat(node_id: string) {
  chat_with.value = node_id;
  api.getP2PMessageList(node_id).then((resp) => {
    message_list.value = resp.data.map((val) => {
      let name = undefined;
      let direction = 1;
      if (Math.random() < 0.5) {
        let fr = friends.value?.find((f) => f.node_id == val.node_id);
        name = fr?.remark || val.node_id.toString();
        direction = 1;
      } else {
        name = "Me";
        direction = 2;
      }
      return {
        message_id: val.message_id,
        content: val.message,
        sendTime: val.time,
        seenTime: "not implemented",
        direction,
        name,
        undefined,
      } as MessageItemType;
    });
  });
}

function send_message() {
  console.log(text.value);
  text.value = "";
}

api.getFriendList().then((resp) => {
  friends.value = resp.data;
});
</script>
