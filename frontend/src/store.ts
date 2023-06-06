import { ref, type Ref, computed } from "vue";
import { useLocalStorage } from "@vueuse/core";

import type { IFriendInfo, IPeerInfo, IP2PSession } from "./interfaces";
import { api } from "./api";

export const currentPage = ref<"main" | "other">("main");
export const p2pStage = ref<"add_backend" | "add_friend" | null>(null);
export const friendsInfo = ref<{ [nodeID: string]: IFriendInfo }>({});
export const peersInfo = ref<{ [nodeID: string]: IPeerInfo }>({});
export const selfID = ref("");
/**
 * Map node ID to p2p session ref. This is shallow ref and it should be entirely updated.
 */
export const p2pSessions: { [key: string]: Ref<IP2PSession> } = {};

export const currentBackend = useLocalStorage("currentBackend", "");
export const backends = useLocalStorage<{
  [addr: string]: {
    addr: string;
    token: string;
  };
}>("backends", {});

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
    if (nodeID in p2pSessions) {
      p2pSessions[nodeID].value = resp.data.data;
    }
    p2pSessions[nodeID] = ref(resp.data.data);
  } catch (error) {
    console.error(error);
  }
}

// TODO: scheduler
actionGetSelfID();
actionGetFriends();
actionGetPeersInfo();
