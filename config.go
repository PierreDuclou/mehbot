package main

var config = struct {
	prefix       string            // command prefix
	baseroles    map[string]string // known discord roles
	messageColor int
	successColor int
	errorColor   int
}{
	prefix: "!",
	baseroles: map[string]string{
		"Guez":      "131380194207465472",
		"Guezt":     "241631468432916480",
		"Worms":     "598963803014561857",
		"Superguez": "336920690030673921",
	},
	messageColor: 0x5f6eff,
	successColor: 0x34cc41,
	errorColor:   0xf34747,
}

type role struct {
	id   string
	name string
}
