// jogo.go - Estruturas centrais do jogo
package main

import (
	"bufio"
	"os"
)

// Elemento representa qualquer objeto do mapa
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool
	id       int // usado para identificar inimigos
}

// Estrutura principal do jogo
type Jogo struct {
	Mapa      [][]Elemento
	PosX      int
	PosY      int
	StatusMsg string
}

// Elementos do jogo
var (
	Personagem = Elemento{'☺', CorCinza, CorPadrao, true, 0}
	InimigoEl  = Elemento{'☠', CorVermelho, CorPadrao, false, 0}
	Parede     = Elemento{'▤', CorParede, CorPadrao, true, 0}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false, 0}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false, 0}
	Tesouro    = Elemento{'$', CorVerde, CorPadrao, false, 0}
)

func jogoNovo() Jogo {
	return Jogo{}
}

// Carrega mapa de arquivo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y
				e = Vazio
			case Tesouro.simbolo:
				e = Tesouro
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
