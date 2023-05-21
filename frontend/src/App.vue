<template>
  <div class="container mx-auto h-screen">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { useDataStore } from './stores/data'
import { faker } from '@faker-js/faker'

const datastore = useDataStore()

for (let i = 1; i <= Math.round(Math.random() * 30); ++i) {
  datastore.friends.push({
    node_id: faker.string.uuid(),
    id: i,
    remark: faker.person.fullName(),
    avatar: faker.image.urlLoremFlickr({ width: 64, height: 64 })
  })
}

for (let i = 1; i <= Math.round(Math.random() * 100); ++i) {
  datastore.messages.push({
    friend_id: Math.round(Math.random() * datastore.friends.length),
    direction: Math.random() < 0.5 ? 1 : 2,
    message_id: i,
    message: faker.lorem.words({ min: 1, max: 20 }),
    time: faker.date.recent().toString()
  })
}

datastore.my_info = {
  name: faker.person.fullName(),
  avatar: faker.image.urlLoremFlickr({ width: 64, height: 64 })
}
</script>
