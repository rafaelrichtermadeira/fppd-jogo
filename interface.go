// interface.go - Interface gráfica do jogo usando termbox
// Implementa a interface gráfica do jogo com a biblioteca termbox-go.

package main

import (
	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsular as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo (apenas cores válidas em termbox)
const (
	CorPadrao   Cor = termbox.ColorDefault
	CorCinza    Cor = termbox.ColorWhite    // substituto para "cinza"
	CorVermelho Cor = termbox.ColorRed
	CorVerde    Cor = termbox.ColorGreen
	CorParede   Cor = termbox.ColorBlue     // parede azul
	CorFundo    Cor = termbox.ColorBlack
	CorTexto    Cor = termbox.ColorYellow
)

// EventoTeclado representa uma ação detectada do teclado
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' || ev.Ch == 'E' {
		return EventoTeclado{Tipo: "interagir"}
	}
	if ev.Ch == 'w' || ev.Ch == 'a' || ev.Ch == 's' || ev.Ch == 'd' {
		return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
	}
	return EventoTeclado{}
}

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha o personagem sobre o mapa
	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)

	// Atualiza a tela
	interfaceAtualizarTela()
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	// Linha de status dinâmica
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	// Instruções fixas
	msg := "Use WASD para mover, E para interagir, ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}
