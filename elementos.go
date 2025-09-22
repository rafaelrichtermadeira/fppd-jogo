// elementos.go - Inimigos, tesouro e portais
package main

import (
	"math/rand"
	"time"
)

// Portais
var (
	PortalAberto  = Elemento{'🌀', CorVerde, CorPadrao, false, 0}
	PortalFechado = Elemento{'🌀', CorVermelho, CorPadrao, true, 0}
)

// ---------------- Tesouro Vivo ----------------
func iniciarTesouroVivo(x, y int, acoes chan func(*Jogo)) {
	acoes <- func(j *Jogo) {
		if y < len(j.Mapa) && x < len(j.Mapa[y]) {
			j.Mapa[y][x] = Tesouro
		}
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			dxdy := []struct{ dx, dy int }{
				{1, 0}, {-1, 0}, {0, 1}, {0, -1},
			}
			move := dxdy[rand.Intn(len(dxdy))]
			acoes <- func(j *Jogo) {
				for yy := range j.Mapa {
					for xx := range j.Mapa[yy] {
						if j.Mapa[yy][xx].simbolo == Tesouro.simbolo {
							nx, ny := xx+move.dx, yy+move.dy
							if ny < 0 || ny >= len(j.Mapa) || nx < 0 || nx >= len(j.Mapa[ny]) {
								return
							}
							dest := j.Mapa[ny][nx]
							if dest.tangivel || dest.simbolo == InimigoEl.simbolo {
								return
							}
							j.Mapa[yy][xx] = Vazio
							j.Mapa[ny][nx] = Tesouro
							if j.PosX == nx && j.PosY == ny && victoryCh != nil {
								victoryCh <- true
							}
							return
						}
					}
				}
			}
		}
	}()
}

// ---------------- Inimigos ----------------
func iniciarInimigoComID(id, startX, startY int, acoes chan func(*Jogo)) {
	acoes <- func(j *Jogo) {
		el := InimigoEl
		el.id = id
		j.Mapa[startY][startX] = el
	}
	go func() {
		dirs := []struct{ dx, dy int }{
			{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		}
		for {
			time.Sleep(800 * time.Millisecond)
			d := dirs[rand.Intn(len(dirs))]
			idd := id
			acoes <- func(j *Jogo) {
				for yy := range j.Mapa {
					for xx := range j.Mapa[yy] {
						cell := j.Mapa[yy][xx]
						if cell.simbolo == InimigoEl.simbolo && cell.id == idd {
							nx, ny := xx+d.dx, yy+d.dy
							if ny < 0 || ny >= len(j.Mapa) || nx < 0 || nx >= len(j.Mapa[ny]) {
								return
							}
							dest := j.Mapa[ny][nx]
							if dest.tangivel || dest.simbolo == Tesouro.simbolo {
								return
							}
							j.Mapa[yy][xx] = Vazio
							newEl := InimigoEl
							newEl.id = idd
							j.Mapa[ny][nx] = newEl
							if j.PosX == nx && j.PosY == ny && derrotaCh != nil {
								derrotaCh <- true
							}
							return
						}
					}
				}
			}
		}
	}()
}

func iniciarInimigos(acoes chan func(*Jogo)) {
	posicoes := [][3]int{
		{1, 5, 5}, {2, 12, 3}, {3, 8, 10}, {4, 18, 7},
		{5, 22, 4}, {6, 3, 14}, {7, 16, 12}, {8, 25, 6},
		{9, 10, 15}, {10, 20, 10},
	}
	for _, p := range posicoes {
		iniciarInimigoComID(p[0], p[1], p[2], acoes)
	}
}

// ---------------- Portais ----------------
func iniciarPortais(x1, y1, x2, y2 int, acoes chan func(*Jogo)) {
	acoes <- func(j *Jogo) {
		j.Mapa[y1][x1] = PortalAberto
		j.Mapa[y2][x2] = PortalAberto
	}
	go func() {
		aberto := true
		for {
			time.Sleep(3 * time.Second)
			isOpen := aberto
			acoes <- func(j *Jogo) {
				if isOpen {
					j.Mapa[y1][x1] = PortalAberto
					j.Mapa[y2][x2] = PortalAberto
				} else {
					j.Mapa[y1][x1] = PortalFechado
					j.Mapa[y2][x2] = PortalFechado
				}
			}
			aberto = !aberto
		}
	}()
}
