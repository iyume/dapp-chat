<template>
  <div class="h-full flex flex-row bg-base-100">
    <!-- Friend list -->
    <div
      class="w-80 flex-none bg-base-100 overflow-x-hidden no-scrollbar shadow-sm"
    >
      <div class="h-4"></div>
      <div class="px-4">
        <div class="flex flex-wrap gap-2">
          <BackendSelector />
          <div class="w-full"></div>
          <!-- TODO: add tooltip on ID -->
          <!-- See also: https://github.com/saadeghi/daisyui/issues/1899 -->
          <div class="badge badge-info whitespace-nowrap cursor-pointer">
            ID: {{ selfIDText }}
          </div>
          <div class="badge badge-success whitespace-nowrap cursor-pointer">
            活跃连接: {{ connInfo.stats.connected }}
          </div>
          <div class="badge badge-info whitespace-nowrap cursor-pointer">
            好友数量: {{ connInfo.stats.friendCount }}
          </div>
        </div>
      </div>
      <ul class="menu menu-compact menu-vertical px-4">
        <li class="menu-title">
          <div class="flex items-center justify-between">
            <span class="inline-block">好友节点</span>
            <button class="btn btn-ghost btn-xs">
              <!-- TODO: add button action -->
              <MiniPlusIcon />
            </button>
          </div>
        </li>
        <li
          v-for="f in connInfo.friends"
          class="w-full hover:bg-base-200 rounded cursor-pointer"
        >
          <div
            class="flex gap-x-4 py-0.5 rounded w-full"
            :class="{ 'bg-base-300': f.node_id == selectedNodeID }"
            @click="selectNodeID(f.node_id)"
          >
            <div
              class="flex-none avatar placeholder"
              :class="cssAvatarStatusTable[f.status]"
            >
              <div
                class="bg-neutral-focus text-neutral-content rounded-full w-8"
              >
                <span class="text-xs">{{ firstChar(f.remark) }}</span>
              </div>
            </div>
            <div class="flex-1 min-w-0 pb-1">
              <p class="text-base font-medium truncate">
                {{ f.remark }}
                <span class="text-xs font-light opacity-60">
                  {{ f.remote_addr }}
                </span>
              </p>
              <p class="text-xs font-light opacity-60 truncate">
                0x{{ f.node_id }}
              </p>
            </div>
          </div>
        </li>
      </ul>
      <ul class="menu menu-compact menu-vertical px-4">
        <li></li>
        <li class="menu-title">
          <span>节点列表</span>
        </li>
        <li
          v-for="p in connInfo.peers"
          class="w-full hover:bg-base-200 rounded cursor-pointer"
        >
          <div
            class="flex gap-x-4 py-0.5 rounded w-full"
            :class="{ 'bg-base-300': p.node_id == selectedNodeID }"
            @click="selectNodeID(p.node_id)"
          >
            <div class="flex-1 min-w-0">
              <p class="text-xs font-normal truncate">
                {{ p.remote_addr }}
                <!-- replace with icon? -->
                <span class="font-light opacity-60"
                  >({{ p.active ? "active" : "inactive" }})</span
                >
              </p>
              <p class="text-xs font-light opacity-60 truncate">
                0x{{ p.node_id }}
              </p>
            </div>
          </div>
        </li>
      </ul>
      <div class="h-4"></div>
    </div>
    <!-- Chat panel, add friend, etc. -->
    <div class="container flex flex-col">
      <!-- exit button? -->
      <div class="h-full px-2">
        <!-- TODO: add friend search page (search by pubkey, node id, remote addr, etc.) -->
        <div v-if="p2pStage != null" class="pt-4">
          <div class="btn btn-circle" @click="resetStage">
            <ArrowLeftIcon />
          </div>
          <AddBackend
            v-if="p2pStage == 'add_backend'"
            :exit="resetStage"
          ></AddBackend>
          <div v-else-if="p2pStage == 'add_friend'"></div>
        </div>
        <Messager v-else-if="selectedNodeID != ''" :node-id="selectedNodeID" />
        <div v-else>选择节点以发送消息</div>
      </div>
      <div class="flex-none h-4"></div>
    </div>
    <div class="w-96 flex-none">
      <!-- stats here? peer acitvity here? -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";

import Messager from "./Messager.vue";
import MiniPlusIcon from "./icons/MiniPlusIcon.vue";
import BackendSelector from "./BackendSelector.vue";
import AddBackend from "./AddBackend.vue";
import ArrowLeftIcon from "./icons/ArrowLeftIcon.vue";

import type { IPeerInfo } from "@/interfaces";
import {
  friendsPeerInfo,
  peersInfo,
  FriendStatus,
  p2pStage,
  selfID,
} from "@/store";

function resetStage() {
  p2pStage.value = null;
}

const selfIDText = computed(() => {
  if (selfID.value == "") {
    return "无法连接后端以获取 ID";
  }
  return selfID.value.slice(0, 7) + "...";
});

const selectedNodeID = ref("");

function selectNodeID(nodeID: string) {
  selectedNodeID.value = nodeID;
  p2pStage.value = null;
}

const cssAvatarStatusTable: Record<FriendStatus, string> = {
  [FriendStatus.Connected]: "online",
  [FriendStatus.Disconnected]: "offline",
  [FriendStatus.Notconnected]: "",
};

const connInfo = computed(() => {
  const friendsPeerInfo_ = friendsPeerInfo.value;
  const peersInfo_ = peersInfo.value;
  // the peers list with friends removed
  const resPeersInfo: IPeerInfo[] = [];
  const resFriends = Object.values(friendsPeerInfo_);
  for (let p of Object.values(peersInfo_)) {
    if (!(p.node_id in friendsPeerInfo_)) {
      resPeersInfo.push(p);
    }
  }
  resFriends.sort((a, b) => {
    if (a.status == b.status) {
      // NOTE: option sensitivity behaves strange
      return a.remark.toLowerCase().localeCompare(b.remark.toLowerCase());
    }
    return a.status - b.status;
  });
  resPeersInfo.sort((a, b) => {
    return Number(b.active) - Number(a.active);
  });
  const stats = {
    connected: 0,
    friendCount: resFriends.length,
  };
  for (let p of Object.values(peersInfo_)) {
    if (p.active) {
      stats.connected += 1;
    }
  }
  return { peers: resPeersInfo, friends: resFriends, stats };
});

const firstChar = (remark: string) => (remark ? remark[0].toUpperCase() : "?");
</script>
