import { faker } from "@faker-js/faker";
import MockAdaptor from "axios-mock-adapter";
import { type IFriend, type IP2pMessage } from "./interfaces";
import type { AxiosStatic } from "axios";

export function use_mock(axios: AxiosStatic) {
  console.warn("api mock installed");
  const mock = new MockAdaptor(axios);

  const friends: IFriend[] = [];
  const messages: IP2pMessage[] = [];
  for (let i = 1; i <= Math.round(Math.random() * 30); ++i) {
    friends.push({
      node_id: faker.string.uuid(),
      remark: faker.person.fullName(),
    });
  }

  let msgid = 1;

  friends.forEach((f) => {
    for (let i = 1; i <= Math.round(Math.random() * 20); ++i) {
      messages.push({
        message_id: ++msgid,
        node_id: f.node_id,
        message: faker.lorem.words({ min: 1, max: 30 }),
        time: faker.date.recent().toISOString(),
      });
    }
  });

  mock.onGet("/get_friend_list").reply<IFriend[]>(200, friends);
  mock.onGet("/get_p2p_msg_list").reply<IP2pMessage[]>((req) => {
    let ret = messages.filter((m) => m.node_id == req.params.node_id);
    return [200, ret];
  });
}
