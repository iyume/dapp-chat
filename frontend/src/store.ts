import { ref } from "vue";
import type { IFriend } from "./interfaces";
import { api } from "./api";

export const currentPage = ref<"main" | "p2p">("main");
export const friendList = ref<IFriend[]>([]);

async function commitFriendList() {
  try {
    const resp = await api.getFriendList();
    if (resp.status != 200 || resp.data.retcode != 0) {
      throw "request failed";
    }
    friendList.value = resp.data.data;
  } catch (error) {
    console.error(error);
  }
}

commitFriendList();
