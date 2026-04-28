import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useRoomStore = defineStore('room', () => {
  const room      = ref(null)
  const myId      = ref('')
  const connected = ref(false)
  const error     = ref(null)

  let ws = null

  function connect(roomId, name, creatorId = '') {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const storedId = sessionStorage.getItem(`poker_player_${roomId}`) || creatorId
    const pidParam = storedId ? `&player_id=${encodeURIComponent(storedId)}` : ''
    const url = `${protocol}//${location.host}/api/rooms/${roomId}/ws?name=${encodeURIComponent(name)}${pidParam}`

    ws = new WebSocket(url)

    ws.onopen = () => { connected.value = true; error.value = null }

    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data)
      if (msg.type === 'init') {
        myId.value = msg.player_id
        sessionStorage.setItem(`poker_player_${roomId}`, msg.player_id)
        room.value = msg.room
      } else if (msg.type === 'state') {
        room.value = msg.room
      } else if (msg.type === 'error') {
        error.value = msg.message
        setTimeout(() => { error.value = null }, 4000)
      }
    }

    ws.onclose  = () => { connected.value = false }
    ws.onerror  = () => { error.value = 'Conexão perdida. Recarregue a página.' }
  }

  function disconnect() {
    ws?.close()
    ws = null
    room.value = null
    connected.value = false
    myId.value = ''
  }

  function send(msg) {
    if (ws?.readyState === WebSocket.OPEN) ws.send(JSON.stringify(msg))
  }

  const vote      = (value)    => send({ action: 'vote',       value })
  const reveal    = ()         => send({ action: 'reveal' })
  const reset     = ()         => send({ action: 'reset' })
  const setMaster = (playerId) => send({ action: 'set_master', player_id: playerId })

  const myPlayer  = computed(() => room.value?.players?.find(p => p.id === myId.value) ?? null)
  const isCreator = computed(() => room.value?.creator_id === myId.value)
  const isMaster  = computed(() => room.value?.master_id  === myId.value)
  const allVoted  = computed(() =>
    (room.value?.players?.length ?? 0) > 0 &&
    room.value.players.every(p => p.has_voted)
  )

  const winners = computed(() => {
    if (!room.value?.revealed) return []
    const counts = {}
    for (const p of room.value.players) {
      if (p.vote) counts[p.vote] = (counts[p.vote] || 0) + 1
    }
    const max = Math.max(...Object.values(counts))
    return Object.entries(counts).filter(([, c]) => c === max).map(([v]) => v)
  })

  const average = computed(() => {
    if (!room.value?.revealed) return null
    const nums = room.value.players
      .map(p => parseFloat(p.vote))
      .filter(n => !isNaN(n))
    if (!nums.length) return null
    return (nums.reduce((a, b) => a + b, 0) / nums.length).toFixed(1)
  })

  return {
    room, myId, connected, error,
    myPlayer, isCreator, isMaster, allVoted, winners, average,
    connect, disconnect, vote, reveal, reset, setMaster
  }
})
