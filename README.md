
![logo](logo.jpg)

A bot that allows you to remove spam from chat rooms in Utopia Messenger.

You can use a ready-made bot by adding it to your contacts with a public key: `CA963CF9120FBF1987AB4275524EFFF0BD057FACF659D66C0FAF3D553F7BDD78`

Procedure:
1. add the bot to your contacts;
2. add the bot as a moderator of your chat, giving him the right to delete messages;
3. send the ID of your channel to the bot in your personal message;
4. enable required anti-spam filters.

## build from source

The ready build can be found on the [releases page](https://github.com/Sagleft/utopia-protectron/releases).

```bash
git clone https://github.com/Sagleft/utopia-protectron.git protectron && cd protectron
go build
cp config.example.json config.json
```

to cross-platform build:
```bash
bash build.sh
```

The parameters for connecting to Utopia client are specified according to the example for connecting to the [docker container](https://github.com/Sagleft/utopia-api-docker).

## TODO

1. sticker filter;
2. repetitive message filter;
3. banning a user if he violates too often
4. no-images filter;

## useful links

* [Forum thread](https://talk.u.is/viewtopic.php?pid=5269)
* [uDocs](https://udocs.gitbook.io/utopia-api/)

---
[![udocs](https://github.com/Sagleft/ures/blob/master/udocs-btn.png?raw=true)](https://udocs.gitbook.io/utopia-api/)
