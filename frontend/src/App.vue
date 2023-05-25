<script setup lang="ts">
import { currentPage } from "./store";
import PanelP2P from "./components/PanelP2P.vue";
import Navbar from "./components/Navbar.vue";
import Sidebar from "./components/Sidebar.vue";
</script>

<template>
  <Navbar drawer-id="my-drawer-main" />
  <Sidebar drawer-id="my-drawer-main" class="h-[calc(100vh-4rem)]">
    <template #content>
      <template v-if="currentPage == 'main'">
        <p>this is main page</p>
      </template>
      <template v-if="currentPage == 'p2p'">
        <PanelP2P />
      </template>
    </template>
  </Sidebar>
</template>
