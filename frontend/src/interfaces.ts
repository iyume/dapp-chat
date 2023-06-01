export interface IResp<T> {
  retcode: number;
  data: T;
}

export interface IPeerInfo {
  node_id: string;
  active: boolean;
  version: number;
}

export interface IFriend {
  remark: string;
  node_id: string;
}

export interface IP2PSession {
  events: IMessage[];
}

export interface IMessage {
  time: string;
  node_id: string;
  message_id: number;
  message: string;
}
