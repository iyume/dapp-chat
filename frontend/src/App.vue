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
        <PanelP2P />
      </template>
      <template v-if="currentPage == 'other'">
        <p>this is {{ currentPage }} page</p>
      </template>
    </template>
  </Sidebar>
</template>
