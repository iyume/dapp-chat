import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { IFriend, IMessage } from '@/interfaces'

export const useDataStore = defineStore('datastore', () => {
  const messages = ref<IMessage[]>([])
  const friends = ref<IFriend[]>([])
  const my_info = ref<{ name: string; avatar: string | undefined } | undefined>()
  return { messages, friends, my_info }
})
