package api

// some json tests

// func TestJsonEncode(t *testing.T) {
// 	assert := assert.New(t)

// 	now := time.Now()
// 	event := P2PMessageEvent{
// 		MessageEvent: MessageEvent{
// 			Event:   Event{Time: now, Type: "message", DetailType: "private"},
// 			Message: Message{},
// 		},
// 	}
// 	bytes, err := json.Marshal(event)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	assert.Equal(fmt.Sprintf(`{"time":"%v","type":"message","detail_type":"private","message":[]}`, now.Format(time.RFC3339Nano)), string(bytes))
// }

// func TestJsonDecode(t *testing.T) {
// 	assert := assert.New(t)

// 	now := time.Now()
// 	// extra field will be ignored
// 	bytes := []byte(fmt.Sprintf(`{"id":"1","time":"%v","type":"message","detail_type":"private","message":[],"extra":""}`, now.Format(time.RFC3339Nano)))
// 	event := P2PMessageEvent{}
// 	if err := json.Unmarshal(bytes, &event); err != nil {
// 		log.Fatalln(err)
// 	}
// 	// remove time field and use Equal() to compare
// 	eventtime := event.Time
// 	event.Time = time.Time{}
// 	assert.Equal(event, P2PMessageEvent{
// 		MessageEvent: MessageEvent{
// 			Event:   Event{Type: "message", DetailType: "private"},
// 			Message: Message{},
// 		},
// 	})
// 	if !eventtime.Equal(now) {
// 		log.Fatalln("Time not equal")
// 	}
// }
