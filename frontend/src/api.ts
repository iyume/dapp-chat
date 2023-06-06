import axios from "axios";
import { computed } from "vue";

import type { IResp, IFriendInfo, IPeerInfo, IP2PSession } from "./interfaces";
import { currentBackend, backends } from "./store";

const p2pApiRequest = computed(() =>
  axios.create({
    baseURL: currentBackend.value,
    timeout: 10000,
    headers: {
      Authorization: `Bearer ${backends.value[currentBackend.value].token}`,
    },
  })
);

export const api = {
  async getSelfID() {
    return p2pApiRequest.value.get<IResp<string>>("/get_self_id");
  },
  async getPeersInfo() {
    return p2pApiRequest.value.get<IResp<IPeerInfo[]>>("/get_peers_info");
  },
  async getFriendList() {
    return p2pApiRequest.value.get<IResp<IFriendInfo[]>>("/get_friend_list");
  },
  async addFriend(node_id: string, remark: string) {
    p2pApiRequest.value.get("/add_friend", { params: { node_id, remark } });
  },
  async getP2PSession(node_id: string) {
    return p2pApiRequest.value.get<IResp<IP2PSession>>("/get_p2p_session", {
      params: { node_id },
    });
  },
  async sendP2PMessage(node_id: string, message: string) {
    // TODO: message json struct; send message by user id
    p2pApiRequest.value.get("/send_p2p_message", {
      params: { node_id, message },
    });
  },
};
