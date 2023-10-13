Basicamente esse projeto consome a API da Blaze do Jogo Double.
O algoritmo indetifica quando se esta disponível para apostar, salva os resultado e recebe so sinais por telegram.
O sinais são recebido por UDP, e são enviados por projeto em python aqui no meu github.
Foi usado UDP pois foi a forma simples e fácil de trabalhar com a questão assincrona.
Esse projeto tem como base o estudo da programação async em golang e apenas isso.
Pois ao mesmo tempo que se abre a janela para jogar se recebe sinais vindo por UDP oriundos do telegram.
