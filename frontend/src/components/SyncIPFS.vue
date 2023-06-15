<template>
  <div class="flex flex-wrap gap-4">
    <div class="form-control w-full max-w-xs min-w-max">
      <label class="label">
        <span class="label-text">IPFS API 地址 (不需要 /api/v0)</span>
      </label>
      <input
        v-model="apiRoot"
        type="text"
        placeholder="hostname:port"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <label class="label">
        <span class="label-text">IPFS MFS 路径 (可选)</span>
      </label>
      <input
        v-model="mfsPath"
        type="text"
        placeholder="/.p2pchat"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <label class="label">
        <span class="label-text-alt text-red-500"
          >* 此路径下的文件会被全部清除</span
        >
      </label>
      <!-- FIXME: horizontal text not wrapped -->
      <div v-if="syncStatus" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-info">
          <InfoIcon />
          <span>同步状态: {{ syncStatusMessage[syncStatus] }}</span>
        </div>
      </div>
      <div v-for="e in infos" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-info">
          <InfoIcon />
          <span>注意: {{ e }}</span>
        </div>
      </div>
      <div v-for="e in errors" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-error">
          <ErrorIcon />
          <span>错误: {{ e }}</span>
        </div>
      </div>
      <div class="h-4"></div>
      <div class="btn btn-primary" @click="syncIPFS">开始同步</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

import ErrorIcon from "./icons/ErrorIcon.vue";
import InfoIcon from "./icons/InfoIcon.vue";

import { api } from "@/api";

enum SyncStatus {
  NotStarted,
  Syncing,
  Success,
  Failed,
}

const syncStatusMessage = {
  [SyncStatus.Syncing]: "同步中...",
  [SyncStatus.Success]: "同步成功",
  [SyncStatus.Failed]: "同步失败",
};

// FIXME: not context
const syncStatus = ref<SyncStatus>(SyncStatus.NotStarted);
const infos = ref<string[]>([]);
const errors = ref<string[]>([]);

const apiRoot = ref("");
const mfsPath = ref("");

function doClean() {
  syncStatus.value = SyncStatus.NotStarted;
  apiRoot.value = mfsPath.value = "";
}

async function syncIPFS() {
  syncStatus.value = SyncStatus.Syncing;
  try {
    const resp = await api.syncIPFS(apiRoot.value, mfsPath.value);
    if (resp.data.retcode != 0) {
      throw "<Backend> " + resp.data.reason;
    }
    syncStatus.value = SyncStatus.Success;
  } catch (e) {
    errors.value.push("" + e);
    setTimeout(() => {
      errors.value.shift();
    }, 5000);
    syncStatus.value = SyncStatus.Failed;
  }
}
</script>
