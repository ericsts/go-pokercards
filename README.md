# 🃏 Planning Poker

> Versão em inglês: [README.en.md](README.en.md)

Aplicação de **Planning Poker** (Scrum Poker) para times ágeis estimarem histórias em tempo real. Criada com o objetivo de aprendizado através de um exemplo prático completo: backend em **Go**, frontend em **Vue 3**, comunicação via **WebSocket**, contêineres com **Docker Compose** e deploy automatizado com **GitHub Actions**.

---

## Sumário

- [Como funciona](#como-funciona)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
- [Estrutura do projeto](#estrutura-do-projeto)
- [Rodando localmente](#rodando-localmente)
- [Rodando com Docker](#rodando-com-docker)
- [Deploy no droplet](#deploy-no-droplet)
- [CI/CD com GitHub Actions](#cicd-com-github-actions)
- [Protocolo WebSocket](#protocolo-websocket)
- [Decisões de design](#decisões-de-design)

---

## Como funciona

1. O **criador** entra na home, digita seu nome e clica em **"Criar nova sala"**.
2. Ele recebe um link que pode compartilhar com o time.
3. Cada participante acessa o link, digita seu nome e entra na sala.
4. O criador pode delegar o papel de **Scrum Master** a qualquer pessoa.
5. Todos selecionam uma carta em segredo (os votos ficam ocultos).
6. O Scrum Master clica em **"Revelar cartas"** — todos os votos aparecem simultaneamente.
7. A carta mais votada e a média numérica são exibidos na mesa.
8. O Scrum Master clica em **"Nova rodada"** para reiniciar.

---

## Tecnologias

| Camada | Tecnologia | Por quê |
|---|---|---|
| Backend | Go 1.22+ | Concorrência nativa, excelente para WebSocket |
| WebSocket | gorilla/websocket | Biblioteca padrão da comunidade Go |
| IDs únicos | google/uuid | Geração de UUIDs v4 para salas e jogadores |
| Frontend | Vue 3 + Vite | Reatividade fina com Composition API, build rápido |
| Estado | Pinia | Store oficial do Vue 3, simples e tipado |
| Roteamento | Vue Router 4 | Histórico HTML5, rotas `/room/:id` |
| Contêineres | Docker + Compose | Isolamento, build reproduzível |
| Proxy | nginx | Serve o frontend e encaminha `/api/` ao backend |
| CI/CD | GitHub Actions | Deploy automático a cada push na `main` |

---

## Arquitetura

```
Browser  ──HTTP──►  nginx (porta 3001 no droplet)
                      │
                      ├── /          ──► Frontend (SPA Vue 3, arquivos estáticos)
                      └── /api/      ──► Backend Go (porta 8080, interna ao Docker)
                            │
                            ├── POST /api/rooms         cria sala
                            ├── GET  /api/rooms/:id     consulta sala
                            ├── GET  /api/cards         lista valores válidos
                            └── GET  /api/rooms/:id/ws  upgrade para WebSocket
```

### Pacotes do backend

```
internal/
  room/       Lógica pura de negócio. Nenhuma dependência de rede.
              Room, Player, View + métodos Vote, Reveal, Reset, SetMaster.

  session/    Camada de infraestrutura WebSocket.
              Client  — envolve uma *websocket.Conn com goroutine de escrita dedicada.
              Session — agrupa Room + mapa de Clients + mutex embutido.
              Store   — mapa thread-safe de roomID → *Session.

  handler/    HTTP e WebSocket. Usa session.Store.
              Rotas registradas no http.ServeMux do Go 1.22.
```

**Por que separar `room` de `session`?**
O pacote `room` não importa nada de WebSocket. Isso permite testar toda a lógica de jogo com testes unitários simples, sem nenhuma conexão de rede. O pacote `session` é o único que conhece gorilla/websocket.

### Componentes do frontend

```
views/
  HomeView    Formulário de criação de sala.
  RoomView    Sala completa: mesa, baralho, controles.

components/
  RoundTable  Mesa redonda em CSS puro. Posiciona jogadores com trigonometria.
  PlayerSeat  Assento de cada jogador: avatar, carta, nome, coroa de Scrum Master.
  CardDeck    Baralho selecionável. Mostra badge da carta escolhida.
  PokerCard   Carta individual: frente (valor), verso (votando), vazia, vencedora.

stores/
  room.js     Pinia store: abre/fecha WebSocket, processa mensagens, expõe estado reativo.
```

---

## Estrutura do projeto

```
go-pokercards/
├── .github/
│   └── workflows/
│       └── deploy.yml          Pipeline CI/CD
├── backend/
│   ├── cmd/server/main.go      Entry point
│   ├── internal/
│   │   ├── room/               Lógica de negócio + testes
│   │   ├── session/            WebSocket (Client, Session, Store) + testes
│   │   └── handler/            HTTP + WS handlers + testes
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── views/              HomeView, RoomView
│   │   ├── components/         RoundTable, PlayerSeat, CardDeck, PokerCard
│   │   ├── stores/room.js      WebSocket + estado global
│   │   └── router/index.js
│   ├── nginx.conf              Proxy interno do contêiner
│   ├── Dockerfile
│   └── package.json
├── nginx/
│   └── poker.ericsantos.eu.conf   Bloco nginx para o droplet
└── docker-compose.yml
```

---

## Rodando localmente

### Pré-requisitos

- Go 1.22+
- Node.js 20+

### Backend

```bash
cd backend
go run ./cmd/server
# Servidor sobe em http://localhost:8080
```

### Frontend

O `vite.config.js` já está configurado para redirecionar `/api/` → `localhost:8080`, então não é necessário nenhuma configuração extra.

```bash
cd frontend
npm install
npm run dev
# Aplicação disponível em http://localhost:5173
```

### Rodando os testes Go

```bash
cd backend
go test ./...
```

---

## Rodando com Docker

```bash
# Na raiz do projeto
docker compose up --build
# Acesse http://localhost:3001
```

Para parar:

```bash
docker compose down
```

---

## Deploy no droplet

### Configuração inicial (feita uma única vez)

```bash
# 1. No seu droplet, instale Docker (se não tiver)
apt update && apt install -y docker.io docker-compose-plugin git

# 2. Clone o repositório
git clone https://github.com/<seu-usuario>/go-pokercards.git /opt/go-pokercards

# 3. Primeiro build
cd /opt/go-pokercards
docker compose up -d --build
```

### nginx no droplet

Copie o arquivo de configuração e ative o subdomínio:

```bash
cp nginx/poker.ericsantos.eu.conf /etc/nginx/sites-available/
ln -s /etc/nginx/sites-available/poker.ericsantos.eu.conf /etc/nginx/sites-enabled/

# Emite o certificado TLS com Certbot
certbot --nginx -d poker.ericsantos.eu

nginx -s reload
```

> O arquivo `nginx/poker.ericsantos.eu.conf` já contém os blocos HTTP e HTTPS com os headers necessários para o WebSocket (`Upgrade`, `Connection`).

---

## CI/CD com GitHub Actions

A cada `git push origin main`:

1. **Job `test`** — executa `go test ./...` no runner do GitHub.
2. **Job `deploy`** (só roda se os testes passaram) — faz SSH no droplet, executa `git pull` e `docker compose up --build`.

### Secrets necessários no GitHub

Vá em **Settings → Secrets and variables → Actions** do seu repositório e cadastre:

| Secret | Valor |
|---|---|
| `DROPLET_HOST` | IP ou hostname do droplet |
| `DROPLET_USER` | Usuário SSH (ex: `root`) |
| `DROPLET_SSH_KEY` | Conteúdo da chave privada SSH (ex: `~/.ssh/id_rsa`) |

> O caminho do projeto no droplet está fixo como `/opt/go-pokercards` no workflow. Se você clonou em outro lugar, edite a linha `cd /opt/go-pokercards` em [.github/workflows/deploy.yml](.github/workflows/deploy.yml).

### Como adicionar a chave SSH

```bash
# Gere um par de chaves dedicado para o deploy (opcional mas recomendado)
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/deploy_key

# Adicione a chave pública ao droplet
ssh-copy-id -i ~/.ssh/deploy_key.pub root@<ip-do-droplet>

# O conteúdo de deploy_key (privada) vai no secret DROPLET_SSH_KEY
cat ~/.ssh/deploy_key
```

---

## Protocolo WebSocket

A conexão é aberta em `GET /api/rooms/:id/ws?name=Alice[&player_id=uuid]`.

### Servidor → Cliente

| `type` | Quando | Campos extras |
|---|---|---|
| `init` | Logo após conectar | `player_id`, `room` |
| `state` | A cada mudança de estado | `room` |
| `error` | Ação inválida | `message` |

### Cliente → Servidor

| `action` | Quem pode | Campos extras |
|---|---|---|
| `vote` | Qualquer jogador | `value` (ex: `"8"`) |
| `reveal` | Scrum Master | — |
| `reset` | Scrum Master | — |
| `set_master` | Criador da sala | `player_id` |

### Objeto `room`

```jsonc
{
  "id": "a1b2c3d4",
  "creator_id": "uuid",
  "master_id":  "uuid",
  "revealed":   false,
  "round":      1,
  "players": [
    { "id": "uuid", "name": "Alice", "has_voted": true,  "vote": "" },
    { "id": "uuid", "name": "Bob",   "has_voted": false, "vote": "" }
  ]
}
// Quando revealed=true, o campo "vote" é preenchido com o valor escolhido.
```

---

## Decisões de design

### Goroutine de escrita por conexão

O gorilla/websocket **não permite escritas concorrentes** na mesma conexão. A solução padrão é ter uma goroutine dedicada por cliente que drena um canal e escreve em série:

```go
// session/client.go
func (c *Client) writePump() {
    for msg := range c.send {   // bloqueia até chegar mensagem
        c.conn.WriteJSON(msg)
    }
}
```

Broadcasts simplesmente enviam para o canal (`c.send <- msg`) sem se preocupar com concorrência.

### `PrepareClient` antes de `RegisterClient`

O cliente novo precisa receber o `init` (com seu próprio ID) **antes** de qualquer broadcast. Por isso o handler cria o cliente, envia o init e só então o registra na sessão:

```go
client := sess.PrepareClient(conn)  // goroutine de escrita sobe aqui
client.Send(initMsg)                // init vai para o canal antes de qualquer broadcast
sess.RegisterClient(playerID, client)
sess.BroadcastState()               // agora sim aparece para todos
```

### Mesa redonda em CSS puro

Não há biblioteca de layout. A posição de cada jogador é calculada com trigonometria básica:

```js
// components/RoundTable.vue
function seatStyle(i, total) {
  const angle = (i / total) * 2 * Math.PI - Math.PI / 2  // começa no topo
  const r = 210  // raio em px
  return {
    transform: `translate(calc(-50% + ${Math.round(r * Math.cos(angle))}px),
                          calc(-50% + ${Math.round(r * Math.sin(angle))}px))`
  }
}
```

### Votos ocultos até o reveal

O método `ToView()` no backend nunca expõe o valor do voto enquanto `Revealed == false` — só o campo `has_voted`. Isso garante que nem mesmo inspecionando o tráfego de rede um jogador vê os votos dos outros antes da revelação.

### Estado em memória

Não há banco de dados. O estado de todas as salas vive no processo Go. Isso simplifica o projeto para fins de aprendizado, mas significa que as salas são perdidas se o servidor reiniciar. Para produção real, bastaria substituir o `Store` por uma implementação com Redis.
