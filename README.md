# sayonara

leave all discord group dms or only ones you were added to by a specific user. fork of [leavemealone](https://github.com/diamondburned/leavemealone).

## installation

```sh
go install github.com/fawni/sayonara@latest
```

## usage

```sh
sayonara -t "TOKEN" [-u USERID]
```

## notes

- you will be shown which groups to be left and asked for confirmation beforehand.
- USERID is optional. if provided, only groups whose owner's id is equal to USERID will be left.
