<template>
  <!-- Name entry modal (shown to guests who land directly on the room URL) -->
  <div v-if="showNameModal" class="modal-backdrop">
    <div class="modal">
      <div class="modal__logo">🃏</div>
      <h2 class="modal__title">Entrar na sala</h2>
      <label class="modal__label">Seu nome</label>
      <input
        v-model="pendingName"
        placeholder="Ex: Bob"
        maxlength="30"
        autofocus
        @keyup.enter="joinRoom"
      />
      <p v-if="notFound" class="modal__error">Sala não encontrada.</p>
      <button class="modal__btn" :disabled="!pendingName.trim()" @click="joinRoom">
        Entrar
      </button>
    </div>
  </div>

  <!-- Main room UI -->
  <div v-else class="room">
    <!-- Header -->
    <header class="room__header">
      <div class="room__brand">🃏 Planning Poker</div>
      <div class="room__share">
        <span class="room__share-label">Compartilhe:</span>
        <code class="room__url">{{ roomUrl }}</code>
        <button class="room__copy-btn" @click="copyUrl" :class="{ 'room__copy-btn--ok': copied }">
          {{ copied ? '✓ Copiado' : 'Copiar' }}
        </button>
      </div>
      <div class="room__status" :class="store.connected ? 'room__status--on' : 'room__status--off'">
        {{ store.connected ? 'Conectado' : 'Desconectado' }}
      </div>
    </header>

    <!-- Error toast -->
    <Transition name="toast">
      <div v-if="store.error" class="toast">{{ store.error }}</div>
    </Transition>

    <!-- Table -->
    <main class="room__main">
      <RoundTable
        v-if="store.room"
        :players="store.room.players"
        :masterid="store.room.master_id"
        :my-id="store.myId"
        :revealed="store.room.revealed"
        :round="store.room.round"
        :winners="store.winners"
        :average="store.average"
        :all-voted="store.allVoted"
      />
      <div v-else class="room__loading">Conectando…</div>
    </main>

    <!-- Controls -->
    <footer class="room__footer">
      <!-- Card deck (hidden after reveal) -->
      <CardDeck
        v-if="store.room && !store.room.revealed"
        :cards="cards"
        :selected="store.myPlayer?.vote ?? ''"
        :disabled="false"
        @vote="store.vote"
      />

      <!-- Set master dropdown (creator only) -->
      <div v-if="store.isCreator && store.room" class="room__master-panel">
        <label class="room__master-label">Scrum Master:</label>
        <select class="room__master-select" @change="changeMaster($event.target.value)">
          <option
            v-for="p in store.room.players"
            :key="p.id"
            :value="p.id"
            :selected="p.id === store.room.master_id"
          >{{ p.name }}</option>
        </select>
      </div>

      <!-- Master actions -->
      <div v-if="store.isMaster && store.room" class="room__actions">
        <button
          v-if="!store.room.revealed"
          class="room__btn room__btn--reveal"
          :disabled="!store.allVoted"
          @click="store.reveal()"
        >
          👁 Revelar cartas
        </button>
        <button
          v-if="store.room.revealed"
          class="room__btn room__btn--reset"
          @click="store.reset()"
        >
          🔄 Nova rodada
        </button>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRoomStore } from '../stores/room.js'
import RoundTable from '../components/RoundTable.vue'
import CardDeck   from '../components/CardDeck.vue'

const route  = useRoute()
const router = useRouter()
const store  = useRoomStore()

const roomId      = route.params.id
const showNameModal = ref(false)
const pendingName   = ref('')
const notFound      = ref(false)
const copied        = ref(false)
const cards         = ref([])

const roomUrl = computed(() => `${location.origin}/room/${roomId}`)

async function fetchCards() {
  const res = await fetch('/api/cards')
  if (res.ok) {
    const data = await res.json()
    cards.value = data.cards
  }
}

async function checkRoom() {
  const res = await fetch(`/api/rooms/${roomId}`)
  if (!res.ok) { notFound.value = true; return false }
  return true
}

onMounted(async () => {
  await fetchCards()

  const storedName      = sessionStorage.getItem('poker_name') || ''
  const storedCreatorId = sessionStorage.getItem(`poker_creator_${roomId}`) || ''

  if (storedName) {
    const exists = await checkRoom()
    if (!exists) { showNameModal.value = true; return }
    store.connect(roomId, storedName, storedCreatorId)
  } else {
    showNameModal.value = true
  }
})

