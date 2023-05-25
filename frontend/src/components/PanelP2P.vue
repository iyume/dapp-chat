<template>
  <div class="h-full flex bg-gray-100">
    <!-- Friend list -->
    <div class="p-4">
      <div class="max-w-xs h-full p-3 bg-white rounded">
        <div class="p-3 flex items-center justify-between">
          <p class="font-bold inline-block align-middle">好友节点</p>
          <button class="btn btn-ghost btn-sm"><SvgSmallPlus /></button>
        </div>
        <ul class="divide-y divide-gray-100">
          <li
            v-for="friend in friendList"
            class="gap-x-6 rounded-md cursor-pointer hover:bg-gray-100"
          >
            <div class="flex gap-x-4 min-w-0 p-2">
              <div class="avatar online placeholder">
                <div
                  class="bg-neutral-focus text-neutral-content rounded-full w-12"
                >
                  <span>M</span>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import Messager from "@/components/Messager.vue";
import SvgSmallPlus from "@/components/SvgSmallPlus.vue";
import { friendList } from "@/store";
</script>
