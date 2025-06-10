<template>
	<div class="home-container">
	  <div class="home-header">
		<h2>Welcome, {{ username }} 👋</h2>
		<button @click="logout" class="logout-button">Logout</button>
	  </div>
  
	  <div v-if="error" class="error">{{ error }}</div>
	  <div v-if="loading">Loading your stream...</div>
  
	  <div v-if="images.length === 0 && !loading" class="empty-message">
		No posts to show yet. Start following users!
	  </div>
  
	  <div v-for="image in images" :key="image.imageurl" class="image-card">
		<img :src="image.imageurl" alt="Posted image" class="image-preview" />
		<div class="image-meta">
			<div class="image-footer">
				<p><strong>Posted by:</strong> {{ image.username }}</p>

				<div class="like-section">
					<span>{{ image.likes }}</span>
					<button class="like-button" @click="likeImage(image.imageurl)">
					❤️
					</button>
				</div>
			</div>

			<div>
				<strong>Comments:</strong>
				<ul class="comment-list">
					<li v-for="(comment, idx) in image.comments.split('~')" :key="idx">{{ comment }}</li>
				</ul>
			</div>

		</div>
	  </div>
	</div>
  </template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const images = ref([])
const error = ref('')
const loading = ref(true)

const logout = () => {
localStorage.removeItem('username')
router.push('/')
}

const fetchStream = async () => {
	try {
		const res = await axios.get(`/users/${username.value}/stream`)
		images.value = res.data || []
	} catch (err) {
		error.value = 'Failed to load stream'
		console.error(err)
	} finally {
		loading.value = false
	}
}

const likeImage = async (imageurl) => {
  try {
    const encodedURL = encodeURIComponent(imageurl)
	console.log('Encoded URL:', encodedURL)
    await axios.put(`/images/${encodedURL}/like`)
    const img = images.value.find(img => img.imageurl === imageurl)
    if (img) img.likes += 1
  } catch (err) {
    console.error('Failed to like image:', err)
  }
}



onMounted(() => {
if (!username.value) {
	router.push('/')
} else {
	fetchStream()
}
})
</script>

<style scoped>
.home-container {
max-width: 800px;
margin: 2rem auto;
padding: 1rem;
}

.home-header {
display: flex;
justify-content: space-between;
align-items: center;
}

.logout-button {
background: #ff4d4d;
border: none;
padding: 0.5rem 1rem;
color: white;
border-radius: 6px;
cursor: pointer;
}

.image-card {
border: 1px solid #ddd;
padding: 1rem;
border-radius: 8px;
margin-top: 1.5rem;
background: white;
}

.image-preview {
width: 100%;
border-radius: 6px;
margin-bottom: 0.5rem;
}

.image-meta p {
margin: 0.2rem 0;
}

.image-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.5rem;
}

.like-section {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.like-button {
  background: transparent;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  transition: transform 0.2s;
}

.like-button:hover {
  transform: scale(1.2);
}


.error {
color: red;
margin-top: 1rem;
}

.empty-message {
margin-top: 2rem;
font-style: italic;
}
.comment-list {
  margin: 0.3rem 0 0 1rem;
  padding: 0;
  list-style-type: disc;
}

</style>
