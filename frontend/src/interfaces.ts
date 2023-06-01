export interface IResp<T> {
  retcode: number;
  data: T;
}

export interface IPeerInfo {
  node_id: string;
  active: boolean;
  version: number;
}

// TODO: add remote addr

export interface IFriend {
  node_id: string;
  remark: string;
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
