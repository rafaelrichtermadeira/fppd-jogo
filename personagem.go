// personagem.go - Movimentação e ações do jogador
package main

import "fmt"

func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w', 'W':
		dy = -1
	case 'a', 'A':
		dx = -1
	case 's', 'S':
		dy = 1
	case 'd', 'D':
		dx = 1
	default:
		return
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	if ny < 0 || ny >= len(jogo.Mapa) || nx < 0 || nx >= len(jogo.Mapa[ny]) {
		return
	}

	cell := jogo.Mapa[ny][nx]

	if cell.simbolo == InimigoEl.simbolo && derrotaCh != nil {
		derrotaCh <- true
		return
	}
	if cell.simbolo == Tesouro.simbolo && victoryCh != nil {
		victoryCh <- true
		return
	}
	if cell.tangivel && cell.simbolo == '🌀' {
		return
	}
	if cell.tangivel {
		return
	}

	jogo.PosX, jogo.PosY = nx, ny

	if cell.simbolo == '🌀' && !cell.tangivel {
		for yy := range jogo.Mapa {
			for xx := range jogo.Mapa[yy] {
				if (xx != nx || yy != ny) && jogo.Mapa[yy][xx].simbolo == '🌀' && !jogo.Mapa[yy][xx].tangivel {
					jogo.PosX, jogo.PosY = xx, yy
					jogo.StatusMsg = "Teletransportado!"
					return
				}
			}
		}
	}
}

func personagemInteragir(jogo *Jogo) {
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d,%d)", jogo.PosX, jogo.PosY)
}

func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		personagemInteragir(jogo)
	case "mover":
		personagemMover(ev.Tecla, jogo)
	}
	return true
}
