export const mock_request = {
  async post(endpoint: string, params?: any): Promise<any> {
    switch (endpoint) {
      case '/send_p2p_message':
        return {
          message_id: 1
        }
      default:
        return {}
    }
  },
  async get<T = any>(endpoint: string, config?: { params: any }): Promise<T> {
    switch (endpoint) {
      case '/get_p2p_msg_list':
        return [
          {
            sender_id: config?.params.user_id,
            time: '2023-05-16T19:30:33Z',
            message_id: 1,
            message: 'test message 1'
          },
          {
            sender_id: config?.params.user_id,
            time: '2023-05-16T19:31:33Z',
            message_id: 2,
            message: 'test message 2'
          },
          {
            sender_id: config?.params.user_id,
            time: '2023-05-16T19:38:33Z',
            message_id: 3,
            message: 'test message 3'
          }
        ] as T
      case '/get_friend_list':
        return [
          {
            id: 12,
            remark: 'friend 1',
            node_id: '3e45'
          },
          {
            id: 13,
            remark: 'friend 1',
            node_id: '3e45'
          },
          {
            id: 14,
            remark: 'friend 1',
            node_id: '3e45'
          },
          {
            id: 15,
            remark: 'friend 1',
            node_id: '3e45'
          },
          {
            id: 16,
            remark: 'friend 1',
            node_id: '3e45'
          }
        ] as T
      default:
        return {} as T
    }
  },
  async put(endpoint: string, params?: any): Promise<any> {
    switch (endpoint) {
      case '/add_friend':
        return {}
      default:
        return {}
    }
  }
}
