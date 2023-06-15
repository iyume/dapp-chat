<template>
  <div class="h-2"></div>
  <div class="tabs tabs-boxed">
    <a class="tab tab-active">IPFS</a>
    <a class="tab">Other</a>
  </div>
  <div class="py-2 text-red-500">
    <p>Self ID: {{ selfID }}</p>
    <p>如果需要为自身添加 Gateway，请从上面复制 Self ID</p>
  </div>
  <div class="flex flex-wrap gap-4">
    <div class="form-control w-full max-w-xs min-w-max">
      <label class="label">
        <span class="label-text">节点 ID</span>
      </label>
      <input
        v-model="nodeID"
        type="text"
        placeholder="32 bytes"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <label class="label">
        <span class="label-text">IPFS Gateway</span>
      </label>
      <input
        v-model="gateway"
        type="text"
        placeholder="127.0.0.1:port/ipns/Qmxxxx"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <div v-for="e in errors" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-error break-all">
          <ErrorIcon />
          <span>错误: {{ e }}</span>
        </div>
      </div>
      <div class="h-4"></div>
      <div class="btn btn-primary" @click="action">添加</div>
    </div>
    <div class="overflow-x-auto">
      <span>已有数据</span>
      <table class="table">
        <!-- head -->
        <thead>
          <tr>
            <th></th>
            <th>Node ID</th>
            <th>IPFS Gateway</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(providers, _, b) in storageProviders">
            <th>{{ b }}</th>
            <td>{{ providers.SelfID }}</td>
            <td>{{ providers.IPFS.Gateway }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

import ErrorIcon from "./icons/ErrorIcon.vue";

import {
  chattingNodeID,
  setIPFSGateway,
  selfID,
  storageProviders,
} from "@/store";

const props = defineProps({
  nodeID: {
    type: String,
    default: "",
  },
});

const errors = ref<string[]>([]);

const nodeID = ref<string>(props.nodeID || chattingNodeID.value);
const gateway = ref<string>("");

function action() {
  const nodeID_ = nodeID.value.replace("0x", "");
  if (nodeID_.length != 64) {
    errors.value.push("node ID must be string of length 64");
    setTimeout(() => {
      errors.value.shift();
    }, 5000);
    return;
  }
  setIPFSGateway(nodeID_, gateway.value);
}
</script>
