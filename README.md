# scdownload

Golang cli to download latest assets of your Supercell game.

## Installation
Using Scdownload is easy. First, use go get to install the latest version of the library.
```console
go get -u github.com/yautah/scdownload@latest
```

## Windows User
I have a compiled binary for windows users, if you don't have golang environment.
- [scdownload.exe](https://www.aliyundrive.com/s/NYTL8AsnHr9)

## ToDo
- Add file extension filter.
- Support multi threads download.
- Add log files for each download/update task.
- Support other Supercell's games, eg: Clash of Clan, Brawl Stars.

## Usage
### scdownload clone
The `scdownload clone` command will download supercell game assets for you.
```console
克隆一个完整的assets到本地

Usage:
  scdownload clone [flags]

Global Flags:
  -g, --game string   指定游戏类型 (default "cr")

Flags:
  -d, --domain string        资源cdn域名 (default "game-assets.clashroyaleapp.com")
  -f, --fingerprint string   fingerprint文件中的hash值 (default "acf932573295414ef92479e9240aecb0854a70a7")
  -o, --output string        资源下载路径 (default "./")
  -e, --extension string     todo: 仅下载指定扩展名的文件 (default "all")
  -h, --help                 help for clone

```

### scdownload pull
The `scdownload pull` command will checked your downloaded files, and update your files from the latest fingerprint files.

Just enter the download assets directory, execute `scdownload pull`.
```console
更新当前资源包内资源

Usage:
  scdownload pull [flags]

Flags:
  -h, --help   help for pull
```

## Supercell Game Cdn lists
* Clash Royale: game-assets.clashroyaleapp.com
* Clash Royale China: cr-cdn.supercellgame.cn
