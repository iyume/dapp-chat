<template>
  <div class="tabs tabs-boxed">
    <a class="tab tab-active">IPFS</a>
    <a class="tab">Other</a>
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
        placeholder="hostname:port"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <div v-for="e in errors" class="max-w-xs">
        <div class="h-4"></div>
        <div class="alert alert-error">
          <ErrorIcon />
          <span>错误: {{ e }}</span>
        </div>
      </div>
      <div class="h-4"></div>
      <div class="btn btn-primary" @click="action">添加</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

import { chattingNodeID, setIPFSGateway } from "@/store";

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
