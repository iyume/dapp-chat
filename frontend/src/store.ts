import { ref, type Ref, computed } from "vue";
import { useLocalStorage } from "@vueuse/core";

import type { IFriendInfo, IPeerInfo, IP2PSession } from "./interfaces";
import { api } from "./api";

export const currentPage = ref<"main" | "other">("main");
export const p2pStage = ref<"add_backend" | "add_friend" | "sync_ipfs" | null>(
  null
);
export const chattingNodeID = ref("");

export const selfID = ref("");
export const friendsInfo = ref<{ [nodeID: string]: IFriendInfo }>({});
export const peersInfo = ref<{ [nodeID: string]: IPeerInfo }>({});
/**
 * Map node ID to p2p session ref. This is shallow ref and it should be entirely updated.
 */
export var p2pSessions = ref<{ [nodeID: string]: IP2PSession }>({});

const _currentBackend = useLocalStorage("currentBackend", "");
export const currentBackend = computed(() => _currentBackend.value);
export const backends = useLocalStorage<{
  [addr: string]: {
    addr: string;
    token: string;
  };
}>("backends", {});

export async function resetBackendStores() {
  // in future, we could cache stores for each backends
  selfID.value = "";
  friendsInfo.value = {};
  peersInfo.value = {};
  p2pSessions.value = {};
  await Promise.all([
    actionGetSelfID(),
    actionGetFriends(),
    actionGetPeersInfo(),
  ]);
}

if (_currentBackend.value != "") {
  resetBackendStores();
}

export async function setBackend(addr: string) {
  if (addr == "") {
    _currentBackend.value = "";
    return;
  }
  if (!(addr in backends.value)) {
    console.error(addr);
    console.error(backends.value);
    console.error(`cannot set backend ${addr} without configuration`);
    return;
  }
  _currentBackend.value = addr;
  chattingNodeID.value = "";
  await resetBackendStores();
}

export enum FriendStatus {
  Connected,
  Disconnected,
  Notconnected,
}

/**
 * The friends list with peer info.
 */
export const friendsPeerInfo = computed(() => {
  const res: {
    [nodeID: string]: {
      status: FriendStatus;
      remote_addr: string;
    } & IFriendInfo;
  } = {};
  const peersInfo_ = peersInfo.value;
  for (let f of Object.values(friendsInfo.value)) {
    let status = FriendStatus.Notconnected;
    let remote_addr = "";
    if (f.node_id in peersInfo_) {
      let p = peersInfo_[f.node_id];
      status = p.active ? FriendStatus.Connected : FriendStatus.Disconnected;
      remote_addr = p.remote_addr;
    }
    res[f.node_id] = { ...f, status, remote_addr };
  }
  return res;
});

export async function actionGetSelfID() {
  try {
    const resp = await api.getSelfID();
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    selfID.value = resp.data.data;
  } catch (error) {
    console.error(error);
  }
}

export async function actionGetFriends() {
  try {
    const resp = await api.getFriendList();
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    const res: typeof friendsInfo.value = {};
    for (let p of resp.data.data) {
      res[p.node_id] = p;
    }
    friendsInfo.value = res;
  } catch (error) {
    console.error(error);
  }
}

export async function actionGetPeersInfo() {
  try {
    const resp = await api.getPeersInfo();
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    const res: typeof peersInfo.value = {};
    for (let p of resp.data.data) {
      res[p.node_id] = p;
    }
    peersInfo.value = res;
  } catch (error) {
    console.error(error);
  }
}

export async function actionGetP2PSession(nodeID: string) {
  try {
    const resp = await api.getP2PSession(nodeID);
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    p2pSessions.value[nodeID] = resp.data.data;
  } catch (error) {
    console.error(error);
  }
}

export async function actionSendP2PMessage(nodeID: string, message: string) {
  try {
    const resp = await api.sendP2PMessage(nodeID, message);
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    actionGetP2PSession(nodeID);
  } catch (error) {
    console.error(error);
  }
}

// TODO: scheduler
