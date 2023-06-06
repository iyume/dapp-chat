<template>
  <div class="dropdown">
    <label tabindex="0" class="btn btn-sm btn-primary whitespace-nowrap">
      <!-- TODO: sche check backend status and show -->
      后端: {{ currentBackend || "未选择" }}
    </label>
    <ul
      tabindex="0"
      class="dropdown-content menu p-2 shadow bg-base-100 rounded-box backdrop-blur bg-opacity-50 z-20 whitespace-nowrap"
    >
      <li v-for="b in backends" @click="setBackend(b.addr)">
        <button class="flex items-center gap-x-1 text-sm px-2">
          <MiniCheckIcon />{{ b.addr }}
        </button>
      </li>
      <li @click="p2pStage = 'add_backend'">
        <button class="flex items-center gap-x-1 text-sm px-2">
          <MiniPlusIcon />添加后端
        </button>
      </li>
    </ul>
  </div>
</template>
<script setup lang="ts">
import MiniCheckIcon from "./icons/MiniCheckIcon.vue";
import MiniPlusIcon from "./icons/MiniPlusIcon.vue";
import { currentBackend, setBackend, backends, p2pStage } from "@/store";
</script>
