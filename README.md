# fairy-tales-bot

_Telegram bot for a simple fairy tales storage with audio records within a telegram dialog_

## Setting up

Since current project version uses long polling and not a webhook approach, we can run it both locally and on a server.

1. Create telegram bot with using [@BotFather](https://t.me/botfather)
2. Install [Docker](https://docs.docker.com/engine/install/), [Docker Compose](https://docs.docker.com/compose/install/)
3. Clone project
4. Copy environment variables file and edit it (with using [vim](https://www.vim.org/)/[neovim](https://neovim.io/), [emacs](https://www.gnu.org/software/emacs/), [nano](https://www.nano-editor.org/), [micro](https://micro-editor.github.io/), whatever) by assigning proper values:

```bash
cp ./.env.example ./.env
```

5. Run build and run project:

```bash
make up
```

To stop containers run (if you need or want it):

```bash
make down
```

For other operations with containers consider using `docker` command as you usually do. If not, but you need, feel free to check [the docs](https://docs.docker.com/) (_LMAO, ya aint know docker in 2k23?! What the dev r u?! Hope ya know how to use git and bash at least_)

## Interesting facts:

- The bot docker image uses the simplest and smallest possibly docker image base — [scratch](https://hub.docker.com/_/scratch/) as its last build stage, so it has the weight of a compiled binary (and a pem cerificate to use TLS) only and weights only around 13-14MB! I used [this article](https://habr.com/ru/articles/460535/) to do so ([the article to do the same with Rust in docker](https://habr.com/ru/articles/766620/) came out just some days ago)
