
package main

import (
	"math/rand"
	"time"
)

// Novos elementos
var (
	Tesouro          = Elemento{'$', CorVerde, CorPadrao, false}
	PortalFechado    = Elemento{'🌀', CorVermelho, CorPadrao, true}
	PortalAberto     = Elemento{'🌀', CorVerde, CorPadrao, false}
	ArmadilhaAtiva   = Elemento{'^', CorVermelho, CorPadrao, true}
	ArmadilhaInativa = Vazio
)

// ---------------- Inimigo (☠) ----------------
func iniciarInimigo(x, y int, acoes chan func(*Jogo)) {
	go func() {
		dirs := []struct{ dx, dy int }{
			{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		}
		for {
			time.Sleep(1 * time.Second)
			d := dirs[rand.Intn(len(dirs))]
			acoes <- func(j *Jogo) {
				if jogoPodeMoverPara(j, x+d.dx, y+d.dy) {
					jogoMoverElemento(j, x, y, d.dx, d.dy)
					x, y = x+d.dx, y+d.dy
				}
			}
		}
	}()
}

// ---------------- Tesouro ($) ----------------
func iniciarTesouro(acoes chan func(*Jogo)) {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			px, py := rand.Intn(50), rand.Intn(20)

			acoes <- func(j *Jogo) {
				if !j.Mapa[py][px].tangivel {
					j.Mapa[py][px] = Tesouro
				}
			}

			select {
			case <-time.After(5 * time.Second):
				acoes <- func(j *Jogo) {
					if j.Mapa[py][px].simbolo == Tesouro.simbolo {
						j.Mapa[py][px] = Vazio
					}
				}
			}
		}
	}()
}

// ---------------- Portal (🌀) ----------------
func iniciarPortal(x, y int, acoes chan func(*Jogo), eventos chan string) {
	go func() {
		aberto := false
		for {
			select {
			case <-time.After(3 * time.Second):
				acoes <- func(j *Jogo) {
					if aberto {
						j.Mapa[y][x] = PortalFechado
					} else {
						j.Mapa[y][x] = PortalAberto
					}
					aberto = !aberto
				}
			case msg := <-eventos:
				acoes <- func(j *Jogo) {
					if msg == "forcar-fechar" {
						j.Mapa[y][x] = PortalFechado
						aberto = false
					}
					if msg == "forcar-abrir" {
						j.Mapa[y][x] = PortalAberto
						aberto = true
					}
				}
			}
		}
	}()
}

// ---------------- Armadilha (^) ----------------
func iniciarArmadilha(x, y int, acoes chan func(*Jogo)) {
	go func() {
		ativa := false
		for {
			time.Sleep(4 * time.Second)
			acoes <- func(j *Jogo) {
				if ativa {
					j.Mapa[y][x] = ArmadilhaAtiva
				} else {
					j.Mapa[y][x] = ArmadilhaInativa
				}
				ativa = !ativa

				if ativa && j.PosX == x && j.PosY == y {
					j.StatusMsg = "Você pisou em uma armadilha! (-HP)"
				}
			}
		}
	}()
}
