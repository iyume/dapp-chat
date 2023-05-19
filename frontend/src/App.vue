<script lang="ts" setup>
import Navbar from "./components/Navbar.vue";
import Sidebar from "./components/Sidebar.vue";
import { currentPage } from "./store";
</script>

<template>
  <!-- <div class="bg-red-500">
    <h1 class="text-3xl font-bold underline">Hello world!</h1>
    <button class="btn">Hello daisyUI</button>
  </div> -->
  <Navbar drawer-id="my-drawer-main" />
  <Sidebar drawer-id="my-drawer-main">
    <template #content>
      <template v-if="currentPage == 'main'">
        <p>this is main page</p>
      </template>
      <template v-if="currentPage == 'p2p'">
        <p>this is p2p page</p>
      </template>
    </template>
  </Sidebar>
</template>
