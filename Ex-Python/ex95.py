pp={}
cad=[]
qntp = 0
pt =[]
xpart = ''
while True:
    print('---'*30)
    pp['Nome']=str(input('Digite o nome do jogador a ser cadastrado :').title().strip())
    qntp = int(input(f'Quantas partidas foram jogadas por {pp['Nome']} :'))
    for c in range (0,qntp):
        pt.append(int(input(f'Quantos gols foram feitos por {pp['Nome']} na {c+1}* partida :')))
    pp['Partidas']=pt[:]
    pp['Total']=int(sum(pt))
    pt.clear()
    cad.append(pp.copy())
    pp.clear()
    r=str(input('Deseja cadastrar mais um jogador (S/N) ?')).strip().upper()
    while r!='N' and r!='S':
        print('Resposta inválida, ', end='')
        r=str(input('Deseja cadastrar mais um jogador (S/N) ?')).strip().upper()
    if r=='N':
        break
print('-=-'*30)
print(f'{'Nome':^30} {'gols':^30} {'total de gols':^30}')
print (f'___'*30)
for c,x in enumerate(cad):
    print(f'{x['Nome']:^30} ', end='|')
    for k in range(0,len(x['Partidas'])):
        xpart += str(cad[c]['Partidas'][k])
    print(f'{xpart:^30}', end='|')
    xpart = ''
    print(f'{cad[c]['Total']:^27}', end ='')
    print('|')
    print('---'*30)
while True:
    r=str(input('Deseja fazer o levantamento de dados de algum jogador ?').upper().strip())
    while r!='S' and r!='N':
        r=str(input('Resposta inválida, tente novamente : '))
    if r=='N':
        break
    ldd = str(input('Digite o nome do jogador a ser analisado :'))
    if ldd not in cad:
        print('Jogador não cadastrado !')
    for x,c in enumerate(cad):
        if c['Nome']==ldd:
            for k in range(0,len(c['Partidas'])):
                print(f'O jogador {ldd} fez {c['Partidas'][k]} gols na {k+1}* partida.')
            print(f'Obtendo um total de {sum(c['Partidas'])/len(c['Partidas'])} média de gols por partida.')
print(f'{'PROGRAMA ENCERRADO':-^90}')