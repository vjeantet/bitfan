go-pop3
==========

This is a simple POP3 client package in Go language.

##Usage
```go
if err := pop3.ReceiveMail(address, user, pass,
	func(number int, uid, data string, err error) (bool, error) {
		log.Printf("%d, %s\n", number, uid)

    // implement your own logic here

		return false, nil
	}); err != nil {
	log.Fatalf("%v\n", err)
}
```
or use the method that implements the command  
```go
client, err := pop3.Dial(address)

if err != nil {
	log.Fatalf("Error: %v\n", err)
}

defer func() {
	client.Quit()
	client.Close()
}()

if err = client.User(user); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

if err = client.Pass(pass); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

var count int
var size uint64

if count, size, err = client.Stat(); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

log.Printf("Count: %d, Size: %d\n", count, size)

if count, size, err = client.List(6); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

log.Printf("Number: %d, Size: %d\n", count, size)

var mis []pop3.MessageInfo

if mis, err = client.ListAll(); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

for _, mi := range mis {
	log.Printf("Number: %d, Size: %d\n", mi.Number, mi.Size)
}

var number int
var uid string

if number, uid, err = client.Uidl(6); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

log.Printf("Number: %d, Uid: %s\n", number, uid)

if mis, err = client.UidlAll(); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

for _, mi := range mis {
	log.Printf("Number: %d, Uid: %s\n", mi.Number, mi.Uid)
}

var content string

if content, err = client.Retr(8); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

log.Printf("Content:\n%s\n", content)

if err = client.Dele(6); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

if err = client.Noop(); err != nil {
	log.Printf("Error: %v\n", err)
	return
}

if err = client.Rset(); err != nil {
	log.Printf("Error: %v\n", err)
	return
}
```

##License
[MIT License](https://github.com/taknb2nch/go-pop3/blob/master/LICENSE)
