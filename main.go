package main

import (
	"math/rand"
	"os"
	"time"
)

func main() {
	// seed aleatório
	rand.Seed(time.Now().UnixNano())

	// Inicializa a interface
	interfaceIniciar()
	defer interfaceFinalizar()

	// Arquivo do mapa
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// canais
	victoryCh = make(chan bool)
	acoes := make(chan func(*Jogo), 50) // buffer para reduzir chance de bloqueio

	// Goroutine que consome ações e aplica no estado (exclusão mútua via canal)
	go func() {
		for f := range acoes {
			if f != nil {
				f(&jogo)
			}
			interfaceDesenharJogo(&jogo)
		}
	}()

	// Goroutine que aguarda vitória
	go func() {
		<-victoryCh
		// finalize e saia
		interfaceFinalizar()
		println("🎉 Parabéns, você pegou o tesouro! (vitória)")
		os.Exit(0)
	}()

	// Inicia elementos concorrentes
	iniciarInimigo(10, 5, acoes)
	iniciarTesouro(acoes)
	iniciarPortais(15, 8, 2, 12, acoes) // entrada em (15,8), saída em (2,12)

	// Desenha estado inicial
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada: pega eventos e os transforma em ações via canal
	for {
		ev := interfaceLerEventoTeclado()
		if ev.Tipo == "sair" {
			// finaliza imediatamente
			interfaceFinalizar()
			os.Exit(0)
		}
		cont := personagemProcessarEvento(ev, acoes)
		if !cont {
			interfaceFinalizar()
			os.Exit(0)
		}
		// não desenhamos aqui — a goroutine que aplica ações já redesenha
	}
}
