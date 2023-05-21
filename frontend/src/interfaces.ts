export interface IFriend {
  id: number
  remark: string
  node_id: string
  avatar?: string
}

export interface IMessage {
  friend_id: number
  direction: 1 | 2 // 1: receive, 2: send
  time: string
  message_id: number
  message: string
}
