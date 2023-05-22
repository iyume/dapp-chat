<script setup lang="ts">
import { currentPage } from "./store";
import P2PView from "./views/P2PView.vue";
import Navbar from "./components/Navbar.vue";
import Sidebar from "./components/Sidebar.vue";
</script>

<template>
  <!-- <div class="bg-red-500">
    <h1 class="text-3xl font-bold underline">Hello world!</h1>
    <button class="btn">Hello daisyUI</button>
  </div> -->
  <Navbar drawer-id="my-drawer-main" />
  <Sidebar drawer-id="my-drawer-main" class="h-[calc(100vh-4rem)]">
    <template #content>
      <template v-if="currentPage == 'main'">
        <p>this is main page</p>
      </template>
      <template v-if="currentPage == 'p2p'">
        <P2PView />
      </template>
    </template>
  </Sidebar>
</template>
