package main

import "fmt"

// Canal global de vitória
var victoryCh chan bool

// Enqueue event -> cria uma função que modifica o estado do jogo e envia para acoes.
// Retorna false se for evento de sair (para o main finalizar).
func personagemProcessarEvento(ev EventoTeclado, acoes chan func(*Jogo)) bool {
	switch ev.Tipo {
	case "sair":
		// main trata sair chamando interfaceFinalizar e os.Exit
		return false
	case "interagir":
		acoes <- func(j *Jogo) {
			j.StatusMsg = fmt.Sprintf("Interagindo em (%d,%d)", j.PosX, j.PosY)
		}
	case "mover":
		var dx, dy int
		switch ev.Tecla {
		case 'w', 'W':
			dy = -1
		case 'a', 'A':
			dx = -1
		case 's', 'S':
			dy = 1
		case 'd', 'D':
			dx = 1
		default:
			// tecla inválida
			return true
		}

		// envia a ação ao canal (será executada pelo applier)
		acoes <- func(j *Jogo) {
			nx := j.PosX + dx
			ny := j.PosY + dy

			// verifica limites e tangibilidade
			if ny < 0 || ny >= len(j.Mapa) || nx < 0 || nx >= len(j.Mapa[ny]) {
				return
			}
			elem := j.Mapa[ny][nx]
			if elem.tangivel {
				// bloqueado
				j.StatusMsg = "Caminho bloqueado!"
				return
			}

			// movimento válido: atualiza posição do jogador
			j.PosX = nx
			j.PosY = ny
			j.StatusMsg = "" // limpa mensagem

			// se pisou em tesouro
			if elem.simbolo == Tesouro.simbolo {
				// remove e sinaliza vitória
				j.Mapa[ny][nx] = Vazio
				if victoryCh != nil {
					// envio não bloqueante: há goroutine esperando no main
					victoryCh <- true
				}
				return
			}

			// se pisou em portal aberto -> teleporta para o outro portal aberto
			if elem.simbolo == PortalAberto.simbolo && !elem.tangivel {
				// busca outro portal aberto
				for yy := range j.Mapa {
					for xx := range j.Mapa[yy] {
						p := j.Mapa[yy][xx]
						if xx == nx && yy == ny {
							continue
						}
						if p.simbolo == PortalAberto.simbolo && !p.tangivel {
							// teleportar jogador para (xx,yy)
							j.PosX = xx
							j.PosY = yy
							j.StatusMsg = "Teletransportado!"
							return
						}
					}
				}
			}
		}
	}
	return true
}
