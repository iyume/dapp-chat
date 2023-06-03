<template>
  <div class="flex flex-col w-full h-full">
    <div class="h-full overflow-y-scroll p-8">
      <h2 class="text-3xl font-extrabold pb-2">{{ userRemark }}</h2>
      <span class="text-sm flex-1 block text-gray-500 min-w-0 shrink truncate">
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
            <!-- FIXME: implement utils.sentTimeChat -->
            <time class="text-xs text-gray-500">{{ e.time }}</time>
          </div>
          <div class="chat-bubble">{{ utils.extractPlainText(e.message) }}</div>
        </div>
      </template>
    </div>
    <div class="flex-none overflow-hidden px-4 py-2">
      <!-- TODO: auto resize textarea to fit content -->
      <textarea
        class="textarea h-16 resize-none w-full bg-base-300 no-scrollbar leading-5"
        placeholder="输入消息发送"
      ></textarea>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import type { IP2PSession } from "@/interfaces";
import {
  actionGetP2PSession,
  p2pSessions,
  selfID,
  friendsPeerInfo,
  FriendStatus,
  peersInfo,
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
  if (id in p2pSessions) {
    return p2pSessions[id].value;
  }
  const newref = ref({ events: [] });
  p2pSessions[id] = newref;
  actionGetP2PSession(id);
  return newref.value;
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
</script>
