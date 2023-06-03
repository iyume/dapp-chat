<template>
  <div class="h-full flex flex-row bg-base-100">
    <!-- Friend list -->
    <div
      class="w-80 flex-none bg-base-100 overflow-x-hidden no-scrollbar shadow-sm"
    >
      <div class="h-4"></div>
      <div class="px-4">
        <div class="dropdown dropdown-bottom">
          <label tabindex="0" class="btn btn-sm btn-primary whitespace-nowrap"
            >后端: 00.00.00.00:00000</label
          >
          <ul
            tabindex="0"
            class="dropdown-content shadow bg-base-200 rounded-box p-3 backdrop-blur bg-opacity-60 menu"
          >
            <li class="rounded">
              <button
                class="flex items-center gap-x-1 text-gray-700 text-sm px-2"
              >
                <MiniCheckIcon />00.00.00.00:00
              </button>
            </li>
            <li class="rounded">
              <button
                class="flex items-center gap-x-1 text-gray-700 text-sm px-2"
              >
                <MiniPlusIcon />添加后端
              </button>
            </li>
          </ul>
        </div>
        <div class="w-full flex gap-x-2 py-3">
          <div class="badge badge-success whitespace-nowrap">
            活跃连接: {{ connInfo.stats.connected }}
          </div>
          <div class="badge badge-info whitespace-nowrap">
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
        <li v-for="f in connInfo.friends" class="w-full">
          <div
            class="flex gap-x-4 py-0.5 rounded w-full"
            :class="{ 'bg-base-300': f.node_id == selectedNodeID }"
            @click="selectedNodeID = f.node_id"
          >
            <div
              class="flex-none avatar placeholder"
              :class="cssAvatarStatusTable[f.status]"
            >
              <div class="bg-neutral-focus text-white rounded-full w-8">
                <span class="text-xs">{{ firstChar(f.remark) }}</span>
              </div>
            </div>
            <div class="flex-1 min-w-0 pb-1">
              <p class="text-base font-medium truncate text-gray-700">
                {{ f.remark }}
                <span class="text-xs font-light text-gray-500">{{
                  f.remote_addr
                }}</span>
              </p>
              <p class="text-xs font-light truncate text-gray-500">
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
        <li v-for="p in connInfo.peers" class="w-full">
          <div
            class="flex gap-x-4 py-0.5 rounded w-full"
            :class="{ 'bg-base-300': p.node_id == selectedNodeID }"
            @click="selectedNodeID = p.node_id"
          >
            <div class="flex-1 min-w-0">
              <p class="text-xs font-normal truncate text-gray-700">
                {{ p.remote_addr }}
                <!-- replace with icon? -->
                <span class="font-light text-gray-500"
                  >({{ p.active ? "active" : "inactive" }})</span
                >
              </p>
              <p class="text-xs font-light truncate text-gray-500">
                0x{{ p.node_id }}
              </p>
            </div>
          </div>
        </li>
      </ul>
      <div class="h-4"></div>
    </div>
    <!-- Chat panel, add friend, etc. -->
    <div class="container flex flex-col h-full">
      <!-- TODO: exit button -->
      <div class="h-full">
        <!-- TODO: add friend search page (search by pubkey, node id, remote addr, etc.) -->
        <Messager v-if="selectedNodeID != ''" :node-id="selectedNodeID" />
      </div>
      <div class="flex-none h-4"></div>
    </div>
    <div class="w-96 flex-none">
      <!-- stats here? peer acitvity here? -->
    </div>
  </div>
</template>

<script setup lang="ts">
import Messager from "@/components/Messager.vue";
import MiniPlusIcon from "@/components/icons/MiniPlusIcon.vue";
import MiniCheckIcon from "@/components/icons/MiniCheckIcon.vue";
import { friendsPeerInfo, peersInfo, FriendStatus } from "@/store";
import { computed, ref } from "vue";
import type { IPeerInfo } from "@/interfaces";

const selectedNodeID = ref("");

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

<style scoped></style>
