<template>
  <div class="table-wrap">
    <div class="table">
      <!-- players positioned around the rim -->
      <PlayerSeat
        v-for="(player, i) in players"
        :key="player.id"
        :player="player"
        :is-master="player.id === masterid"
        :is-me="player.id === myId"
        :revealed="revealed"
        :winners="winners"
        :style="seatStyle(i, players.length)"
        class="table__seat"
      />

      <!-- centre panel -->
      <div class="table__center">
        <template v-if="revealed">
          <div class="center__result">
            <div class="center__winners">
              <span
                v-for="w in winners"
                :key="w"
                class="center__winner-value"
              >{{ w }}</span>
            </div>
            <div class="center__label">mais votado</div>
            <div v-if="average" class="center__avg">média: {{ average }}</div>
          </div>
        </template>

        <template v-else>
          <div class="center__status">
            <div class="center__round">Rodada {{ round }}</div>
            <div class="center__votes">
              {{ votedCount }} / {{ players.length }}
              <span class="center__votes-label">votaram</span>
            </div>
            <div v-if="allVoted" class="center__hint">pronto!</div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import PlayerSeat from './PlayerSeat.vue'

const props = defineProps({
  players:  { type: Array,   required: true },
  masterid: { type: String,  default: '' },
  myId:     { type: String,  default: '' },
  revealed: { type: Boolean, default: false },
  round:    { type: Number,  default: 1 },
  winners:  { type: Array,   default: () => [] },
  average:  { type: String,  default: null },
  allVoted: { type: Boolean, default: false }
})

const votedCount = computed(() => props.players.filter(p => p.has_voted).length)

// Position i-th seat around a circle. Radius is defined in CSS via --table-r.
function seatStyle(i, total) {
  const angle = (i / Math.max(total, 1)) * 2 * Math.PI - Math.PI / 2
  const r = 210   // px, must match --table-r in CSS
  const x = Math.round(r * Math.cos(angle))
  const y = Math.round(r * Math.sin(angle))
  return {
    position: 'absolute',
    left: '50%',
    top: '50%',
    transform: `translate(calc(-50% + ${x}px), calc(-50% + ${y}px))`
  }
}
</script>

<style scoped>
.table-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 60px 20px;
}

.table {
  position: relative;
  width: 500px;
  height: 500px;
  border-radius: 50%;
  background: radial-gradient(circle at 40% 35%, #2d9b5f, var(--color-table) 60%, #134a30);
  box-shadow:
    0 0 0 14px var(--color-table-rim),
    0 0 0 18px #5c4609,
    0 30px 80px rgba(0,0,0,0.6);
  flex-shrink: 0;
}

.table__seat { z-index: 2; }

.table__center {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.center__result,
.center__status {
  text-align: center;
  color: rgba(255,255,255,0.95);
  text-shadow: 0 1px 4px rgba(0,0,0,0.5);
}

.center__winners {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 4px;
}

.center__winner-value {
  font-size: 3rem;
  font-weight: 800;
  color: var(--color-gold);
  text-shadow: 0 2px 8px rgba(0,0,0,0.5);
  line-height: 1;
}

.center__label {
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: rgba(255,255,255,0.6);
}

.center__avg {
  margin-top: 6px;
  font-size: 0.85rem;
  color: rgba(255,255,255,0.7);
}

.center__round {
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(255,255,255,0.5);
  margin-bottom: 8px;
}

.center__votes {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1;
}

.center__votes-label {
  display: block;
  font-size: 0.75rem;
  font-weight: 400;
  color: rgba(255,255,255,0.6);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.center__hint {
  margin-top: 8px;
  font-size: 0.9rem;
  color: #7ddb9e;
}

@media (max-width: 600px) {
  .table { width: 340px; height: 340px; }
  .table-wrap { padding: 40px 10px; }
}
</style>
