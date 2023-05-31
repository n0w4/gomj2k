# GOMJ2K

this project is a very simple way to populate topics on kafka with json data.

this application do not be intended to read files directly, just stdin

## Example usage

```bash
cat /some/file.txt | gomj2k -bootstrap-server localhost:9092  
```

## With structured json

need be a json per line without any pretty thing  

```json
{
    "topic": "log",
    "headers": [{"header1":"value1"}, {"header2":"value2"}],
    "key": "some key",
    "payload": {...}
}
```

the correct way

```json
{"topic": "log","headers": [{"header1":"value1"},...],"key": "some key","payload": {...}}
```

When used structured json with a topic field, can be used the flag **-to-topic** to send a copy of all messages to a extra topic

This example send messages to topic **user.profile** and **user.choices**:

```json
{"topic": "user.profile","headers": [{"header1":"value1"},...],"key": "some key","payload": {...}}
{"topic": "user.choices","headers": [{"header1":"value1"},...],"key": "some key","payload": {...}}
```

adding the flag **-to-topic log**, all messages are sent to topic **log** too

```bash
head /some/file.txt | gomj2k -bootstrap-server localhost:9092 -topic log
```

## When used in free mode

in this way each message will be put on topic as is, but **-to-topic ...** and **-free-mode** need be used together, the message key and headers cannot be used and will be implemented soon

```bash
head /some/file.txt | gomj2k -bootstrap-server localhost:9092  -free-mode -topic some.topic
```

### **Atention:** This solve some personal needs, and some improvements can be delayed sometimes
