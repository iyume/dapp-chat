<template>
  <h1 class="text-red-500 text-lg font-bold">
    注意，您正在将当前后端 (0x{{ selfID.slice(0, 7) }}...) 上传至 IPFS，默认会将
    /.p2pchat 作为数据目录，并且使用 self 密钥发布
  </h1>
  <div class="flex flex-wrap gap-4">
    <div class="form-control w-full max-w-xs min-w-max">
      <label class="label">
        <span class="label-text">IPFS API 地址 (不需要加 /api/v0)</span>
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
      <label class="label">
        <span class="label-text">Publish Key (可选)</span>
      </label>
      <input
        v-model="key"
        type="text"
        placeholder="self"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <div v-if="syncStatus" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-info break-all">
          <InfoIcon />
          <span>同步状态: {{ syncStatusMessage[syncStatus] }}</span>
        </div>
      </div>
      <div v-if="successResp" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-success break-all">
          <SuccessIcon />
          <span
            >成功将 {{ successResp.value }} <br />发布至 /ipns/{{
              successResp.name
            }}<br />现在您可以从 IPFS 配置的 Gateway 进行访问</span
          >
        </div>
      </div>
      <div v-for="e in infos" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-info break-all">
          <InfoIcon />
          <span>注意: {{ e }}</span>
        </div>
      </div>
      <div v-for="e in errors" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-error break-all">
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
import SuccessIcon from "./icons/SuccessIcon.vue";

import { api } from "@/api";
import type { PublishedIPFS } from "@/interfaces";
import { selfID } from "@/store";

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
const successResp = ref<PublishedIPFS>();

const apiRoot = ref("");
const mfsPath = ref("");
const key = ref("");

async function syncIPFS() {
  successResp.value = undefined;
  syncStatus.value = SyncStatus.Syncing;
  try {
    const resp = await api.syncIPFS(
      apiRoot.value,
      mfsPath.value,
      key.value,
      60000
    );
    if (resp.data.retcode != 0) {
      throw "<Backend> " + resp.data.reason;
    }
    syncStatus.value = SyncStatus.Success;
    successResp.value = resp.data.data;
  } catch (e) {
    errors.value.push("" + e);
    setTimeout(() => {
      errors.value.shift();
    }, 5000);
    syncStatus.value = SyncStatus.Failed;
  }
}
</script>
