<template>
  <div class="flex flex-wrap gap-4">
    <div class="form-control w-full max-w-xs min-w-max">
      <label class="label">
        <span class="label-text">后端地址</span>
      </label>
      <input
        v-model="addr"
        type="text"
        placeholder="hostname:port"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <label class="label">
        <span class="label-text">Token</span>
      </label>
      <input
        v-model="token"
        type="text"
        placeholder="server token"
        class="input input-bordered input-primary w-full max-w-xs"
      />
      <div v-if="addr in backends">
        <div class="h-4"></div>
        <div class="alert alert-warning text-sm">
          <WarningIcon />
          <span>警告: 将覆盖已有的 {{ backends[addr] }}</span>
        </div>
      </div>
      <div v-for="e in errors">
        <div class="h-4"></div>
        <div class="alert alert-error">
          <ErrorIcon />
          <span>错误: {{ e }}</span>
        </div>
      </div>
      <div class="h-4"></div>
      <div class="btn btn-primary" @click="addBackend">添加后端</div>
    </div>
    <div class="form-control w-full max-w-xs">
      <label class="label">
        <span class="label-text">选择后端进行删除</span>
      </label>
      <select v-model="delAddr" class="select select-primary w-full max-w-xs">
        <option disabled value="" hidden></option>
        <option v-for="b in backends">
          {{ b.addr }}
        </option>
      </select>
      <div class="h-4"></div>
      <div class="btn btn-primary" @click="deleteBackend">删除后端</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, type PropType } from "vue";

import WarningIcon from "./icons/WarningIcon.vue";
import ErrorIcon from "./icons/ErrorIcon.vue";
import { backends, setBackend, currentBackend } from "@/store";

const props = defineProps({
  // as unknown as () => void
  exit: Function as PropType<() => void>,
});

// FIXME: not context
const errors = ref<string[]>([]);

const addr = ref("");
const token = ref("");

function addBackend() {
  if (addr.value == "" || token.value == "") {
    errors.value.push("表单不能为空");
    setTimeout(() => errors.value.shift(), 2000);
    return;
  }
  backends.value[addr.value] = { addr: addr.value, token: token.value };
  addr.value = token.value = "";
}

const delAddr = ref(Object.keys(backends.value)[0]);

function deleteBackend() {
  delete backends.value[delAddr.value];
  if (delAddr.value == currentBackend.value) {
    setBackend("");
  }
}
</script>
