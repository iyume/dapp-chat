<template>
  <div class="flex flex-col h-full">
    <div class="flex-1 overflow-y-scroll px-4">
      <div class="h-0">
        <div class="h-8"></div>
        <h2 class="text-3xl font-extrabold">
          {{ userRemark }}
        </h2>
        <div class="h-2"></div>
        <span class="text-sm min-w-0 truncate opacity-60">
          <!-- FIXME: shrink not work -->
          0x{{ nodeId }}
        </span>
        <div class="h-4"></div>
        <div
          class="badge whitespace-nowrap"
          :class="connBadgeTable[status].badge"
        >
          {{ connBadgeTable[status].label }}
        </div>
        <div class="divider"></div>
        <template v-for="e in selectedSession.events">
          <div
            class="chat"
            :class="e.user_id == selfID ? 'chat-end' : 'chat-start'"
          >
            <div class="chat-header">
              {{ e.user_id == selfID ? "me" : userRemark }}
              <time class="text-xs opacity-50">
                {{ utils.sentTimeChat(e.time_iso) }}
              </time>
            </div>
            <div class="chat-bubble">
              {{ utils.extractPlainText(e.message) }}
            </div>
            <div
              v-if="e.hash in verifiedMessages"
              class="chat-footer text-xs font-medium text-green-600"
            >
              √ Verified
            </div>
            <div v-else class="chat-footer text-xs font-medium text-red-600">
              × Not Verified
            </div>
          </div>
        </template>
      </div>
    </div>
    <div class="flex-none overflow-hidden px-4 py-2 relative">
      <!-- TODO: auto resize textarea to fit content -->
      <textarea
        v-model="inputMessage"
        @keypress.enter="sendMessageByEnter"
        class="textarea h-16 resize-none w-full bg-base-300 no-scrollbar leading-5"
        placeholder="输入消息发送"
      ></textarea>
      <button
        class="btn btn-square btn-primary absolute right-6 top-4"
        @click="sendMessage"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
          />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";

import type { IP2PSession } from "@/interfaces";
import {
  actionGetP2PSession,
  actionSendP2PMessage,
  p2pSessions,
  selfID,
  friendsPeerInfo,
  FriendStatus,
  peersInfo,
  verifiedMessages,
} from "@/store";
import utils from "@/utils";

const props = defineProps({
  nodeId: { type: String, required: true },
});

const selectedSession = computed<IP2PSession>(() => {
  let id = props.nodeId;
  if (id == "") {
    // no chat selected, should not be render
    return { events: [] };
  }
  if (!(props.nodeId in p2pSessions.value)) {
    // task done and responsively updates
    actionGetP2PSession(props.nodeId);
    return { events: [] };
  }
  return p2pSessions.value[id];
});

// undefined indicates that this is anonymous chat
const friendInfo = computed<(typeof friendsPeerInfo.value)[string] | undefined>(
  () => friendsPeerInfo.value[props.nodeId]
);

const connBadgeTable: Record<FriendStatus, { badge: string; label: string }> = {
  [FriendStatus.Connected]: { badge: "badge-success", label: "已连接" },
  [FriendStatus.Disconnected]: { badge: "badge-warning", label: "不活跃" },
  [FriendStatus.Notconnected]: { badge: "badge-error", label: "未连接" },
};

const status = computed(() => {
  if (friendInfo.value != undefined) {
    return friendInfo.value.status;
  }
  if (props.nodeId in peersInfo.value) {
    return peersInfo.value[props.nodeId].active
      ? FriendStatus.Connected
      : FriendStatus.Disconnected;
  }
  return FriendStatus.Notconnected;
});

const userRemark = computed<string>(() => {
  if (friendInfo.value == undefined) {
    return "匿名";
  }
  return friendInfo.value.remark;
});

const inputMessage = ref("");

function sendMessage() {
  const message = inputMessage.value;
  inputMessage.value = "";
  actionSendP2PMessage(props.nodeId, message);
}

function sendMessageByEnter(e: KeyboardEvent) {
  if (!e.getModifierState("Shift")) {
    e.preventDefault();
    sendMessage();
  }
}
</script>
