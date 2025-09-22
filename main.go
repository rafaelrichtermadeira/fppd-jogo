// main.go - Ponto de entrada
package main

import (
	"math/rand"
	"os"
	"time"
)

var victoryCh chan bool
var derrotaCh chan bool

func main() {
	rand.Seed(time.Now().UnixNano())
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	victoryCh = make(chan bool)
	derrotaCh = make(chan bool)
	acoes := make(chan func(*Jogo), 100)

	// Aplica mudanças concorrentes
	go func() {
		for f := range acoes {
			if f != nil {
				f(&jogo)
			}
			interfaceDesenharJogo(&jogo)
		}
	}()

	// Vitória
	go func() {
		<-victoryCh
		interfaceFinalizar()
		println("🎉 Você pegou o tesouro! Vitória!")
		os.Exit(0)
	}()

	// Derrota
	go func() {
		<-derrotaCh
		interfaceFinalizar()
		println("💀 Um inimigo te pegou! Game Over.")
		os.Exit(0)
	}()

	// Elementos concorrentes
	iniciarInimigos(acoes)               // vários inimigos
	iniciarTesouroVivo(20, 5, acoes)     // tesouro vivo
	iniciarPortais(15, 8, 18, 5, acoes)  // portais

	interfaceDesenharJogo(&jogo)

	for {
		ev := interfaceLerEventoTeclado()
		if !personagemExecutarAcao(ev, &jogo) {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
