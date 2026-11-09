def ficha(a,b):
    if a=='':
        a='<DESCONHECIDO>'
    if b=='':
        b=0
    print(f'O jogador {a} fez {b} gol(s) no campeonato.', end='')



b=0
print('---'*10)
n=str(input('Digite o nome do jogador a ser cadastrado :')).title().strip()
g=str(input(f'Digite a quantidade de gols feitos por {n} :'))
ficha(n,g)