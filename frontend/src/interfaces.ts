export interface IResp<T> {
  retcode: number;
  data: T;
}

export interface IPeerInfo {
  node_id: string;
  active: boolean;
  version: number;
  remote_addr: string;
}

export interface IFriendInfo {
  node_id: string;
  remark: string;
}

export type IMessage = ISegment[];

interface ISegment {
  type: string;
  data: any;
}

export interface ITextSegment {
  text: string;
}

interface IEvent {
  time: string; // RFC3339Nano
  type: string;
  detail_type: string;
}

interface IMessageEvent extends IEvent {
  message: IMessage;
}

export interface IP2PMessageEvent extends IMessageEvent {
  user_id: string;
}

export interface IP2PSession {
  events: IP2PMessageEvent[];
}
