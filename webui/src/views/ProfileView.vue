<template>
  <div class="profile-container">
    <div class="user-info">
      <h2>{{ username }}</h2>
      <button class="change-btn" @click="changeUsername">Change Username</button>
      <div class="stats">
        <span @click="showFollowingModal = true">Following ({{ followingList.length }})</span>
        <span @click="showBannedModal = true">Banned ({{ bannedList.length }})</span>
      </div>
    </div>

    <div class="image-grid">
      <div class="upload-tile" @click="uploadImage">⬆️</div>
      <div v-for="img in images" :key="img.id" class="image-item">
        <img :src="img.imageurl" alt="photo" />
      </div>
    </div>

    <div v-if="showFollowingModal" class="modal-overlay" @click.self="showFollowingModal = false">
      <div class="modal">
        <h3>Following</h3>
        <ul>
          <li v-for="user in followingList" :key="user">
            <span>{{ user }}</span>
            <button @click="toggleFollow(user)" :class="{ active: !isFollowing(user) }">
              {{ isFollowing(user) ? 'Unfollow' : 'Follow' }}
            </button>
          </li>
        </ul>
        <button @click="showFollowingModal = false">Close</button>
      </div>
    </div>

    <div v-if="showBannedModal" class="modal-overlay" @click.self="showBannedModal = false">
      <div class="modal">
        <h3>Banned</h3>
        <ul>
          <li v-for="user in bannedList" :key="user">
            <span>{{ user }}</span>
            <button @click="toggleBan(user)" :class="{ active: !isBanned(user) }">
              {{ isBanned(user) ? 'Unban' : 'Ban' }}
            </button>
          </li>
        </ul>
        <button @click="showBannedModal = false">Close</button>
      </div>
    </div>
  </div>
</template>

<script setup>
const showFollowingModal = ref(false);
const showBannedModal = ref(false);

import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const followingList = ref([])
const bannedList = ref([])
const images = ref([])


const fetchProfile = async () => {
  try {
    console.log("fetching profile")
    const res = await axios.get(`/users/${username.value}`)
    //console.log(res.data.Following)
    followingList.value = res.data.Following ? res.data.Following.split(',').filter(v => v) : []
    //console.log(followingList)
    bannedList.value = res.data.Banned ? res.data.Banned.split(',').filter(v => v) : []
  } catch (err) {
    console.error('Failed to load profile', err)
  }
}

const fetchImages = async () => {
  try {
    console.log("fetching images")
    const res = await axios.get(`/users/${username.value}/photos`)
    images.value = res.data || []
  } catch (err) {
    console.error('Failed to load images', err)
  }
}

const changeUsername = async () => {
  const newName = prompt('New username:')
  if (!newName) return
  try {
    await axios.put(`/users/${username.value}`, { username: newName })
    localStorage.setItem('username', newName)
    username.value = newName
    await fetchProfile()
    await fetchImages()
  } catch (err) {
    console.error('Failed to change username', err)
  }
}

const uploadImage = async () => {
  const newUrl = prompt('Image url:')
  if (!newUrl) return
  try {
    await axios.post('/images', { username: username.value, imageurl: newUrl })
    await fetchImages()
  } catch (err) {
    console.error('Failed to upload image', err)
  }
}

const isFollowing = (user) => followingList.value.includes(user)
const isBanned = (user) => bannedList.value.includes(user)

const toggleFollow = async (user) => {
  try {
    if (isFollowing(user)) {
      await axios.delete(`/users/${username.value}/follow`, { data: { username: user } })
      followingList.value = followingList.value.filter(u => u !== user)
    } else {
      await axios.put(`/users/${username.value}/follow`, { username: user })
      followingList.value.push(user)
    }
  } catch (err) {
    console.error('Failed to toggle follow', err)
  }
}

const toggleBan = async (user) => {
  try {
    if (isBanned(user)) {
      await axios.delete(`/users/${username.value}/ban`, { data: { username: user } })
      bannedList.value = bannedList.value.filter(u => u !== user)
    } else {
      await axios.put(`/users/${username.value}/ban`, { username: user })
      bannedList.value.push(user)
    }
  } catch (err) {
    console.error('Failed to toggle ban', err)
  }
}

onMounted(() => {
  if (!username.value) {
    router.push('/login')
    return
  }
  fetchProfile()
  fetchImages()
})
</script>

<style scoped>
.profile-container {
  padding: 1rem;
}

.user-info {
  border: 1px solid #ddd;
  padding: 1rem;
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.change-btn {
  width: fit-content;
}

.stats span {
  margin-right: 1rem;
  cursor: pointer;
  color: white;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 10px;
}

.upload-tile, .image-item img {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.upload-tile {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  cursor: pointer;
  font-size: 2rem;
}


.modal {
  background: white;
  padding: 1rem;
  border-radius: 4px;
  max-width: 400px;
  width: 90%;
}

.modal ul {
  list-style: none;
  padding: 0;
}

.modal li {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.modal button.active {
  background: #4caf50;
  color: white;
}

.modal button {
  padding: 0.3rem 0.6rem;
}
</style>


<style>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 20px;
  border-radius: 10px;
  min-width: 300px;
  max-width: 90vw;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.change-btn,
.stats span {
  cursor: pointer;
  margin: 5px;
  padding: 8px 12px;
  background-color: #000000;
  color: white;
  border: none;
  border-radius: 5px;
  transition: background-color 0.3s ease;
}

.change-btn:hover,
.stats span:hover {
  background-color: #707070;
}
</style>
