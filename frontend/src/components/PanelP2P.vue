<template>
  <div class="h-full flex flex-row bg-gray-100">
    <!-- Friend list -->
    <div class="max-w-xs w-full p-4">
      <div class="h-full p-3 bg-white rounded flex flex-col">
        <div class="p-3 flex items-center justify-between">
          <p class="font-bold inline-block">好友节点</p>
          <button class="btn btn-ghost btn-sm"><SvgSmallPlus /></button>
        </div>
        <ul class="divide-y divide-gray-100 overflow-y-auto">
          <li
            v-for="friend in friendList"
            class="rounded-md cursor-pointer hover:bg-gray-100"
          >
            <div class="flex gap-x-4 min-w-0 p-2">
              <div class="avatar online placeholder">
                <div
                  class="bg-neutral-focus text-neutral-content rounded-full w-12"
                >
                  <span>{{ firstChar(friend.remark) }}</span>
                </div>
              </div>
              <div class="min-w-0 flex flex-col justify-around">
                <p class="text-base font-semibold">{{ friend.remark }}</p>
                <p class="text-xs truncate text-gray-500">
                  0x{{ friend.node_id }}
                </p>
              </div>
            </div>
          </li>
        </ul>
      </div>
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
import { friendList } from "@/store";

const firstChar = (remark: string) => (remark ? remark[0].toUpperCase() : "?");
</script>
