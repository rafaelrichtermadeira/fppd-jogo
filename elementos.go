// elementos.go - Elementos concorrentes do jogo
package main

import (
	"math/rand"
	"time"
)

// Novos elementos
var (
	Tesouro = Elemento{'$', CorVerde, CorPadrao, false}

	// Portal: mesmo símbolo, mas toggla tangivel para bloquear/permissão.
	PortalAberto  = Elemento{'🌀', CorVerde, CorPadrao, false}
	PortalFechado = Elemento{'🌀', CorVermelho, CorPadrao, true}
)

// ---------------- Inimigo (☠) ----------------
func iniciarInimigo(startX, startY int, acoes chan func(*Jogo)) {
	go func() {
		x, y := startX, startY
		dirs := []struct{ dx, dy int }{
			{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		}
		for {
			time.Sleep(1 * time.Second)
			d := dirs[rand.Intn(len(dirs))]
			// capture valores locais para evitar efeitos de captura de loop
			px, py := x, y
			ddx, ddy := d.dx, d.dy
			acoes <- func(j *Jogo) {
				if jogoPodeMoverPara(j, px+ddx, py+ddy) {
					jogoMoverElemento(j, px, py, ddx, ddy)
					// atualiza posição local (variável externa ao closure)
					x = px + ddx
					y = py + ddy
				}
			}
		}
	}()
}

// ---------------- Tesouro ($) ----------------
// gera UM tesouro por vez; desaparece após timeout se não coletado
func iniciarTesouro(acoes chan func(*Jogo)) {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			coordsCh := make(chan [2]int, 1)

			// Pedido para criar um tesouro (será executado pelo applier)
			acoes <- func(j *Jogo) {
				// remove tesouros antigos
				for yy := range j.Mapa {
					for xx := range j.Mapa[yy] {
						if j.Mapa[yy][xx].simbolo == Tesouro.simbolo {
							j.Mapa[yy][xx] = Vazio
						}
					}
				}
				// tenta achar uma posição livre
				h := len(j.Mapa)
				if h == 0 {
					return
				}
				w := len(j.Mapa[0])
				for attempts := 0; attempts < 1000; attempts++ {
					px := rand.Intn(w)
					py := rand.Intn(h)
					// evita colocar em cima do jogador
					if (px == j.PosX && py == j.PosY) {
						continue
					}
					if !j.Mapa[py][px].tangivel && j.Mapa[py][px].simbolo != Tesouro.simbolo {
						j.Mapa[py][px] = Tesouro
						coordsCh <- [2]int{px, py}
						return
					}
				}
				// se não achar, não faz nada
			}

			// espera as coords colocadas quando o closure foi executado
			coords, ok := <-coordsCh
			if !ok {
				// canal fechado inesperado -> segue loop
				continue
			}

			// Após 5s, remove o tesouro se ainda existir
			time.Sleep(5 * time.Second)
			px := coords[0]
			py := coords[1]
			acoes <- func(j *Jogo) {
				if py < len(j.Mapa) && px < len(j.Mapa[py]) {
					if j.Mapa[py][px].simbolo == Tesouro.simbolo {
						j.Mapa[py][px] = Vazio
					}
				}
			}
		}
	}()
}

// ---------------- Portais (entrada/saída) ----------------
// cria dois portais e alterna entre aberto/fechado juntos
func iniciarPortais(x1, y1, x2, y2 int, acoes chan func(*Jogo)) {
	// coloca portais inicialmente abertos
	acoes <- func(j *Jogo) {
		if y1 < len(j.Mapa) && x1 < len(j.Mapa[y1]) {
			j.Mapa[y1][x1] = PortalAberto
		}
		if y2 < len(j.Mapa) && x2 < len(j.Mapa[y2]) {
			j.Mapa[y2][x2] = PortalAberto
		}
	}

	go func() {
		aberto := true
		for {
			time.Sleep(3 * time.Second)
			px1, py1, px2, py2 := x1, y1, x2, y2
			isOpen := aberto
			acoes <- func(j *Jogo) {
				if py1 < len(j.Mapa) && px1 < len(j.Mapa[py1]) {
					if isOpen {
						j.Mapa[py1][px1] = PortalAberto
					} else {
						j.Mapa[py1][px1] = PortalFechado
					}
				}
				if py2 < len(j.Mapa) && px2 < len(j.Mapa[py2]) {
					if isOpen {
						j.Mapa[py2][px2] = PortalAberto
					} else {
						j.Mapa[py2][px2] = PortalFechado
					}
				}
			}
			aberto = !aberto
		}
	}()
}
