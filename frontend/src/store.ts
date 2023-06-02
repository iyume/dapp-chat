import { ref } from "vue";
import type { IFriend, IPeerInfo } from "./interfaces";
import { api } from "./api";

export const currentPage = ref<"main" | "other">("main");
export const friends = ref<IFriend[]>([]);
export const peersInfo = ref<IPeerInfo[]>([]);

async function commitFriends() {
  try {
    const resp = await api.getFriendList();
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    friends.value = resp.data.data;
  } catch (error) {
    console.error(error);
  }
}

async function commitPeersInfo() {
  try {
    const resp = await api.getPeersInfo();
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    peersInfo.value = resp.data.data;
  } catch (error) {
    console.error(error);
  }
}

// TODO: scheduler
commitFriends();
commitPeersInfo();
