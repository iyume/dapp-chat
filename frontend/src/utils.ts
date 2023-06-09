import type { IMessage, ITextSegment } from "./interfaces";

export default {
  extractPlainText: (msg: IMessage) => {
    let text = "";
    for (let seg of msg) {
      if (seg.type == "text") {
        text += (seg.data as ITextSegment).text;
      }
    }
    return text;
  },
  /**
   * Convert time to readable time in message sent.
   * @param time RFC-3339 (ISO-8601) time
   */
  sentTimeChat: (time: string) => {
    const date = new Date(time);
    const now = new Date();
    var resDate = date.toLocaleDateString("en-US");
    const resTime = date.toLocaleTimeString("en-US");
    if (resDate == now.toLocaleDateString("en-US")) {
      resDate = "今天";
    }
    return `${resDate} ${resTime}`;
  },
};
