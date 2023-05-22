export interface IFriend {
  remark: string;
  node_id: string;
}

export interface IP2pMessage {
  time: string;
  node_id: string;
  message_id: number;
  message: string;
}