onUnmounted(() => store.disconnect())

async function joinRoom() {
  if (!pendingName.value.trim()) return
  const exists = await checkRoom()
  if (!exists) return
  sessionStorage.setItem('poker_name', pendingName.value.trim())
  showNameModal.value = false
  store.connect(roomId, pendingName.value.trim())
}

function changeMaster(playerId) {
  store.setMaster(playerId)
}

async function copyUrl() {
  await navigator.clipboard.writeText(roomUrl.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}
</script>

<style scoped>
/* Modal */
.modal-backdrop {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: var(--color-bg);
}

.modal {
  background: var(--color-surface);
  border-radius: 16px;
  padding: 32px;
  width: 100%;
  max-width: 360px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 8px 40px rgba(0,0,0,0.5);
}

.modal__logo  { font-size: 2.5rem; text-align: center; }
.modal__title { font-size: 1.4rem; font-weight: 700; text-align: center; }
.modal__label {
  font-size: 0.82rem; font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase; letter-spacing: 0.06em;
}
.modal__btn {
  background: var(--color-gold);
  color: #1a1a1a;
  font-weight: 700;
  padding: 12px;
}
.modal__btn:not(:disabled):hover { background: #e0af20; }
.modal__error { color: var(--color-error); font-size: 0.875rem; }

/* Room layout */
.room {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.room__header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  background: var(--color-surface);
  border-bottom: 1px solid rgba(255,255,255,0.07);
  flex-wrap: wrap;
}

.room__brand { font-weight: 700; font-size: 1rem; white-space: nowrap; }

.room__share {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
  flex-wrap: wrap;
}

.room__share-label {
  font-size: 0.78rem;
  color: var(--color-text-muted);
  white-space: nowrap;
}

.room__url {
  font-size: 0.78rem;
  background: rgba(255,255,255,0.07);
  border-radius: 4px;
  padding: 3px 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 260px;
}

.room__copy-btn {
  font-size: 0.78rem;
  padding: 4px 10px;
  background: rgba(255,255,255,0.1);
  color: var(--color-text);
  white-space: nowrap;
}
.room__copy-btn--ok { background: #2a9d8f; }

.room__status {
  font-size: 0.75rem;
  padding: 3px 10px;
  border-radius: 999px;
  white-space: nowrap;
}
.room__status--on  { background: rgba(42,157,143,0.3); color: #7ddb9e; }
.room__status--off { background: rgba(230,57,70,0.2);  color: #f4a0a6; }

.room__main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.room__loading {
  color: var(--color-text-muted);
  font-size: 1rem;
}

/* Footer */
.room__footer {
  background: var(--color-surface);
  border-top: 1px solid rgba(255,255,255,0.07);
}

.room__master-panel {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  border-top: 1px solid rgba(255,255,255,0.05);
}

.room__master-label {
  font-size: 0.8rem;
  color: var(--color-text-muted);
  white-space: nowrap;
}

.room__master-select {
  background: rgba(255,255,255,0.08);
  color: var(--color-text);
  border: 1px solid rgba(255,255,255,0.15);
  border-radius: 6px;
  padding: 5px 10px;
  font-family: inherit;
  font-size: 0.9rem;
  cursor: pointer;
}

.room__actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  padding: 12px 20px;
  border-top: 1px solid rgba(255,255,255,0.05);
}

.room__btn {
  padding: 10px 28px;
  font-weight: 700;
  font-size: 0.95rem;
}
.room__btn--reveal { background: var(--color-gold); color: #1a1a1a; }
.room__btn--reveal:not(:disabled):hover { background: #e0af20; }
.room__btn--reset  { background: #2a9d8f; color: white; }
.room__btn--reset:hover { background: #238d80; }

/* Toast */
.toast {
  position: fixed;
  top: 56px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--color-error);
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 0.9rem;
  z-index: 100;
  box-shadow: 0 4px 16px rgba(0,0,0,0.3);
}
.toast-enter-active, .toast-leave-active { transition: opacity 0.3s, transform 0.3s; }
.toast-enter-from  { opacity: 0; transform: translateX(-50%) translateY(-8px); }
.toast-leave-to    { opacity: 0; transform: translateX(-50%) translateY(-8px); }
</style>
