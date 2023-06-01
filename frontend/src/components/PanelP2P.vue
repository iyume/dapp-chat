<template>
  <div class="h-full flex flex-row bg-base-100">
    <!-- Friend list -->
    <div class="w-80 bg-base-100 overflow-hidden">
      <div class="h-4"></div>
      <!-- TODO: backend select -->
      <ul class="menu menu-compact menu-vertical px-4">
        <li class="menu-title">
          <div class="flex items-center justify-between">
            <span class="inline-block">好友节点</span>
            <button class="btn btn-ghost btn-xs">
              <SvgSmallPlus />
            </button>
          </div>
        </li>
        <!-- TODO: refactor friends struct -->
        <li v-for="friend in friends" class="w-full">
          <div class="flex gap-x-4 py-0.5 rounded w-full">
            <div class="flex-none avatar online placeholder">
              <div class="bg-neutral-focus text-white rounded-full w-8">
                <span class="text-xs">{{ firstChar(friend.remark) }}</span>
              </div>
            </div>
            <div class="flex-1 min-w-0 pb-1">
              <p class="text-base font-medium">{{ friend.remark }}</p>
              <p class="text-xs font-light truncate text-gray-500">
                0x{{ friend.node_id }}
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
        <!-- TODO: refactor style -->
        <li v-for="peer in peersInfo" class="w-full">
          <div class="flex gap-x-4 py-0.5 rounded w-full">
            <div class="flex-none avatar placeholder" :class="'online'">
              <div class="bg-neutral-focus text-white rounded-full w-8">
                <span class="text-xs">?</span>
              </div>
            </div>
            <div class="flex-1 min-w-0 pb-1">
              <p class="text-base font-medium"></p>
              <p class="text-xs font-light truncate text-gray-500">
                0x{{ peer.node_id }}
              </p>
            </div>
          </div>
        </li>
      </ul>
    </div>
    <div class="h-full">
      <!-- Messager -->
      <p>Messaging</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import Messager from "@/components/Messager.vue";
import SvgSmallPlus from "@/components/svgs/SvgSmallPlus.vue";
import { friends, peersInfo } from "@/store";

const firstChar = (remark: string) => (remark ? remark[0].toUpperCase() : "?");
</script>
