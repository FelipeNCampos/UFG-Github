pp={}
gols = 0
g=[]
pp['Nome']=input('Digite o nome do jogador a ser analisado :').strip().title()
qntp=int(input(f'Digite a quantidade de partidas jogadas por {pp['Nome']}:'))
for c in range(0,qntp):
    g.append(int(input(f'Digite a quantidade de golds feitos na {c+1}* partida {c+1} :')))
    gols += g[c]
pp['Gols']=g
pp['totalg']=gols
print('-=-'*10)
print(pp)
for c, x in enumerate(pp['Gols']):
    print(f'Na partida {c+1},{pp['Nome']} fez {pp['Gols'][c]} gols. ')
print(f'{pp['Nome']} teve um total de {pp['totalg']} gols nesse campeonato .')