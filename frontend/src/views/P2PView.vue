<template>
  <div class="h-full">
    <div
      class="w-1/4 h-full float-left p-2 drawer-side bg-gray-100 overflow-y-auto"
    >
      <ul class="w-full menu menu-vertical">
        <li v-for="fr in friends" :key="fr.node_id">
          <FriendListItem
            @click="select_chat(fr.node_id)"
            :active="chat_with == fr.node_id"
            :name="fr.remark"
            :avatar="undefined"
          />
        </li>
      </ul>
    </div>
    <div class="w-3/4 h-full float-right p-2 drawer-content overflow-y-auto">
      <div v-if="chat_with != null" class="h-full w-full">
        <div class="h-3/4 w-full overflow-y-auto">
          <div v-for="item in message_list" :key="item.message_id">
            <MessageItem v-bind="item" />
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
      <div v-else class="container h-full w-full">Select a Chat</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { api } from "@/api";
import type { IFriend } from "@/interfaces";
import FriendListItem from "@/components/FriendListItem.vue";
import MessageItem from "@/components/MessageItem.vue";
import type { MessageItemType } from "@/components/MessageItem.vue";

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
