import axios from 'axios'
import config from '@/config'
import { type IFriend, type IMessage } from './interfaces'
import { mock_request } from './api-mock'

if (!config.p2pApiUrl || !config.p2pToken) {
  console.warn('p2p api url or token is not provided.')
}

// p2p backend should be pre-configured
const p2pApiRequest = (() => {
  // if (import.meta.env.API_MOCK) {
  //   return mock_request as any
  // } else {
  return axios.create({
    baseURL: config.p2pApiUrl,
    timeout: 2000,
    headers: {
      Authorization: `Bearer ${config.p2pToken}`
    }
  })
  // }
})()

export const api = {
  async login(username: string, password: string) {},
  async sendP2PMessage(user_id: string, message: string) {
    // TODO: message json struct; send message by user id
    p2pApiRequest.post('/send_p2p_message', { user_id, message })
  },
  async getP2PMessageList(user_id: string) {
    return p2pApiRequest.get<IMessage[]>('/get_p2p_msg_list', {
      params: { user_id }
    })
  },
  async getFriendList() {
    return p2pApiRequest.get<IFriend[]>('/get_friend_list')
  },
  async addFriend(node_id: string, remark: string) {
    p2pApiRequest.put('/add_friend', { node_id, remark })
  }
}
