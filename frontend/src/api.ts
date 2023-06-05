import axios from "axios";

import type { IResp, IFriendInfo, IPeerInfo, IP2PSession } from "./interfaces";
import config from "@/config";

if (!config.p2pApiUrl || !config.p2pToken) {
  console.warn("p2p api url or token is not provided.");
}

// TODO: use store.ts/currentBackend to refactor it
const p2pApiRequest = axios.create({
  baseURL: config.p2pApiUrl,
  timeout: 10000,
  headers: {
    Authorization: `Bearer ${config.p2pToken}`,
  },
});

export const api = {
  async login(username: string, password: string) {},
  async getPeersInfo() {
    return p2pApiRequest.get<IResp<IPeerInfo[]>>("/get_peers_info");
  },
  async getFriendList() {
    return p2pApiRequest.get<IResp<IFriendInfo[]>>("/get_friend_list");
  },
  async addFriend(node_id: string, remark: string) {
    p2pApiRequest.get("/add_friend", { params: { node_id, remark } });
  },
  async getP2PSession(node_id: string) {
    return p2pApiRequest.get<IResp<IP2PSession>>("/get_p2p_session", {
      params: { node_id },
    });
  },
  async sendP2PMessage(user_id: string, message: string) {
    // TODO: message json struct; send message by user id
    p2pApiRequest.get("/send_p2p_message", { params: { user_id, message } });
  },
};
